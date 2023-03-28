package manager

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/jinzhu/copier"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/animego/downloader/qbittorrent"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const (
	DownloadChanDefaultCap = 10 // 下载通道默认容量
	DownloadStateChanCap   = 5
	NotFoundExpireHour     = 3
	Name2EntityBucket      = "name2entity"
	Name2StatusBucket      = "name2status"
	Hash2NameBucket        = "hash2name"
	SleepUpdateMaxCount    = 10
)

type Manager struct {
	client api.Downloader
	cache  api.Cacher
	rename api.Renamer

	// 通过管道传递下载项
	downloadChan chan any
	name2chan    map[string]chan models.TorrentState
	name2status  map[string]*models.DownloadStatus // 同时存在于db和内存中

	sleepUpdateCount int // UpdateList 休眠倒计数，当不存在正在下载、做种以及下载完成的项目时，在 SleepUpdateMaxCount 后停止更新
	sync.Mutex
}

// NewManager
//
//	@Description: 初始化下载管理器
//	@param client api.Downloader 下载客户端
//	@param cache api.Cacher 缓存
//	@param rename api.Renamer 重命名
//	@return *Manager
func NewManager(client api.Downloader, cache api.Cacher, rename api.Renamer) *Manager {
	m := &Manager{
		client:           client,
		cache:            cache,
		rename:           rename,
		name2chan:        make(map[string]chan models.TorrentState),
		name2status:      make(map[string]*models.DownloadStatus),
		sleepUpdateCount: SleepUpdateMaxCount,
	}
	m.downloadChan = make(chan any, DownloadChanDefaultCap)

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
		_ = copier.Copy(nv, v.(*models.DownloadStatus))
		m.name2status[*k.(*string)] = nv
	})

	// 已下载项目，至少从头更新一遍
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
		// 已有下载记录
		if status.State != downloader.StateNotFound {
			// 文件已存在
			if len(status.Path) != 0 && utils.IsExist(xpath.Join(Conf.SavePath, status.Path)) {
				log.Infof("发现已下载「%s」", status.Path)
			} else if status.Init {
				log.Infof("发现正在下载「%s」", name)
			}
			if !Conf.AllowDuplicateDownload {
				log.Infof("取消下载，不允许重复「%s」", name)
				return
			}
		}
	}
	log.Infof("开始下载「%s」", name)
	m.client.Add(&models.ClientAddOptions{
		Urls:        []string{anime.Url},
		SavePath:    Conf.DownloadPath,
		Category:    Conf.Category,
		Tag:         utils.Tag(Conf.Tag, anime.AirDate, anime.Ep),
		SeedingTime: Conf.SeedingTimeMinute,
		Rename:      name,
	})
	m.cache.Put(Hash2NameBucket, anime.Hash, name, 0)
	m.cache.Put(Name2EntityBucket, name, anime, 0)

	status := &models.DownloadStatus{
		Hash:  anime.Hash,
		State: downloader.StateAdding,
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
	minSize := Conf.IgnoreSizeMaxKb * 1024 // 单位 B
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
	WG.Add(1)
	// 刷新信息、接收下载、接收退出指令协程
	go func() {
		defer WG.Done()
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
				case download := <-m.downloadChan:
					anime := download.(*models.AnimeEntity)
					if m.client.Connected() {
						log.Debugf("接收到下载项:「%s」", anime.FullName())
						m.download(anime)
					} else {
						go func() {
							m.downloadChan <- anime
							log.Warnf("无法连接客户端，等待。已接收到%d个下载项", len(m.downloadChan))
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

func (m *Manager) UpdateDownloadItem(status *models.DownloadStatus, anime *models.AnimeEntity, item *models.TorrentItem) {
	status.State = downloader.StateMap(item.State)
	name := anime.FullName()

	if !status.Init {
		content := m.GetContent(&models.ClientGetOptions{
			Hash: status.Hash,
			Item: item,
		})

		m.name2chan[name] = make(chan models.TorrentState, DownloadStateChanCap)
		renameOpt := &models.RenameOptions{
			Src: xpath.Join(Conf.DownloadPath, content.Name),
			Dst: &models.RenameDst{
				Anime:    anime,
				Content:  content,
				SavePath: Conf.SavePath,
			},
			Mode:  Conf.Rename,
			State: m.name2chan[name],
			RenameCallback: func(opts *models.RenameResult) {
				status.Path = opts.Filepath
				status.Scraped = m.scrape(anime, opts.TVShowDir)
			},
			CompleteCallback: func() {
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
		m.rename.AddRenameTask(renameOpt)
		status.Seeded = false
		status.Downloaded = false

		status.Init = true
	}

	// 做种，或未下载完成，但State符合下载完成状态
	if !status.Seeded {
		if status.State == downloader.StateSeeding || status.State == downloader.StateComplete ||
			(status.State == downloader.StateWaiting && item.Progress == 1) {
			go func() {
				m.name2chan[name] <- downloader.StateSeeding
			}()
			status.Seeded = true
		}
	}

	// 未下载完成，但State符合下载完成状态
	if !status.Downloaded {
		// 完成下载
		if status.State == downloader.StateComplete {
			go func() {
				m.name2chan[name] <- downloader.StateComplete
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
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	m.Lock()
	defer m.Unlock()

	// 获取客户端下载列表
	items := m.client.List(&models.ClientListOptions{
		Category: Conf.Category,
	})
	hash2item := make(map[string]*models.TorrentItem)
	for _, item := range items {
		hash2item[item.Hash] = item
		if state := downloader.StateMap(item.State); state == downloader.StateDownloading || state == downloader.StateSeeding || state == downloader.StateComplete {
			m.sleepUpdateCount = SleepUpdateMaxCount
		}
	}
	if m.sleepUpdateCount <= 0 {
		return
	} else {
		m.sleepUpdateCount--
	}

	for name, status := range m.name2status {
		if status.State == downloader.StateAdding {
			continue
		}
		// 文件是否存在
		if len(status.Path) == 0 || utils.IsExist(xpath.Join(Conf.SavePath, status.Path)) ||
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
				status.State = downloader.StateComplete
				status.Init = true
			}
			m.cache.Put(Name2StatusBucket, name, status, 0)
		} else {
			// 文件不存在，检查过期时间
			if status.ExpireAt <= 0 {
				// 未设置过期，设置3小时过期
				status.ExpireAt = time.Now().Add(NotFoundExpireHour * time.Hour).Unix()
				status.State = downloader.StateNotFound
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
			if status.State != downloader.StateNotFound && status.State != downloader.StateAdding {
				// 文件存在，跳过下载
				continue
			}
		}
		status := &models.DownloadStatus{
			Hash:     item.Hash,
			State:    downloader.StateMap(item.State),
			ExpireAt: 0,
		}
		m.name2status[name] = status
		m.cache.Put(Name2StatusBucket, name, status, 0)
	}
}

func (m *Manager) scrape(bangumi *models.AnimeEntity, dir string) bool {
	if len(dir) == 0 {
		return true
	}
	nfo := xpath.Join(Conf.SavePath, dir, "tvshow.nfo")
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
