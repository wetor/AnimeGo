// Package downloader
// @Description: 接收到下载项后，调用client进行下载，并对下载条目做一定处理
package downloader

import (
	"AnimeGo/internal/animego/downloader"
	"AnimeGo/internal/cache"
	"AnimeGo/internal/models"
	"AnimeGo/internal/store"
	"AnimeGo/internal/utils"
	"context"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"go.uber.org/zap"
)

const (
	UpdateWaitMinSecond    = 2  // 允许的最短刷新时间
	DownloadChanDefaultCap = 10 // 下载通道默认容量
)

var (
	Bucket = "manager"
)

type Manager struct {
	client    downloader.Client
	cache     cache.Cache
	bangumi   map[string]*models.AnimeEntity // 同步缓存，主要使用其中的Hash来索引item
	itemState map[string]*models.Torrent     // 存储当前项的状态信息，处理过的
	items     map[string]*models.TorrentItem // 客户端下载项信息，直接获取到的

	downloadQueue []*models.AnimeEntity // 下载队列，存满或者盗下一个刷新时间会进行下载

	// 通过管道传递下载项
	downloadChan chan *models.AnimeEntity

	sync.Mutex
}

// NewManager
//  @Description: 初始化下载管理器
//  @param client downloader.Client 下载客户端
//  @param cache cache.Cache 缓存
//  @param downloadChan chan *models.AnimeEntity 下载传递通道
//  @return *Manager
//
func NewManager(client downloader.Client, cache cache.Cache, downloadChan chan *models.AnimeEntity) *Manager {
	m := &Manager{
		client:        client,
		cache:         cache,
		downloadQueue: make([]*models.AnimeEntity, 0, store.Config.Advanced.MainConf.DownloadQueueMaxNum),
	}
	m.cache.Add(Bucket)
	if downloadChan == nil || cap(downloadChan) <= 1 {
		downloadChan = make(chan *models.AnimeEntity, DownloadChanDefaultCap)
	}
	m.downloadChan = downloadChan
	// 首次运行将同步缓存与下载器下载项
	m.UpdateList()
	return m
}

// Download
//  @Description: 将下载任务加入到下载队列中
//  @Description: 如果队列满，调用此方法会阻塞
//  @receiver *Manager
//  @param bangumi *models.AnimeEntity
//
func (m *Manager) Download(anime *models.AnimeEntity) {
	m.downloadChan <- anime
}

// download
//  @Description: 批量下载队列
//  @receiver *Manager
//  @param animes []*models.AnimeEntity
//  @param finish chan bool 添加下载完成后发送消息
//
func (m *Manager) download(animes []*models.AnimeEntity, ctx context.Context) {
	// 去重
	set := make(map[string]struct{})
	setIndex := make([]int, 0, len(animes))
	for i, anime := range animes {
		if _, has := set[anime.FullName()]; !has {
			set[anime.FullName()] = struct{}{}
			setIndex = append(setIndex, i)
		}
	}
	for _, i := range setIndex {
		anime := animes[i]
		zap.S().Infof("开始下载「%s」", anime.FullName())
		if !m.canDownload(anime) {
			zap.S().Debugf("取消下载，发现重复「%s」", anime.FullName())
			continue
		}

		m.bangumi[anime.Hash] = anime

		m.client.Add(&models.ClientAddOptions{
			Urls:        []string{anime.Url},
			SavePath:    store.Config.SavePath,
			Category:    store.Config.Category,
			Tag:         store.Config.Tag(anime),
			SeedingTime: store.Config.SeedingTime,
			Rename:      anime.FullName(),
		})
		// 通过gb下载的番剧，将存储与缓存中
		m.cache.Put(Bucket, anime.Hash, anime, 0)
		utils.Sleep(store.Config.Advanced.MainConf.DownloadQueueDelaySecond, ctx)
	}
}

// canDownload
//  @Description: 此资源能否下载
//  @Description: 如果hash已存在，则不会下载
//  @Description: 如果hash不存在，会判断bangumi ID、Season和ep，如果相同会判断是否允许重复下载
//  @receiver *Manager
//  @param bangumi *models.AnimeEntity
//  @return bool
//
func (m *Manager) canDownload(anime *models.AnimeEntity) bool {
	for _, b := range m.bangumi {
		if anime.Hash == b.Hash {
			// 同一资源，不重复下载
			return false
		}
		// 同一集不同资源
		// 如果AllowDuplicateDownload == true，即允许同一集重复下载，则返回true，否则则不允许下载
		if anime.ID == b.ID && anime.Season == b.Season && anime.Ep == b.Ep {
			return store.Config.Advanced.MainConf.AllowDuplicateDownload
		}
	}
	return true
}

// Get
//  @Description: 更新种子下载状态
//  @receiver m
//
func (m *Manager) Get(hash string) *models.TorrentItem {
	//conf := store.Config.Setting
	item := m.client.Get(&models.ClientGetOptions{
		Hash: hash,
	})
	return item
}

// GetContent
//  @Description: 返回最大的一个content
//  @receiver *Manager
//  @param hash string
//  @return *models.TorrentContentItem
//
func (m *Manager) GetContent(hash string) *models.TorrentContentItem {
	cs := m.client.GetContent(&models.ClientGetOptions{
		Hash: hash,
	})
	if len(cs) == 0 {
		return nil
	}
	maxSize := 0
	index := -1
	minSize := store.Config.IgnoreSizeMaxKb * 1024 // 单位 B
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
	return cs[index]
}

// Start
//  @Description: 下载管理器主循环
//  @receiver *Manager
//  @param exit chan bool 退出后的回调chan，manager结束后会返回true
//
func (m *Manager) Start(ctx context.Context) {
	store.WG.Add(1)
	// 开始下载协程
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.S().Error(err)
			}
		}()
		for {
			if len(m.downloadQueue) > 0 {
				m.Lock()
				list := make([]*models.AnimeEntity, len(m.downloadQueue))
				copy(list, m.downloadQueue)
				m.downloadQueue = m.downloadQueue[0:0]
				m.Unlock()
				m.download(list, ctx)
			} else {
				utils.Sleep(store.Config.Advanced.MainConf.DownloadQueueDelaySecond, ctx)
			}
		}
	}()
	// 刷新信息、接收下载、接收退出指令协程
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.S().Error(err)
			}
		}()
		defer store.WG.Done()
		for {
			select {
			case <-ctx.Done():
				zap.S().Debug("正常退出")
				return
			case anime := <-m.downloadChan:
				zap.S().Debugf("接收到下载项:「%s」", anime.FullName())
				m.Lock()
				m.downloadQueue = append(m.downloadQueue, anime)
				m.Unlock()
			default:
				m.UpdateList()

				delay := store.Config.Advanced.MainConf.UpdateDelaySecond
				if delay < UpdateWaitMinSecond {
					delay = UpdateWaitMinSecond
				}
				utils.Sleep(delay, ctx)
			}
		}
	}()
}

// UpdateList
//  @Description: 遍历下载器列表，与缓存中的数据进行对比、合并
//  @receiver *Manager
//
func (m *Manager) UpdateList() {
	// 涉及到bangumi list的清空与重建，运行在协程中，需要加锁
	m.Lock()
	defer m.Unlock()
	conf := store.Config.Setting
	// 获取客户端下载列表
	items := m.client.List(&models.ClientListOptions{
		Category: conf.Category,
	})
	if items == nil {
		return
	}

	// 清空bangumi map，重新同步缓存
	m.bangumi = make(map[string]*models.AnimeEntity, len(items))
	// 内存list不清空
	if m.itemState == nil {
		m.itemState = make(map[string]*models.Torrent, len(items))
	}
	if m.items == nil {
		m.items = make(map[string]*models.TorrentItem, len(items))
	}
	// 从缓存中读取对应的信息
	// 遍历下载器的下载列表（包括已完成）
	// 根据下载项hash在缓存中查找记录，如已存在则将信息重新加入到list中
	// 如不存在，则其不是通过gb下载的，忽略
	for _, item := range items {
		bangumi := &models.AnimeEntity{}
		err := m.cache.Get(Bucket, item.Hash, bangumi)
		if err == nil {
			if !m.canDownload(bangumi) {
				m.client.Delete(&models.ClientDeleteOptions{
					Hash: []string{item.Hash},
				})
				zap.S().Debugf("删除下载，发现重复「%s」", bangumi.FullName())
				continue
			}
			// item 缓存
			m.items[item.Hash] = item

			// 把从缓存数据库中存的bangumi信息加入到map中
			if _, has := m.bangumi[bangumi.Hash]; !has {
				m.bangumi[bangumi.Hash] = bangumi
			}

			// 初始化itemState
			if _, has := m.itemState[item.Hash]; !has {
				m.itemState[item.Hash] = &models.Torrent{Hash: item.Hash}
			}
			state := m.itemState[item.Hash]
			state.State = stateMap(item.State)

			// 未完成重命名，或者首次运行（如重启下载器）
			if !state.Renamed {
				// 获取相对路径，删除绝对路径前缀
				state.OldPath = strings.TrimPrefix(item.ContentPath, path.Clean(conf.SavePath)+"/")
				if state.OldPath == item.ContentPath {
					// 删除前缀失败，读取name
					if c := m.GetContent(item.Hash); c != nil {
						state.OldPath = c.Name
					}
				}
				zap.S().Infof("发现下载项「%s」", state.OldPath)
				state.Path = path.Join(bangumi.DirName(), bangumi.FileName()+path.Ext(state.OldPath))
				if state.Path != state.OldPath {
					m.client.Rename(&models.ClientRenameOptions{
						Hash:    item.Hash,
						OldPath: state.OldPath,
						NewPath: state.Path,
					})
					zap.S().Infof("重命名「%s」->「%s」", state.OldPath, state.Path)
				}
				state.Renamed = true
			}

			// 移动完成，但未搜刮元数据
			if state.Renamed && !state.Scraped {
				state.Scraped = m.scrape(bangumi)
				if state.Scraped {
					// TODO: 完成，是否删除下载项
				}
			}

			// 未下载完成，但State符合下载完成状态
			if !state.Downloaded {
				if state.State == StateComplete || state.State == StateSeeding ||
					(state.State == StateWaiting && item.Progress == 1) {
					state.Downloaded = true
				}
				zap.S().Debugw("下载进度",
					"名称", bangumi.FullName(),
					"进度", fmt.Sprintf("%.1f", item.Progress*100),
					"qbt状态", item.State,
					"状态", state.State,
				)
			}
		}
	}
}

func (m *Manager) scrape(bangumi *models.AnimeEntity) bool {
	nfo := path.Join(store.Config.SavePath, bangumi.DirName(), "tvshow.nfo")
	zap.S().Infof("写入元数据文件「%s」", nfo)

	_, err := os.Stat(nfo)
	if os.IsNotExist(err) {
		err = os.WriteFile(nfo, []byte(bangumi.Meta()), os.ModePerm)
		if err != nil {
			zap.S().Warn(err)
			return false
		}
	}
	data, err := os.ReadFile(nfo)
	if err != nil {
		zap.S().Warn(err)
		return false
	}
	TmdbRegx := regexp.MustCompile(`<tmdbid>\d+</tmdbid>`)
	BangumiRegx := regexp.MustCompile(`<bangumiid>\d+</bangumiid>`)

	xmlStr := string(data)
	xmlStr = TmdbRegx.ReplaceAllString(xmlStr, fmt.Sprintf("<tmdbid>%d</tmdbid>", bangumi.ThemoviedbID))
	xmlStr = BangumiRegx.ReplaceAllString(xmlStr, fmt.Sprintf("<bangumiid>%d</bangumiid>", bangumi.ID))

	err = os.WriteFile(nfo, []byte(xmlStr), os.ModePerm)
	if err != nil {
		zap.S().Warn(err)
		return false
	}
	return true
}
