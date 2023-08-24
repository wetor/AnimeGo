package database

import (
	"fmt"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/dirdb"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

// Scan
//
//	扫描已完成下载
//	构建name2dir和dir2name缓存
//	如果启用CacheMode，同时会载入所有的文件数据库
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
			m.setAnimeCache(file.Dir, anime)
		case path.Ext(SeasonDBName):
			season := &models.SeasonDBEntity{}
			err = file.DB.Unmarshal(season)
			if err != nil {
				log.DebugErr(err)
				log.Warnf("读取数据文件失败: %s", file.File)
				break
			}
			m.setSeasonCache(file.Dir, season)
		case path.Ext(EpisodeDBFmt):
			ep := &models.EpisodeDBEntity{}
			err = file.DB.Unmarshal(ep)
			if err != nil {
				log.DebugErr(err)
				log.Warnf("读取数据文件失败: %s", file.File)
				break
			}
			m.setEpisodeCache(file.Dir, ep)
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

func (m *Database) setAnimeCache(dir string, anime *models.AnimeDBEntity) {
	m.dir2name[dir] = anime.Name
	if _, ok := m.name2dir[anime.Name]; !ok {
		m.name2dir[anime.Name] = &AnimeDir{
			Dir:       dir,
			SeasonDir: make(map[int]string),
		}
	}
	m.name2dir[anime.Name].Dir = dir
	m.cacheAnimeDBEntity[anime.Name] = anime
}

func (m *Database) setSeasonCache(dir string, season *models.SeasonDBEntity) {
	m.dir2name[dir] = season.Name
	if _, ok := m.name2dir[season.Name]; !ok {
		m.name2dir[season.Name] = &AnimeDir{
			Dir:       path.Dir(dir),
			SeasonDir: make(map[int]string),
		}
	}
	m.name2dir[season.Name].SeasonDir[season.Season] = dir
	if _, ok := m.cacheSeasonDBEntity[season.Name]; !ok {
		m.cacheSeasonDBEntity[season.Name] = make(map[int]*models.SeasonDBEntity)
	}
	m.cacheSeasonDBEntity[season.Name][season.Season] = season
}

func (m *Database) setEpisodeCache(dir string, ep *models.EpisodeDBEntity) {
	m.dir2name[dir] = ep.Name
	if _, ok := m.name2dir[ep.Name]; !ok {
		m.name2dir[ep.Name] = &AnimeDir{
			Dir:       path.Dir(dir), // 上层文件夹
			SeasonDir: make(map[int]string),
		}
	}
	m.name2dir[ep.Name].SeasonDir[ep.Season] = dir
	if _, ok := m.cacheDB[ep.Name]; !ok {
		m.cacheDB[ep.Name] = make(map[int]map[string]*models.EpisodeDBEntity)
	}
	if _, ok := m.cacheDB[ep.Name][ep.Season]; !ok {
		m.cacheDB[ep.Name][ep.Season] = make(map[string]*models.EpisodeDBEntity)
	}
	m.cacheDB[ep.Name][ep.Season][ep.Key()] = ep
}

// Add
//
//	添加数据到缓存中，根据类型决定缓存Bucket和Key
//	Step2
func (m *Database) Add(data any) error {
	m.Lock()
	defer m.Unlock()

	switch value := data.(type) {
	case *models.AnimeEntity:
		name := value.FullName()
		hash := value.Hash()
		m.cache.Put(Name2HashBucket, name, hash, 0)
		m.cache.Put(Hash2EntityBucket, hash, value, 0)
	}
	return nil
}

// GetAnimeEntity
//
//	获取AnimeEntity
//	从bolt中获取
func (m *Database) GetAnimeEntity(hash string) (*models.AnimeEntity, error) {
	anime := &models.AnimeEntity{}
	err := m.cache.Get(Hash2EntityBucket, hash, anime)
	if err != nil {
		return nil, err
	}
	return anime, nil
}

// GetAnimeEntityByName
//
//	获取AnimeEntity，使用name
//	从bolt中获取
func (m *Database) GetAnimeEntityByName(name string) (*models.AnimeEntity, error) {
	var hash string
	err := m.cache.Get(Name2HashBucket, name, &hash)
	if err != nil {
		return nil, err
	}
	return m.GetAnimeEntity(hash)
}

// getAnimeDBEntityByDir
//
//	获取文件数据库AnimeDBEntity，使用dir
//	从内存或文件数据库中获取
func (m *Database) getAnimeDBEntityByDir(dir string) (*models.AnimeDBEntity, error) {
	if _, ok := m.dir2name[dir]; !ok {
		return nil, &exceptions.ErrDatabaseDirNotFound{Dir: dir}
	}
	name := m.dir2name[dir]
	return m.getAnimeDBEntity(name)

}

// getAnimeDBEntity
//
//	获取文件数据库AnimeDBEntity
//	从内存或文件数据库中获取
func (m *Database) getAnimeDBEntity(name string) (*models.AnimeDBEntity, error) {
	if anime, ok := m.cacheAnimeDBEntity[name]; !ok {
		return nil, &exceptions.ErrDatabaseDBNotFound{Name: name}
	} else {
		return anime, nil
	}
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
	if err != nil {
		return err
	}
	m.setAnimeCache(dir, a)
	return nil
}

// getSeasonDBEntityByDir
//
//	获取文件数据库SeasonDBEntity，使用dir
//	从内存或文件数据库中获取
func (m *Database) getSeasonDBEntityByDir(dir string, season int) (*models.SeasonDBEntity, error) {
	if _, ok := m.dir2name[dir]; !ok {
		return nil, &exceptions.ErrDatabaseDirNotFound{Dir: dir}
	}
	name := m.dir2name[dir]
	return m.getSeasonDBEntity(name, season)
}

// getSeasonDBEntity
//
//	获取文件数据库SeasonDBEntity
//	从内存或文件数据库中获取
func (m *Database) getSeasonDBEntity(name string, season int) (*models.SeasonDBEntity, error) {
	if _, ok := m.cacheSeasonDBEntity[name]; !ok {
		return nil, &exceptions.ErrDatabaseDBNotFound{Name: name}
	}
	anime := m.cacheSeasonDBEntity[name]
	if s, ok := anime[season]; ok {
		return nil, &exceptions.ErrDatabaseDBNotFound{Name: name, Season: season}
	} else {
		return s, nil
	}
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
	m.setSeasonCache(dir, s)
	return nil
}

func (m *Database) getEpisodeDBEntity(name string, season int, ep int, epType models.AnimeEpType) (*models.EpisodeDBEntity, error) {
	if _, ok := m.cacheDB[name]; !ok {
		return nil, &exceptions.ErrDatabaseDBNotFound{Name: name}
	}
	if _, ok := m.cacheDB[name][season]; !ok {
		return nil, &exceptions.ErrDatabaseDBNotFound{Name: name, Season: season}
	}
	key := fmt.Sprintf("E%d-%v", ep, epType)
	if _, ok := m.cacheDB[name][season][key]; !ok {
		return nil, &exceptions.ErrDatabaseDBNotFound{Name: name, Season: season, Ep: ep}
	}
	return m.cacheDB[name][season][key], nil
}

func (m *Database) getEpisodeDBEntityList(name string, season int) ([]*models.EpisodeDBEntity, error) {
	if _, ok := m.cacheDB[name]; !ok {
		return nil, &exceptions.ErrDatabaseDBNotFound{Name: name}
	}
	eps := make([]*models.EpisodeDBEntity, 0)
	if season > 0 {
		if _, ok := m.cacheDB[name][season]; !ok {
			return nil, &exceptions.ErrDatabaseDBNotFound{Name: name, Season: season}
		}
		for _, e := range m.cacheDB[name][season] {
			eps = append(eps, e)
		}
	} else {
		for _, s := range m.cacheDB[name] {
			for _, e := range s {
				eps = append(eps, e)
			}
		}
	}
	return eps, nil
}

func (m *Database) setEpisodeDBEntity(filename string, ep *models.EpisodeDBEntity) error {
	ep.UpdateAt = utils.Unix()
	if ep.CreateAt == 0 {
		ep.CreateAt = ep.UpdateAt
	}
	dir := path.Dir(filename)
	err := utils.CreateMutiDir(dir)
	if err != nil {
		return err
	}
	file := fmt.Sprintf(EpisodeDBFmt, strings.TrimSuffix(filename, path.Ext(filename)))
	err = m.write(file, ep)
	if err != nil {
		return err
	}
	m.setEpisodeCache(dir, ep)
	return nil
}
