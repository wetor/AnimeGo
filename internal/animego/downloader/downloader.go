package downloader

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
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

func NewManager(client api.Client, notifier api.ClientNotifier, database api.Database) *Manager {
	m := &Manager{
		client:   client,
		notifier: notifier,
		database: database,
	}
	m.Init()
	return m
}

func (m *Manager) Init() {
	m.cache = make(map[string]*models.AnimeEntity)
	m.hash2stateList = make(map[string]*ItemState)
	m.name2hash = make(map[string]string)
	m.errs = make([]error, 0)
}

func (m *Manager) sleep(ctx context.Context) {
	utils.Sleep(RefreshSecond, ctx)
}

func (m *Manager) addError(err error) {
	m.errMutex.Lock()
	m.errs = append(m.errs, err)
	m.errMutex.Unlock()
}

func (m *Manager) transition(oldTorrentState, newTorrentState constant.TorrentState) (constant.TorrentState, NotifyState) {
	if newTorrentState == constant.StateError {
		// error
		return newTorrentState, NotifyOnError
	}

	state := newTorrentState
	result := NotifyOnStart
	switch oldTorrentState {
	case constant.StateInit:
		// init -> start
		result = NotifyOnStart
		state = constant.StateAdding
	case constant.StateAdding:
		switch newTorrentState {
		case constant.StateDownloading:
			// start -> download
			result = NotifyOnDownload
		case constant.StateSeeding:
			// start -> seed
			// 做种状态下重启后
			result = NotifyOnSeeding
		case constant.StateComplete:
			// start -> complete
			// 完成状态下重启后，需要先经过seeding
			result = NotifyOnSeeding
			state = constant.StateSeeding
		}
	case constant.StateDownloading:
		switch newTorrentState {
		case constant.StateDownloading:
			// download -> download
			// 刷新进度
			result = NotifyOnDownload
		case constant.StateSeeding:
			// download -> seed
			result = NotifyOnSeeding
		case constant.StateComplete:
			// download -> complete
			// 跳过了seeding，需要先经过seeding
			result = NotifyOnSeeding
			state = constant.StateSeeding
		}
	case constant.StateSeeding:
		switch newTorrentState {
		case constant.StateSeeding:
			// seed -> seed
			result = NotifyOnSeeding
		case constant.StateComplete:
			// seed -> complete
			result = NotifyOnComplete
		}
	case constant.StateComplete:
		// complete
		result = NotifyOnComplete
	case constant.StateError:
		switch newTorrentState {
		case constant.StateDownloading:
			// error -> download
			result = NotifyOnDownload
		case constant.StateSeeding:
			// error -> seed
			result = NotifyOnSeeding
		case constant.StateComplete:
			// error -> complete
			result = NotifyOnComplete
		}
	}
	log.Debugf("torrent %v -> %v", oldTorrentState, state)
	return state, result
}

func (m *Manager) notify(oldNotifyState, newNotifyState NotifyState, event []models.ClientEvent) error {
	log.Debugf("notify %v -> %v", oldNotifyState, newNotifyState)
	if newNotifyState == NotifyOnInit {
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
	err := m.add(anime.Hash(), &models.AddOptions{
		Url:      anime.Torrent.Url,
		File:     anime.Torrent.File,
		SavePath: m.client.Config().DownloadPath,
		Category: Category,
		Tag:      utils.Tag(Tag, anime.AirDate, anime.Ep[0].Ep),
		Name:     name,
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

func (m *Manager) add(hash string, opt *models.AddOptions) error {
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
		Torrent: constant.StateInit,
		Notify:  NotifyOnInit,
		Name:    name,
	}
	m.name2hash[name] = hash
	return nil
}

func (m *Manager) delete(hash string, deleteItem bool) error {
	if deleteItem {
		log.Debugf("删除下载项：%v", hash)
		// 删除下载项
		err := m.client.Delete(&models.DeleteOptions{
			Hash:       []string{hash},
			DeleteFile: true,
		})
		if err != nil {
			return err
		}
	} else {
		log.Debugf("删除缓存：%v", hash)
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

func (m *Manager) updateList() {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	m.Lock()
	defer m.Unlock()

	items, err := m.client.List(&models.ListOptions{
		Category: Category,
	})
	if err != nil {
		m.addError(err)
		return
	}

	for _, item := range items {
		state := m.client.State(item.State)
		itemState, ok := m.hash2stateList[item.Hash]
		if !ok {
			// 没有记录状态，可能重启，从最初状态开始计算
			itemState = &ItemState{
				Torrent: constant.StateInit,
				Notify:  NotifyOnInit,
				Name:    item.Name,
			}
		}
		if state != itemState.Torrent {
			// 发送通知
			state, notify := m.transition(itemState.Torrent, state)
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
					m.updateList()
					m.sleep(ctx)
				}
			}()
			if exit {
				return
			}
		}
	}()
}
