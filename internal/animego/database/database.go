package database

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"sync"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/exceptions"

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
	EpisodeDBFmt = "%s.e_json"
)

type AnimeDir struct {
	Dir       string
	SeasonDir map[int]string
}

type Database struct {
	cache         api.Cacher
	rename        api.Renamer
	name2dir      map[string]*AnimeDir // anime name -> anime dir&season dir
	dir2name      map[string]string    // anime dir/season dir -> anime name
	hash2filename map[string][]string  // hash -> dst filenames

	cacheAnimeDBEntity  map[string]*models.AnimeDBEntity
	cacheSeasonDBEntity map[string]map[int]*models.SeasonDBEntity

	cacheDB map[string]map[int]map[string]*models.EpisodeDBEntity
	sync.Mutex
	dirMutex sync.Mutex // 事务控制
}

func NewDatabase(cache api.Cacher, rename api.Renamer) (*Database, error) {
	m := &Database{
		cache:         cache,
		rename:        rename,
		name2dir:      make(map[string]*AnimeDir),
		dir2name:      make(map[string]string),
		hash2filename: make(map[string][]string),
	}
	m.cache.Add(Name2EntityBucket)
	m.cache.Add(Hash2NameBucket)

	m.cacheAnimeDBEntity = make(map[string]*models.AnimeDBEntity)
	m.cacheSeasonDBEntity = make(map[string]map[int]*models.SeasonDBEntity)
	m.cacheDB = make(map[string]map[int]map[string]*models.EpisodeDBEntity)

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
			log.Errorf("%+v", err)
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
			log.Errorf("%+v", err)
			continue
		}
	}
}

func (m *Database) OnDownloadComplete(events []models.ClientEvent) {
	log.Infof("OnDownloadComplete %v", events)
	for _, event := range events {
		err := m.handleDownloadComplete(event.Hash)
		if err != nil {
			log.Errorf("%+v", err)
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
				err = m.writeAllRenamed(_anime, _result)
				if err != nil {
					log.Warnf("写入文件数据库失败: %s", _name)
				}
				if _result.Scrape() {
					// TODO: 无法确保scrape成功
					if m.scrape(_anime, _result) {
						log.Infof("刮削完成: %s", _name)
						err = m.writeAllScraped(_anime, _result)
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
	if err != nil {
		log.DebugErr(err)
		log.Warnf("更新文件数据库失败")
		return errors.Wrapf(err, "处理事件失败: %s", event)
	}
	m.hash2filename[hash] = renameResult.Filenames()

	m.setAnimeCache(path.Join(Conf.SavePath, renameResult.AnimeDir), &models.AnimeDBEntity{
		BaseDBEntity: models.BaseDBEntity{
			Hash: hash,
			Name: name,
		},
	})
	m.setSeasonCache(path.Join(Conf.SavePath, renameResult.SeasonDir), &models.SeasonDBEntity{
		BaseDBEntity: models.BaseDBEntity{
			Hash: hash,
			Name: name,
		},
		Season: anime.Season,
	})
	// 是否启动重命名任务
	eps, err := m.getEpisodeDBEntityList(name, anime.Season)
	if err != nil {
		switch err.(type) {
		case *exceptions.ErrDatabaseDBNotFound:
			eps = make([]*models.EpisodeDBEntity, 0)
		default:
			return err
		}
	}
	enableEpsSet := make(map[string]struct{})
	for _, key := range epKeys {
		enableEpsSet[key] = struct{}{}
	}
	// 剔除已经重命名完成的ep
	for _, ep := range eps {
		key := anime.EpKeyByEp(ep.Ep)
		if _, ok := enableEpsSet[key]; ok && ep.Renamed {
			delete(enableEpsSet, key)
		}
	}
	// 重命名
	for key := range enableEpsSet {
		err = m.rename.EnableTask([]string{key})
		if err != nil {
			log.DebugErr(err)
			log.Warnf("启动重命名任务失败")
			return errors.Wrapf(err, "处理事件失败: %s", event)
		}
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
	// 处理Episode文件数据库
	if filenames, ok := m.hash2filename[hash]; ok {
		err = m.writeAllEpisode(anime, filenames, "seeded", true)
		if err != nil {
			return err
		}
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
	if filenames, ok := m.hash2filename[hash]; ok {
		err = m.writeAllEpisode(anime, filenames, "downloaded", true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Database) writeEpisode(anime *models.AnimeEntity, epIndex int, filename string, field string, value bool) error {
	name := anime.AnimeName()
	// 处理Episode文件数据库
	edit := false
	ep, err := m.getEpisodeDBEntity(name, anime.Season, anime.Ep[epIndex].Ep, anime.Ep[epIndex].Type)
	if err != nil {
		switch err.(type) {
		case *exceptions.ErrDatabaseDBNotFound:
			edit = true
			ep = &models.EpisodeDBEntity{
				BaseDBEntity: models.BaseDBEntity{
					Hash:     anime.Torrent.Hash,
					Name:     name,
					CreateAt: utils.Unix(),
				},
				StateDB: models.StateDB{},
				Season:  anime.Season,
				Type:    anime.Ep[epIndex].Type,
				Ep:      anime.Ep[epIndex].Ep,
			}
		default:
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
		ep.Hash = anime.Torrent.Hash
		err = m.setEpisodeDBEntity(path.Join(Conf.SavePath, filename), ep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Database) writeAllEpisode(anime *models.AnimeEntity, filenames []string, field string, value bool) error {
	// 处理Episode文件数据库
	for i, filename := range filenames {
		err := m.writeEpisode(anime, i, filename, field, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// writeAllRenamed
//
//	重命名完成，更新数据库
func (m *Database) writeAllRenamed(anime *models.AnimeEntity, renameResult *models.RenameAllResult) error {
	m.dirMutex.Lock()
	defer m.dirMutex.Unlock()
	name := anime.AnimeName()
	now := utils.Unix()
	dir := path.Join(Conf.SavePath, renameResult.AnimeDir)
	// 获取Anime文件数据库
	adb, err := m.getAnimeDBEntityByDir(dir)
	if err != nil {
		switch err.(type) {
		case *exceptions.ErrDatabaseDBNotFound:
			adb = &models.AnimeDBEntity{
				BaseDBEntity: models.BaseDBEntity{
					Hash: anime.Torrent.Hash,
					Name: name,
				},
			}
		default:
			return err
		}
	}
	// 写入Anime文件数据库
	err = m.setAnimeDBEntity(dir, adb)
	if err != nil {
		return err
	}

	// 获取Season文件数据库
	seasonDir := path.Join(Conf.SavePath, renameResult.SeasonDir)
	season, err := m.getSeasonDBEntityByDir(seasonDir, anime.Season)
	if err != nil {
		switch err.(type) {
		case *exceptions.ErrDatabaseDBNotFound:
			season = &models.SeasonDBEntity{
				BaseDBEntity: adb.BaseDBEntity,
				Season:       anime.Season,
			}
			season.CreateAt = now
		default:
			return err
		}
	}
	// 写入Season文件数据库
	err = m.setSeasonDBEntity(seasonDir, season)
	if err != nil {
		return err
	}

	delete(m.hash2filename, anime.Torrent.Hash)
	// 处理Episode文件数据库
	err = m.writeAllEpisode(anime, renameResult.Filenames(), "renamed", true)
	if err != nil {
		return err
	}
	return nil
}

// writeAllScraped
//
//	刮削完成，更新数据库
func (m *Database) writeAllScraped(anime *models.AnimeEntity, renameResult *models.RenameAllResult) error {
	m.dirMutex.Lock()
	defer m.dirMutex.Unlock()
	// 处理Episode文件数据库
	err := m.writeAllEpisode(anime, renameResult.Filenames(), "scraped", true)
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
	m.dirMutex.Lock()
	defer m.dirMutex.Unlock()
	m.Lock()
	defer m.Unlock()

	switch value := data.(type) {
	case *models.AnimeEntity:
		name := value.AnimeName()

		// 是否启动重命名任务
		eps, err := m.getEpisodeDBEntityList(name, value.Season)
		if err != nil {
			switch err.(type) {
			case *exceptions.ErrDatabaseDBNotFound:
				eps = make([]*models.EpisodeDBEntity, 0)
			default:
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
func (m *Database) scrape(anime *models.AnimeEntity, result *models.RenameAllResult) bool {
	if len(result.AnimeDir) == 0 {
		return true
	}
	nfo := path.Join(Conf.SavePath, result.AnimeDir, "tvshow.nfo")
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
