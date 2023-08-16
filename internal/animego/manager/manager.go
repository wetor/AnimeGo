package manager

import (
	"context"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/client"
	"github.com/wetor/AnimeGo/pkg/exceptions"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	DownloadChanDefaultCap = 10 // 下载通道默认容量
	NotFoundExpireHour     = 24
	AddingExpireSecond     = 60 // 添加状态超时
	Name2EntityBucket      = "name2entity"
	Name2StatusBucket      = "name2status"
	Hash2NameBucket        = "hash2name"
	SleepUpdateMaxCount    = 10

	FileAllExist   = 0
	FileSomeExist  = 1
	FileAllNoExist = 2
)

type Manager struct {
	client api.Client
	cache  api.Cacher
	rename api.Renamer

	// 通过管道传递下载项
	downloadChan chan any
	name2status  map[string]*models.DownloadStatus // 同时存在于db和内存中

	sleepUpdateCount int // UpdateList 休眠倒计数，当不存在正在下载、做种以及下载完成的项目时，在 SleepUpdateMaxCount 后停止更新

	errs     []error
	errMutex sync.Mutex
	sync.Mutex
}

// NewManager
//
//	@Description: 初始化下载管理器
//	@param client api.Downloader 下载客户端
//	@param cache api.Cacher 缓存
//	@param rename api.Renamer 重命名
//	@return *Manager
func NewManager(client api.Client, cache api.Cacher, rename api.Renamer) *Manager {
	m := &Manager{
		client:           client,
		cache:            cache,
		rename:           rename,
		name2status:      make(map[string]*models.DownloadStatus),
		sleepUpdateCount: SleepUpdateMaxCount,
		errs:             make([]error, 0),
	}
	m.downloadChan = make(chan any, DownloadChanDefaultCap)

	m.cache.Add(Name2EntityBucket)
	m.cache.Add(Name2StatusBucket)
	m.cache.Add(Hash2NameBucket)

	m.loadCache()

	m.UpdateList()
	return m
}

func (m *Manager) addError(err error) {
	m.errMutex.Lock()
	m.errs = append(m.errs, err)
	m.errMutex.Unlock()
}

func (m *Manager) loadCache() {
	// 同步name2status
	keyType := ""
	valueType := models.DownloadStatus{}
	m.cache.GetAll(Name2StatusBucket, &keyType, &valueType, func(k, v interface{}) {
		nv := &models.DownloadStatus{}
		_ = copier.Copy(nv, v.(*models.DownloadStatus))
		m.name2status[*k.(*string)] = nv
	})

	// 已下载项目，至少从头更新一遍
	for _, status := range m.name2status {
		status.Init = false
	}
}

func (m *Manager) Delete(hash []string) {
	err := m.client.Delete(&client.DeleteOptions{
		Hash:       hash,
		DeleteFile: true,
	})
	if err != nil {
		m.addError(err)
	}
}

// Download
//
//	@Description: 将下载任务加入到下载队列中
//	@Description: 如果队列满，调用此方法会阻塞
//	@receiver *Manager
//	@param anime any
func (m *Manager) Download(anime any) {
	m.downloadChan <- anime
}

func (m *Manager) download(anime *models.AnimeEntity) {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	m.Lock()
	defer m.Unlock()
	name := anime.FullName()

	if status, has := m.name2status[name]; has {
		if !m.isCompleted(status) || m.fileExist(Conf.SavePath, status.Path) == FileAllExist {
			if status.Init {
				if status.Renamed {
					log.Infof("发现已下载「%v」", strings.Join(status.Path, ", "))
				} else {
					log.Infof("发现正在下载「%s」", name)
				}
			}
			if !Conf.AllowDuplicateDownload {
				log.Infof("取消下载，不允许重复「%s」", name)
				return
			}
		}
	}
	log.Infof("开始下载「%s」", name)
	err := m.client.Add(&client.AddOptions{
		Url:         anime.Torrent.Url,
		File:        anime.Torrent.File,
		SavePath:    m.client.Config().DownloadPath,
		Category:    Conf.Category,
		Tag:         utils.Tag(Conf.Tag, anime.AirDate, anime.Ep[0].Ep),
		SeedingTime: Conf.SeedingTimeMinute,
		Rename:      name,
	})
	if err != nil {
		m.addError(err)
		return
	}
	m.cache.Put(Hash2NameBucket, anime.Torrent.Hash, name, 0)
	m.cache.Put(Name2EntityBucket, name, anime, 0)

	status := &models.DownloadStatus{
		Hash:     anime.Torrent.Hash,
		State:    downloader.StateAdding,
		ExpireAt: utils.Unix() + AddingExpireSecond,
	}
	m.name2status[name] = status
	m.cache.Put(Name2StatusBucket, name, status, 0)
}

// Start
//
//	@Description: 下载管理器主循环
//	@receiver *Manager
//	@param ctx context.Context
func (m *Manager) Start(ctx context.Context) {
	WG.Add(1)
	// 刷新信息、接收下载、接收退出指令协程
	go func() {
		defer WG.Done()
		for {
			exit := false
			func() {
				defer utils.HandleError(func(err error) {
					log.Errorf("%+v", err)
					m.sleep(ctx)
				})
				defer func() {
					m.errMutex.Lock()
					defer m.errMutex.Unlock()
					for _, err := range m.errs {
						log.Errorf("%s", err)
					}
					m.errs = make([]error, 0)
				}()
				select {
				case <-ctx.Done():
					log.Debugf("正常退出 manager cache")
					exit = true
					return
				case download := <-m.downloadChan:
					anime := download.(*models.AnimeEntity)
					if m.client.Connected() {
						log.Debugf("接收到下载项:「%s」", anime.FullName())
						m.download(anime)
					} else {
						go func() {
							m.downloadChan <- anime
							log.Warnf("等待连接到下载器。已接收到%d个下载项", len(m.downloadChan))
						}()
						m.sleep(ctx)
					}
				default:
					m.UpdateList()
					m.sleep(ctx)
				}
			}()
			if exit {
				return
			}
		}
	}()
}

func (m *Manager) sleep(ctx context.Context) {
	utils.Sleep(Conf.UpdateDelaySecond, ctx)
}

func (m *Manager) updateStatus(status *models.DownloadStatus, anime *models.AnimeEntity) {
	if anime == nil {
		return
	}
	name := anime.FullName()
	keys := anime.EpKeys()
	if status.Path == nil || len(status.Path) == 0 {
		status.Path = make([]string, len(anime.Ep))
	}
	if !status.Init {
		status.Seeded = false
		status.Downloaded = false
		status.Init = true
	}

	if !status.Renamed && len(anime.Ep) > 0 {
		if _, err := m.rename.GetRenameTaskState(keys); err != nil {
			renameOpt := &models.RenameOptions{
				Name:   name,
				Entity: anime,
				SrcDir: Conf.DownloadPath,
				DstDir: Conf.SavePath,
				Mode:   Conf.Rename,
				RenameCallback: func(opts *models.RenameResult) {
					status.Path[opts.Index] = opts.Filename
					// TODO: 无法确保scrape成功
				},
				CompleteCallback: func(opts *models.RenameAllResult) {
					status.Renamed = true
					status.Scraped = m.scrape(anime, opts.AnimeDir)
					log.Infof("移动完成「%s」", opts.Name)
				},
			}
			_, err := m.rename.AddRenameTask(renameOpt)
			if err != nil {
				m.addError(err)
				return
			}
			err = m.rename.EnableTask(keys)
			if err != nil {
				m.addError(err)
				return
			}
		}
	}

	// 做种，或未下载完成，但State符合下载完成状态
	if !status.Seeded {
		if status.State == downloader.StateSeeding || status.State == downloader.StateComplete {
			if !status.Renamed {
				go func() {
					err := m.rename.SetDownloadState(keys, downloader.StateSeeding)
					if err != nil {
						m.addError(err)
					}
				}()
			}
			status.Seeded = true
		}
	}

	// 未下载完成，但State符合下载完成状态
	if !status.Downloaded {
		// 完成下载
		if status.State == downloader.StateComplete {
			if !status.Renamed {
				go func() {
					err := m.rename.SetDownloadState(keys, downloader.StateComplete)
					if err != nil {
						m.addError(err)
					}
				}()
			}
			status.Downloaded = true
		}
	}
}

func (m *Manager) DeleteCache(fullname string) {
	lock := m.TryLock()
	defer func() {
		if lock {
			m.Unlock()
		}
	}()
	delete(m.name2status, fullname)
	err := m.cache.Delete(Name2StatusBucket, fullname)
	if err != nil {
		m.addError(err)
	}
	err = m.cache.Delete(Name2EntityBucket, fullname)
	if err != nil {
		m.addError(err)
	}
}

func (m *Manager) isCompleted(status *models.DownloadStatus) bool {
	return (status.Init && status.Seeded && status.Downloaded && status.Renamed && status.Scraped) ||
		(status.State == downloader.StateNotFound)
}

func (m *Manager) setDownloaded(status *models.DownloadStatus) {
	status.Init = true
	status.Seeded = true
	status.Downloaded = true
	status.Renamed = true
	status.Scraped = true
	status.State = downloader.StateComplete
}

func (m *Manager) fileExist(dir string, ps []string) int {
	existNum := 0
	for _, p := range ps {
		if len(p) != 0 && utils.IsExist(path.Join(dir, p)) {
			existNum++
		}
	}
	if len(ps) == existNum {
		return FileAllExist
	} else if existNum > 0 {
		return FileSomeExist
	} else {
		return FileAllNoExist
	}
}

func (m *Manager) UpdateList() {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	m.Lock()
	defer m.Unlock()

	// 获取客户端下载列表
	items, err := m.client.List(&client.ListOptions{
		Category: Conf.Category,
	})
	if err != nil {
		m.addError(err)
		return
	}
	hash2item := make(map[string]*client.TorrentItem)
	for _, item := range items {
		hash2item[item.Hash] = item
		if state := downloader.StateMap(item.State); state != downloader.StateComplete &&
			state != downloader.StateError && state != downloader.StateUnknown {
			m.sleepUpdateCount = SleepUpdateMaxCount
		}
	}
	if m.sleepUpdateCount <= 0 {
		return
	} else {
		m.sleepUpdateCount--
	}

	deleteNames := make([]string, 0)
	updateStatusKeys := make([]any, 0, len(m.name2status))
	updateStatusValues := make([]any, 0, len(m.name2status))

	now := utils.Unix()
	for name, status := range m.name2status {
		// adding状态的为新加入下载项，检查超时
		if status.State == downloader.StateAdding {
			if status.Expire(now) {
				deleteNames = append(deleteNames, name)
			}
			continue
		}
		// 文件是否存在
		fileExist := m.fileExist(Conf.SavePath, status.Path) == FileAllExist
		isCompleted := m.isCompleted(status)
		item, hasItem := hash2item[status.Hash]
		anime := &models.AnimeEntity{}
		if !isCompleted {
			err := m.cache.Get(Name2EntityBucket, name, anime)
			if err != nil {
				anime = nil
			}
		}
		// 在下载列表中
		if hasItem {
			status.ExpireAt = 0
			status.State = downloader.StateMap(item.State)
			// 未完成
			if !isCompleted {
				// 同步下载列表
				m.updateStatus(status, anime)
			}
			// 完成
			if isCompleted {
				// 文件存在
				if fileExist {
					status.State = downloader.StateComplete
					// 重命名完成
					if status.Renamed {
						// 删除已完成的列表记录
						err = m.client.Delete(&client.DeleteOptions{
							Hash:       []string{status.Hash},
							DeleteFile: false,
						})
						if err != nil {
							m.addError(err)
							return
						}
					}
				}
				// 文件不存在
				if !fileExist {
					err := m.cache.Get(Name2EntityBucket, name, anime)
					if err != nil {
						anime = nil
					}
					// 原文件存在
					if m.fileExist(Conf.DownloadPath, anime.FilePathSrc()) == FileAllExist {
						status.Init = false
						status.Renamed = false
					}
					// 同步下载列表
					m.updateStatus(status, anime)
				}
			}
			// 下载中，打印日志
			if !status.Downloaded {
				log.Debugf("下载进度: %v, 名称: %v, qbt状态: %v, 状态: %v",
					fmt.Sprintf("%.1f", item.Progress*100),
					name,
					item.State,
					status.State,
				)
			}
			updateStatusKeys = append(updateStatusKeys, name)
			updateStatusValues = append(updateStatusValues, status)
		}
		// 不在下载列表中
		if !hasItem {
			// 文件存在
			if fileExist {
				// 未完成
				if !isCompleted {
					// 可能是在下载过程中，在下载器中被手动删除下载项，默认已下载完成
					log.Warnf("存在可能未下载完成的项目「%v」，检查后选择是否删除", strings.Join(status.Path, ", "))
				}
				// 设置已下载完成
				m.setDownloaded(status)
				updateStatusKeys = append(updateStatusKeys, name)
				updateStatusValues = append(updateStatusValues, status)
			}
			// 文件不存在
			if !fileExist {
				// 设置缓存将在 NotFoundExpireHour 小时后过期
				if status.ExpireAt == 0 {
					status.ExpireAt = utils.Unix() + NotFoundExpireHour*60*60
					status.State = downloader.StateNotFound
					updateStatusKeys = append(updateStatusKeys, name)
					updateStatusValues = append(updateStatusValues, status)
				}
				// 未完成
				if !isCompleted {
					// 不在下载列表，文件不存在，未完成。可能是没有移动
					// 同步下载列表
					m.updateStatus(status, anime)
					// 设置已下载完成
					m.setDownloaded(status)
					updateStatusKeys = append(updateStatusKeys, name)
					updateStatusValues = append(updateStatusValues, status)
				}
				// 检查过期
				if status.Expire(now) {
					deleteNames = append(deleteNames, name)
				}
			}
		}
	}
	// 一次性更新状态到缓存中
	if len(updateStatusKeys) > 0 {
		m.cache.BatchPut(Name2StatusBucket, updateStatusKeys, updateStatusValues, 0)
	}
	// 处理删除，将删除 name2status 中数据
	for _, name := range deleteNames {
		m.DeleteCache(name)
	}

	// 处理新增
	var name string
	appendStatusKeys := make([]any, 0, len(m.name2status))
	appendStatusValues := make([]any, 0, len(m.name2status))
	for _, item := range items {
		// 尝试从已下载中查找name
		err := m.cache.Get(Hash2NameBucket, item.Hash, &name)
		if err != nil {
			continue
		}
		// 判断是否已下载
		status, has := m.name2status[name]
		if !has || status.State == downloader.StateAdding {
			// 未下载或状态为NotFound或Adding
			status = &models.DownloadStatus{
				Hash:     item.Hash,
				State:    downloader.StateMap(item.State),
				ExpireAt: 0,
			}
			m.name2status[name] = status
			appendStatusKeys = append(appendStatusKeys, name)
			appendStatusValues = append(appendStatusValues, status)
		}
	}
	// 一次性更新状态到缓存中
	if len(appendStatusKeys) > 0 {
		m.cache.BatchPut(Name2StatusBucket, appendStatusKeys, appendStatusValues, 0)
	}
}

func (m *Manager) scrape(bangumi *models.AnimeEntity, dir string) bool {
	if len(dir) == 0 {
		return true
	}
	nfo := path.Join(Conf.SavePath, dir, "tvshow.nfo")
	log.Infof("写入元数据文件「%s」", nfo)

	if !utils.IsExist(nfo) {
		err := os.WriteFile(nfo, []byte(bangumi.Meta()), os.ModePerm)
		if err != nil {
			err = errors.WithStack(&exceptions.ErrManager{Message: "写入tvshow.nfo元文件失败"})
			log.DebugErr(err)
			m.addError(err)
			return false
		}
	}
	data, err := os.ReadFile(nfo)
	if err != nil {
		err = errors.WithStack(&exceptions.ErrManager{Message: "打开tvshow.nfo元文件失败"})
		log.DebugErr(err)
		m.addError(err)
		return false
	}
	TmdbRegx := regexp.MustCompile(`<tmdbid>\d+</tmdbid>`)
	BangumiRegx := regexp.MustCompile(`<bangumiid>\d+</bangumiid>`)

	xmlStr := string(data)
	xmlStr = TmdbRegx.ReplaceAllString(xmlStr, fmt.Sprintf("<tmdbid>%d</tmdbid>", bangumi.ThemoviedbID))
	xmlStr = BangumiRegx.ReplaceAllString(xmlStr, fmt.Sprintf("<bangumiid>%d</bangumiid>", bangumi.ID))

	err = os.WriteFile(nfo, []byte(xmlStr), os.ModePerm)
	if err != nil {
		err = errors.WithStack(&exceptions.ErrManager{Message: "编辑tvshow.nfo元文件失败"})
		log.DebugErr(err)
		m.addError(err)
		return false
	}
	return true
}
