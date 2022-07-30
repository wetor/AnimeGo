package manager

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/modules/cache"
	"GoBangumi/modules/client"
	"GoBangumi/store"
	"fmt"
	"go.uber.org/zap"
	"path"
	"strings"
	"sync"
	"time"
)

const (
	UpdateWaitMinSecond = 2 // 允许的最短刷新时间
)

type Manager struct {
	client    client.Client
	bangumi   []*models.Bangumi               // 同步缓存，主要使用其中的Hash来索引item
	itemState map[string]*models.TorrentState // 存储当前项的状态信息，处理过的
	items     map[string]*models.TorrentItem  // 客户端下载项信息，直接获取到的

	downloadList []*models.Bangumi // 下载队列，存满或者盗下一个刷新时间会进行下载
	exitChan     chan bool         // 结束标记

	sync.Mutex
}

func NewManager(client client.Client) *Manager {
	m := &Manager{
		client:       client,
		downloadList: make([]*models.Bangumi, 0, config.Advanced().DownloadQueueMaxNum),
		exitChan:     make(chan bool),
	}
	// 首次运行将同步缓存与下载器下载项
	m.UpdateList()
	return m
}

// Download
//  @Description: 将下载任务加入到下载队列中
//  @Description: 如果队列满，会立即下载所有下载项
//  @receiver *Manager
//  @param bgm *models.Bangumi
//
func (m *Manager) Download(bgm *models.Bangumi) {
	m.Lock()
	defer m.Unlock()
	m.downloadList = append(m.downloadList, bgm)
	if len(m.downloadList) == config.Advanced().DownloadQueueMaxNum {
		list := make([]*models.Bangumi, len(m.downloadList))
		copy(list, m.downloadList)
		go m.download(list)
		m.downloadList = m.downloadList[0:0]
	}
}

// download
//  @Description: 批量下载队列，需要运行在协程中
//  @receiver *Manager
//  @param bgms []*models.Bangumi
//
func (m *Manager) download(bgms []*models.Bangumi) {
	conf := config.Setting()
	for _, bgm := range bgms {
		zap.S().Infof("开始下载「%s」", bgm.FullName())
		if !m.canDownload(bgm) {
			continue
		}
		m.client.Add(&models.ClientAddOptions{
			Urls:        []string{bgm.Url},
			SavePath:    conf.SavePath,
			Category:    conf.Category,
			Tag:         conf.Tag(bgm),
			SeedingTime: conf.SeedingTime,
			Rename:      bgm.FullName(),
		})
		// 通过gb下载的番剧，将存储与缓存中
		store.Cache.Put(cache.ClientBangumiBucket, bgm.Hash, bgm, 0)
		time.Sleep(time.Duration(config.Advanced().Main().DownloadQueueDelaySecond) * time.Second)
	}
}

// canDownload
//  @Description: 此资源能否下载
//  @Description: 如果hash已存在，则不会下载
//  @Description: 如果hash不存在，会判断bangumi ID、Season和ep，如果相同会判断是否允许重复下载
//  @receiver *Manager
//  @param bgm *models.Bangumi
//  @return bool
//
func (m *Manager) canDownload(bgm *models.Bangumi) bool {
	for _, b := range m.bangumi {
		if bgm.Hash == b.Hash {
			// 同一资源，不重复下载
			return false
		}
		// 同一集不同资源
		// 如果AllowDuplicateDownload == true，即允许同一集重复下载，则返回true，否则则不允许下载
		if bgm.ID == b.ID && bgm.Season == b.Season && bgm.Ep == b.Ep {
			return config.Advanced().Main().AllowDuplicateDownload
		}
	}
	return true
}

// Get
//  @Description: 更新种子下载状态
//  @receiver m
//
func (m *Manager) Get(hash string) *models.TorrentItem {
	//conf := config.Setting()
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
	minSize := config.Setting().IgnoreSizeMaxKb * 1024 // 单位 B
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
func (m *Manager) Start(exit chan bool) {
	go func() {
		for {
			select {
			case <-m.exitChan:
				exit <- true
				return
			default:
				if store.InitState == store.InitFinish {
					// 如果下载队列有内容则会进行下载
					if len(m.downloadList) > 0 {
						m.Lock()
						list := make([]*models.Bangumi, len(m.downloadList))
						copy(list, m.downloadList)
						go m.download(list)
						m.downloadList = m.downloadList[0:0]
						m.Unlock()
					}
					m.UpdateList()
				}
				delay := config.Advanced().Main().UpdateDelaySecond
				if delay < UpdateWaitMinSecond {
					delay = UpdateWaitMinSecond
				}
				time.Sleep(time.Duration(delay) * time.Second)
			}
		}
	}()
}

// Exit
//  @Description: 结束manager
//  @receiver *Manager
//
func (m *Manager) Exit() {
	m.exitChan <- true
}

// UpdateList
//  @Description: 遍历下载器列表，与缓存中的数据进行对比、合并
//  @receiver *Manager
//
func (m *Manager) UpdateList() {
	// 涉及到bangumi list的清空与重建，运行在协程中，需要加锁
	m.Lock()
	defer m.Unlock()
	conf := config.Setting()
	// 获取客户端下载列表
	items := m.client.List(&models.ClientListOptions{
		Category: conf.Category,
	})
	if items == nil {
		return
	}

	// 清空bangumi list，重新同步缓存
	m.bangumi = make([]*models.Bangumi, 0, len(items))
	// 内存list不清空
	if m.itemState == nil {
		m.itemState = make(map[string]*models.TorrentState, len(items))
	}
	if m.items == nil {
		m.items = make(map[string]*models.TorrentItem, len(items))
	}
	// 从缓存中读取对应的信息
	// 遍历下载器的下载列表（包括已完成）
	// 根据下载项hash在缓存中查找记录，如已存在则将信息重新加入到list中
	// 如不存在，则其不是通过gb下载的，忽略
	for _, item := range items {
		bangumiTemp := store.Cache.Get(cache.ClientBangumiBucket, item.Hash)
		if bangumiTemp != nil {
			if bangumi, ok := bangumiTemp.(*models.Bangumi); ok {
				m.bangumi = append(m.bangumi, bangumi)
				// 初始化itemState
				if _, has := m.itemState[item.Hash]; !has {
					m.itemState[item.Hash] = &models.TorrentState{
						Hash: item.Hash,
					}
				}
				state := m.itemState[item.Hash]
				state.State = item.State
				m.items[item.Hash] = item
				fmt.Println(bangumi.Name, item.Progress, item.State)

				// 是否已经重命名过，重启后第一次也会执行此部分，但不会做修改
				if !state.Renamed {
					// 获取相对路径，删除绝对路径前缀
					oldPath := strings.TrimPrefix(item.ContentPath, path.Clean(conf.SavePath)+"/")
					if oldPath == item.ContentPath {
						// 删除前缀失败，读取name
						if c := m.GetContent(item.Hash); c != nil {
							oldPath = c.Name
						}
					}
					zap.S().Infof("发现下载项「%s」", oldPath)
					newPath := path.Join(bangumi.DirName(), bangumi.FullName()+path.Ext(oldPath))
					if newPath != oldPath {
						m.client.Rename(&models.ClientRenameOptions{
							Hash:    item.Hash,
							OldPath: oldPath,
							NewPath: newPath,
						})
						zap.S().Infof("重命名「%s」->「%s」", oldPath, newPath)
					}
					state.Renamed = true
					state.Path = newPath
				}
			}
		}
	}
}
