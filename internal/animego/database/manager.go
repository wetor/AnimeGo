package database

import (
	"path"

	"github.com/pkg/errors"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/dirdb"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

func (m *Database) Scan() error {
	m.dirMutex.Lock()
	defer m.dirMutex.Unlock()
	m.Lock()
	defer m.Unlock()
	d, err := dirdb.Open(Conf.SavePath)
	if err != nil {
		return errors.Wrap(err, "扫描文件数据库失败")
	}
	files, err := d.ScanAll()
	if err != nil {
		return errors.Wrap(err, "扫描文件数据库失败")
	}
	for _, file := range files {
		err = file.Open()
		if err != nil {
			log.DebugErr(err)
			log.Warnf("打开数据文件失败: %s", file.File)
			continue
		}
		switch file.Ext {
		case path.Ext(AnimeDBName):
			anime := &models.AnimeDBEntity{}
			err = file.DB.Unmarshal(anime)
			if err != nil {
				log.DebugErr(err)
				log.Warnf("读取数据文件失败: %s", file.File)
				break
			}
			m.dir2name[file.Dir] = anime.Name
			if _, ok := m.name2dir[anime.Name]; ok {
				m.name2dir[anime.Name].Dir = file.Dir
			} else {
				m.name2dir[anime.Name] = &AnimeDir{
					Dir:       file.Dir,
					SeasonDir: make(map[int]string),
				}
			}
			if CacheMode {
				m.cacheAnimeDBEntity[anime.Name] = anime
			}
		case path.Ext(SeasonDBName):
			season := &models.SeasonDBEntity{}
			err = file.DB.Unmarshal(season)
			if err != nil {
				log.DebugErr(err)
				log.Warnf("读取数据文件失败: %s", file.File)
				break
			}
			m.dir2name[file.Dir] = season.Name
			if _, ok := m.name2dir[season.Name]; ok {
				m.name2dir[season.Name].SeasonDir[season.Season] = file.Dir
			} else {
				m.name2dir[season.Name] = &AnimeDir{
					Dir:       path.Dir(file.Dir),
					SeasonDir: make(map[int]string),
				}
				m.name2dir[season.Name].SeasonDir[season.Season] = file.Dir
			}
			if CacheMode {
				if m.cacheSeasonDBEntity[season.Name] == nil {
					m.cacheSeasonDBEntity[season.Name] = make(map[int]*models.SeasonDBEntity)
				}
				m.cacheSeasonDBEntity[season.Name][season.Season] = season
			}
		}

		err = file.Close()
		if err != nil {
			log.DebugErr(err)
			log.Warnf("关闭数据文件失败: %s", file.File)
			continue
		}
	}
	return nil
}

// write
//
//	写入文件数据库
func (m *Database) write(file string, data any) error {
	f := dirdb.NewFile(file)
	err := f.Open()
	if err != nil {
		log.DebugErr(err)
		return err
	}
	defer f.Close()
	err = f.DB.Marshal(data)
	if err != nil {
		log.DebugErr(err)
		return err
	}
	log.Debugf("write %s: %+v", file, data)
	return nil
}

// read
//
//	读取文件数据库
func (m *Database) read(file string, data any) error {
	f := dirdb.NewFile(file)
	err := f.Open()
	if err != nil {
		log.DebugErr(err)
		return err
	}
	defer f.Close()
	err = f.DB.Unmarshal(data)
	if err != nil {
		log.DebugErr(err)
		return err
	}
	log.Debugf("read %s: %+v", file, data)
	return nil
}

// getAnimeEntityByHash
//
//	获取AnimeEntity
//	从内存或bolt中获取
func (m *Database) getAnimeEntity(hash string) (*models.AnimeEntity, error) {
	var name string
	err := m.cache.Get(Hash2NameBucket, hash, &name)
	if err != nil {
		return nil, err
	}
	return m.getAnimeEntityByName(name)
}

// getAnimeEntityByName
//
//	获取AnimeEntity，使用name
//	从内存或bolt中获取
func (m *Database) getAnimeEntityByName(name string) (*models.AnimeEntity, error) {
	anime := &models.AnimeEntity{}
	err := m.cache.Get(Name2EntityBucket, name, anime)
	if err != nil {
		return nil, err
	}
	return anime, nil
}

// getAnimeDBEntityByDir
//
//	获取文件数据库AnimeDBEntity，使用dir
//	从内存或文件数据库中获取
func (m *Database) getAnimeDBEntityByDir(dir string) (*models.AnimeDBEntity, error) {
	if CacheMode {
		name, ok := m.dir2name[dir]
		if ok {
			if anime, ok := m.cacheAnimeDBEntity[name]; ok {
				return anime, nil
			}
		}
		log.Debugf("未找到内存缓存，读取文件数据库: %s", dir)
	}
	file := path.Join(dir, AnimeDBName)
	if !utils.IsExist(file) {
		return nil, nil
	}
	a := &models.AnimeDBEntity{}
	err := m.read(file, a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// getAnimeDBEntity
//
//	获取文件数据库AnimeDBEntity
//	从内存或文件数据库中获取
func (m *Database) getAnimeDBEntity(name string) (*models.AnimeDBEntity, error) {
	if CacheMode {
		if anime, ok := m.cacheAnimeDBEntity[name]; ok {
			return anime, nil
		}
		log.Debugf("未找到 %s 内存缓存，读取文件数据库", name)
	}
	dir, ok := m.name2dir[name]
	if !ok {
		return nil, nil
	}
	return m.getAnimeDBEntityByDir(dir.Dir)
}

// setAnimeDBEntity
//
//	设置文件数据库AnimeDBEntity
//	写入内存和文件数据库
func (m *Database) setAnimeDBEntity(dir string, a *models.AnimeDBEntity) error {
	a.UpdateAt = utils.Unix()
	if a.CreateAt == 0 {
		a.CreateAt = a.UpdateAt
	}
	err := utils.CreateMutiDir(dir)
	if err != nil {
		return err
	}
	file := path.Join(dir, AnimeDBName)
	err = m.write(file, a)

	m.dir2name[dir] = a.Name
	m.name2dir[a.Name] = &AnimeDir{
		Dir:       dir,
		SeasonDir: make(map[int]string),
	}
	if CacheMode {
		m.cacheAnimeDBEntity[a.Name] = a
	}
	return nil
}

// getSeasonDBEntityByDir
//
//	获取文件数据库SeasonDBEntity，使用dir
//	从内存或文件数据库中获取
func (m *Database) getSeasonDBEntityByDir(dir string, season int) (*models.SeasonDBEntity, error) {
	if CacheMode {
		name, ok := m.dir2name[dir]
		if ok {
			if anime, ok := m.cacheSeasonDBEntity[name]; ok {
				if s, ok := anime[season]; ok {
					return s, nil
				}
			}
		}
		log.Debugf("未找到内存缓存，读取文件数据库: %s", dir)
	}

	file := path.Join(dir, SeasonDBName)
	if !utils.IsExist(file) {
		return nil, nil
	}
	s := &models.SeasonDBEntity{}
	err := m.read(file, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// getSeasonDBEntity
//
//	获取文件数据库SeasonDBEntity
//	从内存或文件数据库中获取
func (m *Database) getSeasonDBEntity(name string, season int) (*models.SeasonDBEntity, error) {
	if CacheMode {
		if anime, ok := m.cacheSeasonDBEntity[name]; ok {
			if s, ok := anime[season]; ok {
				return s, nil
			}
		}
		log.Debugf("未找到 %s 内存缓存，读取文件数据库", name)
	}
	dir, ok := m.name2dir[name]
	if !ok {
		return nil, nil
	}
	seasonDir, ok := dir.SeasonDir[season]
	if !ok {
		return nil, nil
	}
	return m.getSeasonDBEntityByDir(path.Join(dir.Dir, seasonDir), season)
}

// setSeasonDBEntity
//
//	设置文件数据库SeasonDBEntity
//	写入内存和文件数据库
func (m *Database) setSeasonDBEntity(dir string, s *models.SeasonDBEntity) error {
	s.UpdateAt = utils.Unix()
	if s.CreateAt == 0 {
		s.CreateAt = s.UpdateAt
	}
	err := utils.CreateMutiDir(dir)
	if err != nil {
		return err
	}
	file := path.Join(dir, SeasonDBName)
	err = m.write(file, s)
	if err != nil {
		return err
	}
	m.dir2name[dir] = s.Name
	if _, ok := m.name2dir[s.Name]; ok {
		m.name2dir[s.Name].SeasonDir[s.Season] = path.Base(dir)
	} else {
		m.name2dir[s.Name] = &AnimeDir{
			Dir:       path.Dir(dir),
			SeasonDir: make(map[int]string),
		}
		m.name2dir[s.Name].SeasonDir[s.Season] = dir
	}
	if CacheMode {
		if m.cacheSeasonDBEntity[s.Name] == nil {
			m.cacheSeasonDBEntity[s.Name] = make(map[int]*models.SeasonDBEntity)
		}
		m.cacheSeasonDBEntity[s.Name][s.Season] = s
	}
	return nil
}
