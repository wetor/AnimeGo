package downloader

import (
	"context"
	"fmt"
	"os"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/animego/downloader/qbittorrent"
	"github.com/wetor/AnimeGo/internal/animego/manager"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
)

const (
	UpdateWaitMinSecond    = 2  // 允许的最短刷新时间
	DownloadChanDefaultCap = 10 // 下载通道默认容量
	DownloadStateChan      = 5
	NotFoundExpireHour     = 3
	Name2EntityBucket      = "name2entity"
	Name2StatusBucket      = "name2status"
	Hash2NameBucket        = "hash2name"
	SleepUpdateMaxCount    = 10
)

type Manager struct {
	client downloader.Client
	cache  api.Cacher

	// 通过管道传递下载项
	downloadChan chan *models.AnimeEntity
	name2chan    map[string]chan models.TorrentState
	name2status  map[string]*models.DownloadStatus // 同时存在于db和内存中

	sleepUpdateCount int // UpdateList 休眠倒计数，当不存在正在下载、做种以及下载完成的项目时，在 SleepUpdateMaxCount 后停止更新
	sync.Mutex
}

// NewManager
//
//	@Description: 初始化下载管理器
//	@param client downloader.Client 下载客户端
//	@param cache cache.Cache 缓存
//	@param downloadChan chan *models.AnimeEntity 下载传递通道
//	@return *Manager
func NewManager(client downloader.Client, cache api.Cacher, downloadChan chan *models.AnimeEntity) *Manager {
	m := &Manager{
		client:           client,
		cache:            cache,
		name2chan:        make(map[string]chan models.TorrentState),
		name2status:      make(map[string]*models.DownloadStatus),
		sleepUpdateCount: SleepUpdateMaxCount,
	}

	if downloadChan == nil || cap(downloadChan) <= 1 {
		downloadChan = make(chan *models.AnimeEntity, DownloadChanDefaultCap)
	}
	m.downloadChan = downloadChan

	m.cache.Add(Name2EntityBucket)
	m.cache.Add(Name2StatusBucket)
	m.cache.Add(Hash2NameBucket)

	m.loadCache()

	m.UpdateList()
	return m
}

func (m *Manager) loadCache() {
	// 同步name2status
	keyType := ""
	valueType := models.DownloadStatus{}
	m.cache.GetAll(Name2StatusBucket, &keyType, &valueType, func(k, v interface{}) {
		nv := &models.DownloadStatus{}
		utils.ConvertModel(v.(*models.DownloadStatus), nv)
		m.name2status[*k.(*string)] = nv
	})

	for _, status := range m.name2status {
		status.Init = false
	}
}

func (m *Manager) Delete(hash []string) {
	m.client.Delete(&models.ClientDeleteOptions{
		Hash:       hash,
		DeleteFile: true,
	})
}

// Download
//
//	@Description: 将下载任务加入到下载队列中
//	@Description: 如果队列满，调用此方法会阻塞
//	@receiver *Manager
//	@param anime *models.AnimeEntity
func (m *Manager) Download(anime *models.AnimeEntity) {
	m.downloadChan <- anime
}

func (m *Manager) download(anime *models.AnimeEntity) {
	m.Lock()
	defer m.Unlock()
	name := anime.FullName()

	if status, has := m.name2status[name]; has {
		// 已有下载记录
		if status.State != StateNotFound {
			// 文件已存在
			if len(status.Path) != 0 && utils.IsExist(path.Join(manager.DownloaderConf.SavePath, status.Path)) {
				log.Infof("发现已下载「%s」", status.Path)
			} else if status.Init {
				log.Infof("发现正在下载「%s」", name)
			}
			if !manager.DownloaderConf.AllowDuplicateDownload {
				log.Infof("取消下载，不允许重复「%s」", name)
				return
			}
		}
	}
	log.Infof("开始下载「%s」", name)
	m.client.Add(&models.ClientAddOptions{
		Urls:        []string{anime.Url},
		SavePath:    manager.DownloaderConf.DownloadPath,
		Category:    manager.DownloaderConf.Category,
		Tag:         utils.Tag(manager.DownloaderConf.Tag, anime.AirDate, anime.Ep),
		SeedingTime: manager.DownloaderConf.SeedingTimeMinute,
		Rename:      name,
	})
	m.cache.Put(Hash2NameBucket, anime.Hash, name, 0)
	m.cache.Put(Name2EntityBucket, name, anime, 0)

	status := &models.DownloadStatus{
		Hash:  anime.Hash,
		State: StateAdding,
	}
	m.name2status[name] = status
	m.cache.Put(Name2StatusBucket, name, status, 0)
}

func (m *Manager) GetContent(opt *models.ClientGetOptions) *models.TorrentContentItem {
	cs := m.client.GetContent(opt)
	if len(cs) == 0 {
		return nil
	}
	maxSize := 0
	index := -1
	minSize := manager.DownloaderConf.IgnoreSizeMaxKb * 1024 // 单位 B
	for i, c := range cs {
		if c.Size < minSize {
			continue
		}
		if c.Size > maxSize {
			maxSize = c.Size
			index = i
		}
	}
	if index < 0 {
		return nil
	}
	// TODO: 支持多内容返回
	// TODO: 支持外挂字幕
	return cs[index]
}

// Start
//
//	@Description: 下载管理器主循环
//	@receiver *Manager
//	@param ctx context.Context
func (m *Manager) Start(ctx context.Context) {
	manager.WG.Add(1)
	// 刷新信息、接收下载、接收退出指令协程
	go func() {
		defer manager.WG.Done()
		for {
			exit := false
			func() {
				defer errors.HandleError(func(err error) {
					log.Errorf("", err)
					m.sleep(ctx)
				})
				select {
				case <-ctx.Done():
					log.Debugf("正常退出 manager downloader")
					exit = true
					return
				case anime := <-m.downloadChan:
					if m.client.Connected() {
						log.Debugf("接收到下载项:「%s」", anime.FullName())
						m.download(anime)
					} else {
						log.Warnf("无法连接客户端，等待。已接收到%d个下载项", len(m.downloadChan))
						go func() {
							m.downloadChan <- anime
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
	delay := manager.DownloaderConf.UpdateDelaySecond
	if delay < UpdateWaitMinSecond {
		delay = UpdateWaitMinSecond
	}
	utils.Sleep(delay, ctx)
}

func (m *Manager) UpdateDownloadItem(status *models.DownloadStatus, anime *models.AnimeEntity, item *models.TorrentItem) {
	status.State = stateMap(item.State)
	name := anime.FullName()

	if !status.Init {
		content := m.GetContent(&models.ClientGetOptions{
			Hash: status.Hash,
			Item: item,
		})

		renamePath := path.Join(anime.DirName(), anime.FileName()+path.Ext(content.Name))
		m.name2chan[name] = make(chan models.TorrentState, DownloadStateChan)
		renameOpt := &models.RenameOptions{
			Src:   path.Join(manager.DownloaderConf.DownloadPath, content.Name),
			Dst:   path.Join(manager.DownloaderConf.SavePath, renamePath),
			State: m.name2chan[name],
			RenameCallback: func() {
				status.Path = renamePath
				status.Scraped = m.scrape(anime)
			},
			Callback: func() {
				status.Renamed = true
				if c, ok := m.client.(*qbittorrent.QBittorrent); ok {
					// qbt需要手动删除列表记录，否则无法重复下载
					c.Delete(&models.ClientDeleteOptions{
						Hash:       []string{status.Hash},
						DeleteFile: false,
					})
				}
			},
		}
		RenameAnime(renameOpt)
		if status.State == StateSeeding || status.State == StateComplete {
			go func() {
				m.name2chan[name] <- status.State
			}()
		}
		status.Init = true
	}

	// 移动完成，且搜刮元数据失败
	if status.Renamed && !status.Scraped {
		status.Scraped = m.scrape(anime)
	}

	// 做种，或未下载完成，但State符合下载完成状态
	if !status.Seeded {
		if status.State == StateSeeding ||
			(status.State == StateWaiting && item.Progress == 1) {
			go func() {
				m.name2chan[name] <- status.State
			}()
			status.Seeded = true
		}
	}

	// 未下载完成，但State符合下载完成状态
	if !status.Downloaded {
		// 完成下载
		if status.State == StateComplete {
			go func() {
				m.name2chan[name] <- status.State
			}()
			status.Downloaded = true
		}
		log.Debugf("下载进度: %v, 名称: %v, qbt状态: %v, 状态: %v",
			fmt.Sprintf("%.1f", item.Progress*100),
			name,
			item.State,
			status.State,
		)
	}
}

func (m *Manager) DeleteCache(fullname string) {
	m.Lock()
	defer m.Unlock()

	delete(m.name2status, fullname)
	err := m.cache.Delete(Name2StatusBucket, fullname)
	errors.NewAniErrorD(err).TryPanic()
	err = m.cache.Delete(Name2EntityBucket, fullname)
	errors.NewAniErrorD(err).TryPanic()
}

func (m *Manager) UpdateList() {
	m.Lock()
	defer m.Unlock()

	// 获取客户端下载列表
	items := m.client.List(&models.ClientListOptions{
		Category: manager.DownloaderConf.Category,
	})
	hash2item := make(map[string]*models.TorrentItem)
	for _, item := range items {
		hash2item[item.Hash] = item
		if state := stateMap(item.State); state == StateDownloading || state == StateSeeding || state == StateComplete {
			m.sleepUpdateCount = SleepUpdateMaxCount
		}
	}
	if m.sleepUpdateCount <= 0 {
		return
	} else {
		m.sleepUpdateCount--
	}

	for name, status := range m.name2status {
		if status.State == StateAdding {
			continue
		}
		// 文件是否存在
		if len(status.Path) == 0 || utils.IsExist(path.Join(manager.DownloaderConf.SavePath, status.Path)) ||
			(!status.Init || !status.Renamed || !status.Scraped) {
			// 是否存在于下载列表
			if item, has := hash2item[status.Hash]; has {
				// 同步下载列表
				status.ExpireAt = 0
				anime := &models.AnimeEntity{}
				err := m.cache.Get(Name2EntityBucket, name, anime)
				if err == nil {
					m.UpdateDownloadItem(status, anime, item)
				}
			} else {
				// 不在下载列表中，标记完成
				status.State = StateComplete
			}
			m.cache.Put(Name2StatusBucket, name, status, 0)
		} else {
			// 文件不存在，检查过期时间
			if status.ExpireAt <= 0 {
				// 未设置过期，设置3小时过期
				status.ExpireAt = time.Now().Add(NotFoundExpireHour * time.Hour).Unix()
				status.State = StateNotFound
				m.cache.Put(Name2StatusBucket, name, status, 0)
			} else if status.ExpireAt-time.Now().Unix() <= 0 {
				// 已过期，删除
				m.Unlock()
				m.DeleteCache(name)
				m.Lock()
			}
		}
		delete(hash2item, status.Hash)
	}

	// 处理新增
	for _, item := range items {
		// 尝试从已下载中查找name
		name := ""
		err := m.cache.Get(Hash2NameBucket, item.Hash, &name)
		if err != nil {
			continue
		}
		// 判断是否已下载
		if status, has := m.name2status[name]; has {
			// 已下载
			if status.State != StateNotFound && status.State != StateAdding {
				// 文件存在，跳过下载
				continue
			}
		}
		status := &models.DownloadStatus{
			Hash:     item.Hash,
			State:    stateMap(item.State),
			ExpireAt: 0,
		}
		m.name2status[name] = status
		m.cache.Put(Name2StatusBucket, name, status, 0)
	}
}

func (m *Manager) scrape(bangumi *models.AnimeEntity) bool {
	nfo := path.Join(manager.DownloaderConf.SavePath, bangumi.DirName(), "tvshow.nfo")
	log.Infof("写入元数据文件「%s」", nfo)

	if !utils.IsExist(nfo) {
		err := os.WriteFile(nfo, []byte(bangumi.Meta()), os.ModePerm)
		if err != nil {
			log.Debugf("", errors.NewAniErrorD(err))
			log.Warnf("写入tvshow.nfo元文件失败")
			return false
		}
	}
	data, err := os.ReadFile(nfo)
	if err != nil {
		log.Debugf("", errors.NewAniErrorD(err))
		log.Warnf("打开已存在的tvshow.nfo元文件失败")
		return false
	}
	TmdbRegx := regexp.MustCompile(`<tmdbid>\d+</tmdbid>`)
	BangumiRegx := regexp.MustCompile(`<bangumiid>\d+</bangumiid>`)

	xmlStr := string(data)
	xmlStr = TmdbRegx.ReplaceAllString(xmlStr, fmt.Sprintf("<tmdbid>%d</tmdbid>", bangumi.ThemoviedbID))
	xmlStr = BangumiRegx.ReplaceAllString(xmlStr, fmt.Sprintf("<bangumiid>%d</bangumiid>", bangumi.ID))

	err = os.WriteFile(nfo, []byte(xmlStr), os.ModePerm)
	if err != nil {
		log.Debugf("", errors.NewAniErrorD(err))
		log.Warnf("写入修改的tvshow.nfo元文件失败")
		return false
	}
	return true
}
