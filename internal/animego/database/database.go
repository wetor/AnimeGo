package database

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/google/wire"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/dirdb"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

var Set = wire.NewSet(
	NewDatabase,
)

type Database struct {
	cache    api.Cacher
	name2dir map[string]*models.AnimeDir // anime name -> anime dir&season dir
	dir2name map[string]string           // anime dir/season dir -> anime name

	cacheAnimeDBEntity  map[string]*models.AnimeDBEntity
	cacheSeasonDBEntity map[string]map[int]*models.SeasonDBEntity

	cacheDB map[string]map[int]map[string]*models.EpisodeDBEntity
	sync.Mutex
	dirMutex sync.Mutex // 事务控制

	*models.DatabaseOptions
}

func NewDatabase(opts *models.DatabaseOptions, cache api.Cacher) (*Database, error) {
	dirdb.Init(&dirdb.Options{
		DefaultExt: []string{path.Ext(constant.DatabaseAnimeDBName),
			path.Ext(constant.DatabaseSeasonDBName), path.Ext(constant.DatabaseEpisodeDBFmt)}, // anime, season
	})
	m := &Database{
		cache:           cache,
		DatabaseOptions: opts,
	}
	m.cache.Add(constant.DatabaseHash2EntityBucket)
	m.cache.Add(constant.DatabaseName2HashBucket)
	m.Init()
	err := m.Scan()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Database) Init() {
	m.Lock()
	defer m.Unlock()
	m.dirMutex.Lock()
	defer m.dirMutex.Unlock()

	m.name2dir = make(map[string]*models.AnimeDir)
	m.dir2name = make(map[string]string)

	m.cacheAnimeDBEntity = make(map[string]*models.AnimeDBEntity)
	m.cacheSeasonDBEntity = make(map[string]map[int]*models.SeasonDBEntity)
	m.cacheDB = make(map[string]map[int]map[string]*models.EpisodeDBEntity)
}

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
	d, err := dirdb.Open(m.SavePath)
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
		case path.Ext(constant.DatabaseAnimeDBName):
			anime := &models.AnimeDBEntity{}
			err = file.DB.Unmarshal(anime)
			if err != nil {
				log.DebugErr(err)
				log.Warnf("读取数据文件失败: %s", file.File)
				break
			}
			m.SetAnimeCache(file.Dir, anime)
		case path.Ext(constant.DatabaseSeasonDBName):
			season := &models.SeasonDBEntity{}
			err = file.DB.Unmarshal(season)
			if err != nil {
				log.DebugErr(err)
				log.Warnf("读取数据文件失败: %s", file.File)
				break
			}
			m.SetSeasonCache(file.Dir, season)
		case path.Ext(constant.DatabaseEpisodeDBFmt):
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

func (m *Database) SetAnimeCache(dir string, anime *models.AnimeDBEntity) {
	m.dir2name[dir] = anime.Name
	if _, ok := m.name2dir[anime.Name]; !ok {
		m.name2dir[anime.Name] = &models.AnimeDir{
			Dir:       dir,
			SeasonDir: make(map[int]string),
		}
	}
	m.name2dir[anime.Name].Dir = dir
	if a, ok := m.cacheAnimeDBEntity[anime.Name]; ok {
		anime.CreateAt = a.CreateAt
	}
	m.cacheAnimeDBEntity[anime.Name] = anime
}

func (m *Database) SetSeasonCache(dir string, season *models.SeasonDBEntity) {
	m.dir2name[dir] = season.Name
	if _, ok := m.name2dir[season.Name]; !ok {
		m.name2dir[season.Name] = &models.AnimeDir{
			Dir:       path.Dir(dir),
			SeasonDir: make(map[int]string),
		}
	}
	m.name2dir[season.Name].SeasonDir[season.Season] = dir
	if _, ok := m.cacheSeasonDBEntity[season.Name]; !ok {
		m.cacheSeasonDBEntity[season.Name] = make(map[int]*models.SeasonDBEntity)
	}
	if s, ok := m.cacheSeasonDBEntity[season.Name][season.Season]; ok {
		season.CreateAt = s.CreateAt
	}
	m.cacheSeasonDBEntity[season.Name][season.Season] = season
}

func (m *Database) setEpisodeCache(dir string, ep *models.EpisodeDBEntity) {
	m.dir2name[dir] = ep.Name
	if _, ok := m.name2dir[ep.Name]; !ok {
		m.name2dir[ep.Name] = &models.AnimeDir{
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
		m.cache.Put(constant.DatabaseName2HashBucket, name, hash, 0)
		m.cache.Put(constant.DatabaseHash2EntityBucket, hash, value, 0)
	}
	return nil
}

// Delete
//
//	删除缓存
func (m *Database) Delete(data any) error {
	m.Lock()
	defer m.Unlock()

	switch value := data.(type) {
	case string:
		err := m.cache.Delete(constant.DatabaseName2HashBucket, value)
		if err != nil {
			return err
		}
		err = m.cache.Delete(constant.DatabaseHash2EntityBucket, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAnimeEntity
//
//	获取AnimeEntity
//	从bolt中获取
func (m *Database) GetAnimeEntity(hash string) (*models.AnimeEntity, error) {
	anime := &models.AnimeEntity{}
	err := m.cache.Get(constant.DatabaseHash2EntityBucket, hash, anime)
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
	err := m.cache.Get(constant.DatabaseName2HashBucket, name, &hash)
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
	file := path.Join(dir, constant.DatabaseAnimeDBName)
	err = m.write(file, a)
	if err != nil {
		return err
	}
	m.SetAnimeCache(dir, a)
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
	file := path.Join(dir, constant.DatabaseSeasonDBName)
	err = m.write(file, s)
	if err != nil {
		return err
	}
	m.SetSeasonCache(dir, s)
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

func (m *Database) GetEpisodeDBEntityList(name string, season int) ([]*models.EpisodeDBEntity, error) {
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
	file := fmt.Sprintf(constant.DatabaseEpisodeDBFmt, strings.TrimSuffix(filename, path.Ext(filename)))
	err = m.write(file, ep)
	if err != nil {
		return err
	}
	m.setEpisodeCache(dir, ep)
	return nil
}

func (m *Database) WriteEpisode(anime *models.AnimeEntity, epIndex int, filename string, field string, value bool) error {
	name := anime.AnimeName()
	// 处理Episode文件数据库
	edit := false
	ep, err := m.getEpisodeDBEntity(name, anime.Season, anime.Ep[epIndex].Ep, anime.Ep[epIndex].Type)
	if err != nil {
		if exceptions.IsNotFound(err) {
			edit = true
			ep = &models.EpisodeDBEntity{
				BaseDBEntity: models.BaseDBEntity{
					Hash:     anime.Hash(),
					Name:     name,
					CreateAt: utils.Unix(),
				},
				StateDB: models.StateDB{},
				Season:  anime.Season,
				Type:    anime.Ep[epIndex].Type,
				Ep:      anime.Ep[epIndex].Ep,
			}
		} else {
			return err
		}
	}
	// 修改内容
	switch field {
	case "seeded":
		if ep.Seeded != value {
			edit = true
			ep.Seeded = value
		}
	case "downloaded":
		if ep.Downloaded != value {
			edit = true
			ep.Downloaded = value
		}
	case "renamed":
		if ep.Renamed != value {
			edit = true
			ep.Renamed = value
		}
	case "scraped":
		if ep.Scraped != value {
			edit = true
			ep.Scraped = value
		}
	}
	if edit {
		ep.Hash = anime.Hash()
		err = m.setEpisodeDBEntity(path.Join(m.SavePath, filename), ep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Database) WriteAllEpisode(anime *models.AnimeEntity, filenames []string, field string, value bool) error {
	// 处理Episode文件数据库
	for i, filename := range filenames {
		err := m.WriteEpisode(anime, i, filename, field, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteAllRenamed
//
//	重命名完成，更新数据库
func (m *Database) WriteAllRenamed(anime *models.AnimeEntity, renameResult *models.RenameAllResult) error {
	m.dirMutex.Lock()
	defer m.dirMutex.Unlock()
	name := anime.AnimeName()
	now := utils.Unix()
	dir := path.Join(m.SavePath, renameResult.AnimeDir)
	// 获取Anime文件数据库
	adb, err := m.getAnimeDBEntityByDir(dir)
	if err != nil {
		if exceptions.IsNotFound(err) {
			adb = &models.AnimeDBEntity{
				BaseDBEntity: models.BaseDBEntity{
					Hash: anime.Hash(),
					Name: name,
				},
			}
		} else {
			return err
		}
	}
	// 写入Anime文件数据库
	err = m.setAnimeDBEntity(dir, adb)
	if err != nil {
		return err
	}

	// 获取Season文件数据库
	seasonDir := path.Join(m.SavePath, renameResult.SeasonDir)
	season, err := m.getSeasonDBEntityByDir(seasonDir, anime.Season)
	if err != nil {
		if exceptions.IsNotFound(err) {
			season = &models.SeasonDBEntity{
				BaseDBEntity: adb.BaseDBEntity,
				Season:       anime.Season,
			}
			season.CreateAt = now
		} else {
			return err
		}
	}
	// 写入Season文件数据库
	err = m.setSeasonDBEntity(seasonDir, season)
	if err != nil {
		return err
	}

	// 处理Episode文件数据库
	err = m.WriteAllEpisode(anime, renameResult.Filenames(), "renamed", true)
	if err != nil {
		return err
	}
	return nil
}

// WriteAllScraped
//
//	刮削完成，更新数据库
func (m *Database) WriteAllScraped(anime *models.AnimeEntity, renameResult *models.RenameAllResult) error {
	m.dirMutex.Lock()
	defer m.dirMutex.Unlock()
	// 处理Episode文件数据库
	err := m.WriteAllEpisode(anime, renameResult.Filenames(), "scraped", true)
	if err != nil {
		return err
	}
	return nil
}

// Scrape
//
//	刮削
func (m *Database) Scrape(anime *models.AnimeEntity, result *models.RenameAllResult) bool {
	if len(result.AnimeDir) == 0 {
		return true
	}
	nfo := path.Join(m.SavePath, result.AnimeDir, "tvshow.nfo")
	log.Infof("写入元数据文件「%s」", nfo)

	if !utils.IsExist(nfo) {
		err := os.WriteFile(nfo, []byte(anime.Meta()), os.ModePerm)
		if err != nil {
			log.DebugErr(err)
			log.Warnf("写入元文件失败: tvshow.nfo")
			return false
		}
	}
	data, err := os.ReadFile(nfo)
	if err != nil {
		log.DebugErr(err)
		log.Warnf("打开元文件失败: tvshow.nfo")
		return false
	}
	TmdbRegx := regexp.MustCompile(`<tmdbid>\d+</tmdbid>`)
	BangumiRegx := regexp.MustCompile(`<bangumiid>\d+</bangumiid>`)

	xmlStr := string(data)
	xmlStr = TmdbRegx.ReplaceAllString(xmlStr, fmt.Sprintf("<tmdbid>%d</tmdbid>", anime.ThemoviedbID))
	xmlStr = BangumiRegx.ReplaceAllString(xmlStr, fmt.Sprintf("<bangumiid>%d</bangumiid>", anime.ID))

	err = os.WriteFile(nfo, []byte(xmlStr), os.ModePerm)
	if err != nil {
		log.DebugErr(err)
		log.Warnf("修改元文件失败: tvshow.nfo")
		return false
	}
	return true
}

// IsExist
//
//	遍历本地文件夹数据库，判断是否已下载
//	Step 1
func (m *Database) IsExist(data any) bool {
	m.dirMutex.Lock()
	defer m.dirMutex.Unlock()
	m.Lock()
	defer m.Unlock()

	switch value := data.(type) {
	case *models.AnimeEntity:
		name := value.AnimeName()

		// 是否启动重命名任务
		eps, err := m.GetEpisodeDBEntityList(name, value.Season)
		if err != nil {
			if exceptions.IsNotFound(err) {
				eps = make([]*models.EpisodeDBEntity, 0)
			} else {
				log.DebugErr(err)
				return false
			}
		}
		sum := 0
		for _, ep := range value.Ep {
			for _, e := range eps {
				if ep.Type == e.Type && ep.Ep == e.Ep && e.Downloaded {
					sum++
					break
				}
			}
		}
		// 全部都已存在且下载完成
		if sum == len(value.Ep) {
			return true
		}
	}
	return false
}
