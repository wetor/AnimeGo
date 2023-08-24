package downloader

import (
	"context"
	"github.com/pkg/errors"
	"sync"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/client"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

type Manager struct {
	client   api.Client
	database api.Database
	notifier api.ClientNotifier

	cache map[string]*models.AnimeEntity
	// 保存上一次状态，检查状态是否改变
	hash2stateList map[string]*ItemState
	name2hash      map[string]string

	errs     []error
	errMutex sync.Mutex
	sync.Mutex
}

func NewManager(client api.Client, notifier api.ClientNotifier) *Manager {
	return &Manager{
		client:         client,
		notifier:       notifier,
		cache:          make(map[string]*models.AnimeEntity),
		hash2stateList: make(map[string]*ItemState),
		name2hash:      make(map[string]string),
		errs:           make([]error, 0),
	}
}

func (m *Manager) sleep(ctx context.Context) {
	utils.Sleep(RefreshSecond, ctx)
}

func (m *Manager) addError(err error) {
	m.errMutex.Lock()
	m.errs = append(m.errs, err)
	m.errMutex.Unlock()
}

func (m *Manager) transition(oldTorrentState, newTorrentState models.TorrentState) NotifyState {
	log.Debugf("torrent %v -> %v", oldTorrentState, newTorrentState)
	result := NotifyOnNone
	if newTorrentState == StateError {
		// error
		return NotifyOnError
	}

	switch oldTorrentState {
	case StateAdding:
		switch newTorrentState {
		case StateDownloading:
			// start -> download
			result = NotifyOnStart
		case StateSeeding:
			// start -> seed
			// 非常规
			result = NotifyOnSeeding
		case StateComplete:
			// start -> complete
			// 非常规
			result = NotifyOnComplete
		}
	case StateDownloading:
		switch newTorrentState {
		case StateDownloading:
			// download -> download
			// 刷新进度
			result = NotifyOnDownload
		case StateSeeding:
			// download -> seed
			result = NotifyOnSeeding
		case StateComplete:
			// download -> complete
			result = NotifyOnComplete
		}
	case StateSeeding:
		switch newTorrentState {
		case StateSeeding:
			// seed -> seed
			result = NotifyOnSeeding
		case StateComplete:
			// seed -> complete
			result = NotifyOnComplete
		}
	case StateComplete:
		// complete
		result = NotifyOnComplete
	case StateError:
		switch newTorrentState {
		case StateDownloading:
			// error -> download
			result = NotifyOnStart
		case StateSeeding:
			// error -> seed
			result = NotifyOnSeeding
		case StateComplete:
			// error -> complete
			result = NotifyOnComplete
		}
	}
	return result
}

func (m *Manager) notify(oldNotifyState, newNotifyState NotifyState, event []models.ClientEvent) error {
	log.Debugf("notify %v -> %v", oldNotifyState, newNotifyState)
	if newNotifyState == NotifyOnNone {
		return nil
	}
	switch newNotifyState {
	case NotifyOnStart:
		if oldNotifyState == NotifyOnStart {
			break
		}
		m.notifier.OnDownloadStart(event)
	case NotifyOnDownload:
		if oldNotifyState == NotifyOnDownload {
			// do something
			break
		}
	case NotifyOnSeeding:
		if oldNotifyState == NotifyOnComplete {
			// do something
			break
		}
		m.notifier.OnDownloadSeeding(event)
	case NotifyOnComplete:
		if oldNotifyState == NotifyOnComplete {
			break
		}
		m.notifier.OnDownloadComplete(event)
		m.Lock()
		defer m.Unlock()
		for _, e := range event {
			err := m.delete(e.Hash, false)
			if err != nil {
				return err
			}
		}
	case NotifyOnStop:
		if oldNotifyState == NotifyOnStop {
			break
		}
		m.notifier.OnDownloadStop(event)
	case NotifyOnError:
		if oldNotifyState == NotifyOnError {
			break
		}
		m.notifier.OnDownloadError(event)
	}
	return nil
}

func (m *Manager) Download(anime *models.AnimeEntity) error {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	name := anime.FullName()
	if m.database.IsExist(anime) {
		log.Infof("发现已下载「%s」", name)
		if !AllowDuplicateDownload {
			log.Infof("取消下载，不允许重复「%s」", name)
			return exceptions.ErrDownloadExist{Name: name}
		}
	}
	log.Infof("添加下载「%s」", name)
	err := m.Add(anime.Hash(), &client.AddOptions{
		Url:         anime.Torrent.Url,
		File:        anime.Torrent.File,
		SavePath:    m.client.Config().DownloadPath,
		Category:    Category,
		Tag:         utils.Tag(Tag, anime.AirDate, anime.Ep[0].Ep),
		SeedingTime: SeedingTimeMinute,
		Name:        name,
	})
	if err != nil {
		return errors.Wrap(err, "添加下载项失败")
	}
	err = m.database.Add(anime)
	if err != nil {
		return errors.Wrap(err, "添加下载项失败")
	}
	m.cache[anime.Hash()] = anime
	return nil
}

func (m *Manager) Add(hash string, opt *client.AddOptions) error {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	m.Lock()
	defer m.Unlock()
	name := opt.Name
	if _, has := m.hash2stateList[hash]; has {
		log.Infof("发现正在下载「%s」", name)
		return exceptions.ErrClientExistItem{Client: m.client.Name(), Name: name}
	}
	if _, has := m.name2hash[name]; has {
		log.Infof("发现正在下载「%s」", name)
		return exceptions.ErrClientExistItem{Client: m.client.Name(), Name: name}
	}
	// 添加下载项
	err := m.client.Add(opt)
	if err != nil {
		return err
	}

	m.hash2stateList[hash] = &ItemState{
		Torrent: StateAdding,
		Notify:  NotifyOnNone,
		Name:    name,
	}
	m.name2hash[name] = hash
	return nil
}

func (m *Manager) delete(hash string, deleteItem bool) error {
	if deleteItem {
		// 删除下载项
		err := m.client.Delete(&client.DeleteOptions{
			Hash:       []string{hash},
			DeleteFile: true,
		})
		if err != nil {
			return err
		}
	}
	delete(m.cache, hash)
	if state, ok := m.hash2stateList[hash]; ok {
		delete(m.name2hash, state.Name)
		delete(m.hash2stateList, hash)
	}
	return nil
}

func (m *Manager) Delete(hash string) error {
	m.Lock()
	defer m.Unlock()
	return m.delete(hash, true)
}

func (m *Manager) UpdateList() {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	m.Lock()
	defer m.Unlock()

	items, err := m.client.List(&client.ListOptions{
		Category: Category,
	})
	if err != nil {
		m.addError(err)
		return
	}

	for _, item := range items {
		state := StateMap(item.State)
		itemState, ok := m.hash2stateList[item.Hash]
		if !ok {
			// 没有记录状态，可能重启，从最初状态开始计算
			itemState = &ItemState{
				Torrent: StateAdding,
				Notify:  NotifyOnStart,
				Name:    item.Name,
			}
		}
		if state != itemState.Torrent {
			// 发送通知
			notify := m.transition(itemState.Torrent, state)
			err = m.notify(itemState.Notify, notify, []models.ClientEvent{
				{Hash: item.Hash},
			})
			if err != nil {
				m.addError(err)
				// 失败重试
				continue
			}
			itemState.Notify = notify
			itemState.Torrent = state
		}
		m.hash2stateList[item.Hash] = itemState
		m.name2hash[item.Name] = item.Hash
	}
}

func (m *Manager) Start(ctx context.Context) {
	WG.Add(1)
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
					log.Debugf("正常退出 manager downloader")
					exit = true
					return
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
