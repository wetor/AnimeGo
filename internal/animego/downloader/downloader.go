package downloader

import (
	"context"
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
	notifier api.ClientNotifier

	// 保存上一次状态，检查状态是否改变
	stateList map[string]*ItemState

	errs     []error
	errMutex sync.Mutex
	sync.Mutex
}

func NewManager(client api.Client, notifier api.ClientNotifier) *Manager {
	return &Manager{
		client:    client,
		notifier:  notifier,
		stateList: make(map[string]*ItemState, 0),
		errs:      make([]error, 0),
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

func (m *Manager) Add(hash string, opt *client.AddOptions) error {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	m.Lock()
	defer m.Unlock()

	if _, has := m.stateList[hash]; has {
		log.Infof("发现正在下载「%s」", opt.Rename)
		return exceptions.ErrClientExistItem{Client: m.client.Name(), Name: opt.Rename}
	}
	// 添加下载项
	err := m.client.Add(opt)
	if err != nil {
		return err
	}

	m.stateList[hash] = &ItemState{
		Torrent: StateAdding,
		Notify:  NotifyOnNone,
	}
	return nil
}

func (m *Manager) Delete(hash string) error {
	m.Lock()
	defer m.Unlock()

	// 删除下载项
	err := m.client.Delete(&client.DeleteOptions{
		Hash:       []string{hash},
		DeleteFile: true,
	})
	if err != nil {
		return err
	}
	delete(m.stateList, hash)
	return nil
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
		itemState, ok := m.stateList[item.Hash]
		if !ok {
			// 没有记录状态，可能重启，从最初状态开始计算
			itemState = &ItemState{
				Torrent: StateAdding,
				Notify:  NotifyOnStart,
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
