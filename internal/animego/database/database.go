package database

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"sync"

	"github.com/wetor/AnimeGo/internal/animego/downloader"

	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	Name2EntityBucket = "name2entity"
	Hash2NameBucket   = "hash2name"

	AnimeDBName  = "anime.a_json"
	SeasonDBName = "anime.s_json"
)

type AnimeDir struct {
	Dir       string
	SeasonDir map[int]string
}

type Database struct {
	cache    api.Cacher
	rename   api.Renamer
	name2dir map[string]*AnimeDir
	dir2name map[string]string

	cacheAnimeDBEntity  map[string]*models.AnimeDBEntity
	cacheSeasonDBEntity map[string]map[int]*models.SeasonDBEntity

	sync.Mutex
	dirMutex sync.Mutex // 事务控制
}

func NewDatabase(cache api.Cacher, rename api.Renamer) (*Database, error) {
	m := &Database{
		cache:    cache,
		rename:   rename,
		name2dir: make(map[string]*AnimeDir),
		dir2name: make(map[string]string),
	}
	m.cache.Add(Name2EntityBucket)
	m.cache.Add(Hash2NameBucket)

	if CacheMode {
		m.cacheAnimeDBEntity = make(map[string]*models.AnimeDBEntity)
		m.cacheSeasonDBEntity = make(map[string]map[int]*models.SeasonDBEntity)
	}
	err := m.Scan()
	if err != nil {
		return nil, err
	}
	return m, nil
}

// OnDownloadStart 开始下载事件，重启后首次也会执行
//
//	Step 3
func (m *Database) OnDownloadStart(events []models.ClientEvent) {
	m.Lock()
	defer m.Unlock()
	log.Infof("OnDownloadStart %v", events)
	for _, event := range events {
		err := m.handleDownloadStart(event.Hash)
		if err != nil {
			//log.Errorf("%+v", err)
			continue
		}
	}

}

func (m *Database) OnDownloadPause(events []models.ClientEvent) {
	log.Infof("OnDownloadPause %v", events)
}

func (m *Database) OnDownloadStop(events []models.ClientEvent) {
	log.Infof("OnDownloadStop %v", events)
}

func (m *Database) OnDownloadSeeding(events []models.ClientEvent) {
	log.Infof("OnDownloadSeeding %v", events)
	for _, event := range events {
		err := m.handleDownloadSeeding(event.Hash)
		if err != nil {
			//log.Errorf("%+v", err)
			continue
		}
	}
}

func (m *Database) OnDownloadComplete(events []models.ClientEvent) {
	log.Infof("OnDownloadComplete %v", events)
	for _, event := range events {
		err := m.handleDownloadComplete(event.Hash)
		if err != nil {
			//log.Errorf("%+v", err)
			continue
		}
	}
}

func (m *Database) OnDownloadError(events []models.ClientEvent) {
	log.Infof("OnDownloadError %v", events)
}

func (m *Database) handleDownloadStart(hash string) error {
	event := "OnDownloadStart"
	// 获取缓存中的anime信息
	anime, err := m.getAnimeEntity(hash)
	if err != nil {
		return errors.Wrapf(err, "处理事件失败: %s", event)
	}
	name := anime.AnimeName()
	epKeys := anime.EpKeys()
	// 添加重命名任务
	var renameResult *models.RenameAllResult
	if !m.rename.HasRenameTask(epKeys) {
		renameResult, err = m.rename.AddRenameTask(&models.RenameOptions{
			Name:           name,
			Entity:         anime,
			SrcDir:         Conf.DownloadPath,
			DstDir:         Conf.SavePath,
			Mode:           Conf.Rename,
			RenameCallback: func(opts *models.RenameResult) {},
			CompleteCallback: func(_result *models.RenameAllResult) {
				_name := _result.Name
				log.Infof("移动完成「%s」", _name)
				// 写入文件夹数据库
				_anime, err := m.getAnimeEntityByName(_name)
				if err != nil {
					log.Warnf("获取信息失败: %s", _name)
				}
				err = m.writeRename(_anime, _result)
				if err != nil {
					log.Warnf("写入文件数据库失败: %s", _name)
				}
				if _result.Scrape() {
					// TODO: 无法确保scrape成功
					if m.scrape(_anime, _result.AnimeDir) {
						log.Infof("刮削完成: %s", _name)
						err = m.writeScrape(_result)
						if err != nil {
							log.Warnf("写入文件数据库失败: %s", _name)
						}
					} else {
						log.Warnf("刮削失败: %s", _name)
					}
				}
			},
		})
		if err != nil {
			log.DebugErr(err)
			log.Warnf("添加重命名任务失败")
			return errors.Wrapf(err, "处理事件失败: %s", event)
		}
	}

	animeDb, err := func() (*models.AnimeDBEntity, error) {
		m.dirMutex.Lock()
		defer m.dirMutex.Unlock()
		dir := path.Join(Conf.SavePath, renameResult.AnimeDir)
		// 读取文件夹数据库
		animeDb, err := m.getAnimeDBEntityByDir(dir)
		// 初始化文件数据库
		if animeDb == nil {
			animeDb = &models.AnimeDBEntity{
				BaseDBEntity: models.BaseDBEntity{
					Name: name,
					Hash: hash,
				},
				Init: true,
			}
			err = m.setAnimeDBEntity(dir, animeDb)
			if err != nil {
				return nil, err
			}
		}
		return animeDb, nil
	}()
	if err != nil {
		log.DebugErr(err)
		log.Warnf("更新文件数据库失败")
		return errors.Wrapf(err, "处理事件失败: %s", event)
	}
	// 是否启动重命名任务
	if !animeDb.Renamed {
		err = m.rename.EnableTask(epKeys)
		if err != nil {
			log.DebugErr(err)
			log.Warnf("启动重命名任务失败")
			return errors.Wrapf(err, "处理事件失败: %s", event)
		}
	} else {
		m.rename.DeleteTask(epKeys)
	}
	return nil
}

func (m *Database) handleDownloadSeeding(hash string) error {
	event := "OnDownloadSeeding"
	// 获取缓存中的anime信息
	anime, err := m.getAnimeEntity(hash)
	if err != nil {
		return errors.Wrapf(err, "处理事件失败: %s", event)
	}
	err = m.rename.SetDownloadState(anime.EpKeys(), downloader.StateSeeding)
	if err != nil {
		return err
	}
	err = m.writeField(anime.AnimeName(), "seeded", true)
	if err != nil {
		return err
	}
	return nil
}

func (m *Database) handleDownloadComplete(hash string) error {
	event := "OnDownloadComplete"
	// 获取缓存中的anime信息
	anime, err := m.getAnimeEntity(hash)
	if err != nil {
		return errors.Wrapf(err, "处理事件失败: %s", event)
	}
	err = m.rename.SetDownloadState(anime.EpKeys(), downloader.StateComplete)
	if err != nil {
		return err
	}
	err = m.writeField(anime.AnimeName(), "downloaded", true)
	if err != nil {
		return err
	}
	return nil
}

// writeRename
//
//	重命名完成，更新数据库
func (m *Database) writeRename(anime *models.AnimeEntity, renameResult *models.RenameAllResult) error {
	m.dirMutex.Lock()
	defer m.dirMutex.Unlock()
	dir := path.Join(Conf.SavePath, renameResult.AnimeDir)
	// 获取当前文件数据库
	adb, err := m.getAnimeDBEntityByDir(dir)
	if err != nil {
		return err
	}
	// 修改内容
	adb.Renamed = true
	// 写入文件数据库
	err = m.setAnimeDBEntity(dir, adb)
	if err != nil {
		return err
	}

	// 重新构建Season文件数据库
	season := &models.SeasonDBEntity{
		BaseDBEntity: adb.BaseDBEntity,
		Season:       anime.Season,
		Episodes:     make([]models.EpisodeDBEntity, len(renameResult.Results)),
	}
	for i, res := range renameResult.Results {
		ep := anime.Ep[i]
		season.Episodes[i] = models.EpisodeDBEntity{
			File: path.Base(res.Filename),
			Type: ep.Type,
			Ep:   ep.Ep,
		}
	}
	// 写入Season文件数据库
	err = m.setSeasonDBEntity(path.Join(Conf.SavePath, renameResult.SeasonDir), season)
	if err != nil {
		return err
	}
	return nil
}

// writeScrape
//
//	刮削完成，更新数据库
func (m *Database) writeScrape(renameResult *models.RenameAllResult) error {
	m.dirMutex.Lock()
	defer m.dirMutex.Unlock()
	dir := path.Join(Conf.SavePath, renameResult.AnimeDir)
	// 获取当前文件数据库
	adb, err := m.getAnimeDBEntityByDir(dir)
	if err != nil {
		return err
	}
	// 修改内容
	adb.Scraped = true
	// 写入文件数据库
	err = m.setAnimeDBEntity(dir, adb)
	if err != nil {
		return err
	}
	return nil
}

func (m *Database) writeField(name string, field string, value bool) error {
	m.dirMutex.Lock()
	defer m.dirMutex.Unlock()
	animeDir, ok := m.name2dir[name]
	if !ok {
		return nil
	}
	dir := animeDir.Dir
	// 获取当前文件数据库
	adb, err := m.getAnimeDBEntityByDir(dir)
	if err != nil {
		return err
	}
	// 修改内容
	switch field {
	case "init":
		adb.Init = value
	case "renamed":
		adb.Renamed = value
	case "downloaded":
		adb.Downloaded = value
	case "seeded":
		adb.Seeded = value
	case "scraped":
		adb.Scraped = value
	}
	// 写入文件数据库
	err = m.setAnimeDBEntity(dir, adb)
	if err != nil {
		return err
	}
	return nil
}

// IsExist
//
//	遍历本地文件夹数据库，判断是否已下载
//	Step 1
func (m *Database) IsExist(data any) bool {
	m.Lock()
	defer m.Unlock()

	switch value := data.(type) {
	case string:
		a, err := m.getAnimeDBEntity(value)
		if err != nil {
			log.DebugErr(err)
		}
		if a == nil {
			return false
		}
		return a.Downloaded
	case *models.AnimeEntity:
		name := value.AnimeName()
		a, err := m.getAnimeDBEntity(name)
		if err != nil {
			log.DebugErr(err)
			return false
		}
		if a == nil {
			return false
		}
		s, err := m.getSeasonDBEntity(name, value.Season)
		if err != nil {
			log.DebugErr(err)
			return false
		}
		if s == nil {
			return false
		}
		// 待判断的ep数大于已下载的ep数
		if len(value.Ep) > len(s.Episodes) {
			return false
		}
		animeDir, ok := m.name2dir[name]
		if !ok {
			return false
		}
		seasonDir, ok := animeDir.SeasonDir[value.Season]
		if !ok {
			return false
		}
		sameEpSum := 0
		for _, ep := range value.Ep {
			for _, e := range s.Episodes {
				if e.Type == ep.Type && e.Ep == ep.Ep && utils.IsExist(path.Join(seasonDir, e.File)) {
					sameEpSum++
					break
				}
			}
		}
		// 没有全部被下载
		if sameEpSum < len(value.Ep) {
			return false
		}
		return true

	}
	return false
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
		name := value.AnimeName()
		m.cache.Put(Hash2NameBucket, value.Torrent.Hash, name, 0)
		m.cache.Put(Name2EntityBucket, name, value, 0)
	}
	return nil
}

// scrape
//
//	刮削
func (m *Database) scrape(anime *models.AnimeEntity, dir string) bool {
	if len(dir) == 0 {
		return true
	}
	nfo := path.Join(Conf.SavePath, dir, "tvshow.nfo")
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
