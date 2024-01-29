package clientnotifier

import (
	"path"

	"github.com/google/wire"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/animego/database"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
)

var Set = wire.NewSet(
	NewNotifier,
)

type Notifier struct {
	rename api.Renamer

	hash2filename map[string][]string // hash -> dst filenames

	db *database.Database
	*models.NotifierOptions
}

func NewNotifier(opts *models.NotifierOptions, db *database.Database, rename api.Renamer) *Notifier {
	return &Notifier{
		hash2filename:   make(map[string][]string),
		rename:          rename,
		db:              db,
		NotifierOptions: opts,
	}
}

func (n *Notifier) Init() {
	n.hash2filename = make(map[string][]string)
	n.db.Init()
}

// OnDownloadStart 开始下载事件，重启后首次也会执行
//
//	Step 3
//	必须经过的流程
func (n *Notifier) OnDownloadStart(events []models.ClientEvent) {
	n.db.Lock()
	defer n.db.Unlock()
	log.Infof("OnDownloadStart %v", events)
	for _, event := range events {
		err := n.handleDownloadStart(event.Hash)
		if err != nil {
			log.Errorf("%+v", err)
			continue
		}
	}
}

func (n *Notifier) OnDownloadPause(events []models.ClientEvent) {
	log.Infof("OnDownloadPause %v", events)
}

func (n *Notifier) OnDownloadStop(events []models.ClientEvent) {
	log.Infof("OnDownloadStop %v", events)
}

// OnDownloadSeeding 做种事件，重启后首次也会执行
//
//	Step 4
//	必须经过的流程
func (n *Notifier) OnDownloadSeeding(events []models.ClientEvent) {
	log.Infof("OnDownloadSeeding %v", events)
	for _, event := range events {
		err := n.handleDownloadSeeding(event.Hash)
		if err != nil {
			log.Errorf("%+v", err)
			continue
		}
	}
}

// OnDownloadComplete 完成事件，重启后首次也会执行
//
//	Step 5
//	必须经过的流程
func (n *Notifier) OnDownloadComplete(events []models.ClientEvent) {
	log.Infof("OnDownloadComplete %v", events)
	for _, event := range events {
		err := n.handleDownloadComplete(event.Hash)
		if err != nil {
			log.Errorf("%+v", err)
			continue
		}
	}
}

func (n *Notifier) OnDownloadError(events []models.ClientEvent) {
	log.Infof("OnDownloadError %v", events)
}

func (n *Notifier) handleDownloadStart(hash string) error {
	event := "OnDownloadStart"
	// 获取缓存中的anime信息
	anime, err := n.db.GetAnimeEntity(hash)
	if err != nil {
		return errors.Wrapf(err, "处理事件失败: %s", event)
	}
	epKeys := anime.EpKeys()
	// 添加重命名任务
	var renameResult *models.RenameAllResult
	if !n.rename.HasRenameTask(epKeys) {
		renameResult, err = n.rename.AddRenameTask(&models.RenameOptions{
			Name:           anime.FullName(),
			Entity:         anime,
			SrcDir:         n.DownloadPath,
			DstDir:         n.SavePath,
			Mode:           n.Rename,
			RenameCallback: func(opts *models.RenameResult) {},
			CompleteCallback: func(_result *models.RenameAllResult) {
				_name := _result.Name
				// 写入文件夹数据库
				_anime, err := n.db.GetAnimeEntityByName(_name)
				if err != nil {
					log.Warnf("获取信息失败: %s", _name)
				}
				log.Infof("移动完成「%s」", _name)
				err = n.db.WriteAllRenamed(_anime, _result)
				if err != nil {
					log.Warnf("写入文件数据库失败: %s", _name)
				}
				delete(n.hash2filename, anime.Hash())
				if _result.Scrape() {
					// TODO: 无法确保scrape成功
					if n.db.Scrape(_anime, _result) {
						log.Infof("刮削完成: %s", _name)
						err = n.db.WriteAllScraped(_anime, _result)
						if err != nil {
							log.Warnf("写入文件数据库失败: %s", _name)
						}
					} else {
						log.Warnf("刮削失败: %s", _name)
					}
				}
				err = n.Callback.Func(_anime.Hash())
				if err != nil {
					log.Warnf("删除下载项失败: %s", _name)
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
	name := anime.AnimeName()

	n.hash2filename[hash] = renameResult.Filenames()

	n.db.SetAnimeCache(path.Join(n.SavePath, renameResult.AnimeDir), &models.AnimeDBEntity{
		BaseDBEntity: models.BaseDBEntity{
			Hash: hash,
			Name: name,
		},
	})
	n.db.SetSeasonCache(path.Join(n.SavePath, renameResult.SeasonDir), &models.SeasonDBEntity{
		BaseDBEntity: models.BaseDBEntity{
			Hash: hash,
			Name: name,
		},
		Season: anime.Season,
	})
	// 是否启动重命名任务
	eps, err := n.db.GetEpisodeDBEntityList(name, anime.Season)
	if err != nil {
		if exceptions.IsNotFound(err) {
			eps = make([]*models.EpisodeDBEntity, 0)
		} else {
			return err
		}
	}
	enableEpsSet := make(map[string]int)
	for i, key := range epKeys {
		enableEpsSet[key] = i
	}
	// 剔除已经重命名完成的ep
	for _, ep := range eps {
		key := anime.EpKeyByEp(ep.Ep)
		if idx, ok := enableEpsSet[key]; ok && ep.Renamed {
			log.Infof("发现部分已下载，跳过此部分重命名: %v", path.Join(n.DownloadPath, anime.Ep[idx].Src))
			delete(enableEpsSet, key)
		}
	}
	// 重命名
	for key := range enableEpsSet {
		err = n.rename.EnableTask([]string{key})
		if err != nil {
			log.DebugErr(err)
			log.Warnf("启动重命名任务失败")
			return errors.Wrapf(err, "处理事件失败: %s", event)
		}
	}
	return nil
}

func (n *Notifier) handleDownloadSeeding(hash string) error {
	event := "OnDownloadSeeding"
	// 获取缓存中的anime信息
	anime, err := n.db.GetAnimeEntity(hash)
	if err != nil {
		return errors.Wrapf(err, "处理事件失败: %s", event)
	}
	err = n.rename.SetDownloadState(anime.EpKeys(), constant.StateSeeding)
	if err != nil {
		return err
	}
	// 处理Episode文件数据库
	if filenames, ok := n.hash2filename[hash]; ok {
		err = n.db.WriteAllEpisode(anime, filenames, "seeded", true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *Notifier) handleDownloadComplete(hash string) error {
	event := "OnDownloadComplete"
	// 获取缓存中的anime信息
	anime, err := n.db.GetAnimeEntity(hash)
	if err != nil {
		return errors.Wrapf(err, "处理事件失败: %s", event)
	}
	err = n.rename.SetDownloadState(anime.EpKeys(), constant.StateComplete)
	if err != nil {
		return err
	}
	if filenames, ok := n.hash2filename[hash]; ok {
		err = n.db.WriteAllEpisode(anime, filenames, "downloaded", true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *Notifier) IsExist(data any) bool {
	return n.db.IsExist(data)
}

func (n *Notifier) Add(data any) error {
	return n.db.Add(data)
}

func (n *Notifier) Delete(data any) error {
	return n.db.Delete(data)
}

func (n *Notifier) Scan() error {
	return n.db.Scan()
}
