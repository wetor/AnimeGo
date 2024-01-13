package qbittorrent

import (
	"context"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/pkg/client"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/third_party/qbapi"
)

const (
	Name = "qBittorrent"
)

type QBittorrent struct {
	option      []qbapi.Option
	connectFunc func() bool
	retryChan   chan int
	retryNum    int // 重试次数

	connected bool
	auth      *client.AuthOptions
	client    *qbapi.QBAPI
}

func NewQBittorrent(opt *client.AuthOptions) *QBittorrent {
	c := &QBittorrent{
		retryChan: make(chan int, 1),
		retryNum:  1,
		connected: false,
		auth:      opt,
	}
	c.option = make([]qbapi.Option, 0, 3)
	c.option = append(c.option, qbapi.WithAuth(c.auth.Username, c.auth.Password))
	c.option = append(c.option, qbapi.WithHost(c.auth.Url))
	c.option = append(c.option, qbapi.WithTimeout(time.Duration(client.ConnectTimeoutSecond)*time.Second))

	c.connectFunc = func() bool {
		var err error
		c.client, err = qbapi.NewAPI(c.option...)
		if err != nil {
			log.DebugErr(err)
			log.Warnf("初始化 %s 客户端第%d次，失败", Name, c.retryNum)
			return false
		}
		if err = c.client.Login(client.Ctx); err != nil {
			log.DebugErr(err)
			log.Warnf("连接 %s 第%d次，失败", Name, c.retryNum)
			return false
		}
		return true
	}
	c.retryChan <- client.ChanRetryConnect
	return c
}

func (c *QBittorrent) Name() string {
	return Name
}

func (c *QBittorrent) Config() *client.Config {
	return &client.Config{
		ApiUrl:       c.auth.Url,
		DownloadPath: client.DownloadPath,
	}
}

// State
//
//	@Description: 下载器状态转换
//	@param state string
//	@return client.TorrentState
func (c *QBittorrent) State(state string) client.TorrentState {
	switch state {
	case QbtAllocating, QbtMetaDL, QbtStalledDL,
		QbtCheckingDL, QbtCheckingResumeData, QbtQueuedDL,
		QbtForcedUP, QbtQueuedUP:
		// 若进度为100，则下载完成
		return client.StateWaiting
	case QbtDownloading, QbtForcedDL:
		return client.StateDownloading
	case QbtMoving:
		return client.StateMoving
	case QbtUploading, QbtStalledUP:
		// 已下载完成
		return client.StateSeeding
	case QbtPausedDL:
		return client.StatePausing
	case QbtPausedUP, QbtCheckingUP:
		// 已下载完成
		return client.StateComplete
	case QbtError, QbtMissingFiles:
		return client.StateError
	case QbtUnknown:
		return client.StateUnknown
	default:
		return client.StateUnknown
	}
}

func (c *QBittorrent) Connected() bool {
	return c.connected
}

func (c *QBittorrent) clientVersion() string {
	clientResp, err := c.client.GetApplicationVersion(context.Background(), &qbapi.GetApplicationVersionReq{})
	if err != nil {
		return ""
	}
	return clientResp.Version
}

// Start
//
//	@Description: 启动下载器协程
//	@Description: 客户端在线监听、登录重试
//	@Description: 客户端处理下载消息，获取下载进度
//	@receiver *QBittorrent
func (c *QBittorrent) Start() {
	client.WG.Add(1)
	go func() {
		defer client.WG.Done()
		for {
			select {
			case <-client.Ctx.Done():
				log.Debugf("正常退出 %s reconnect listen", Name)
				return
			case msg := <-c.retryChan:
				c.connected = true
				if msg == client.ChanRetryConnect && (c.client == nil || len(c.clientVersion()) == 0) {
					if ok := c.connectFunc(); !ok {
						c.retryNum++
						c.connected = false
						// 重连失败
					} else {
						// 重连成功
						c.retryNum = 0
						c.connected = true
						log.Infof("连接 %s 成功", Name)
					}
				}
			}
		}
	}()
	client.WG.Add(1)
	go func() {
		defer client.WG.Done()
		for {
			select {
			case <-client.Ctx.Done():
				log.Debugf("正常退出 %s check listen", Name)
				return
			default:
				if c.retryNum == 0 {
					c.retryChan <- client.ChanRetryConnect
					// 检查是否在线，时间长
					utils.Sleep(client.CheckTimeSecond, client.Ctx)
				} else if c.retryNum <= client.RetryConnectNum {
					c.retryChan <- client.ChanRetryConnect
					// 失败重试，时间短
					utils.Sleep(client.ConnectTimeoutSecond, client.Ctx)
				} else {
					// 超过重试次数，不在频繁重试
					c.retryNum = 0
				}
			}
		}
	}()
}

func (c *QBittorrent) List(opt *client.ListOptions) ([]*client.TorrentItem, error) {
	if !c.connected {
		return nil, errors.WithStack(&exceptions.ErrClientNoConnected{Client: Name})
	}
	req := &qbapi.GetTorrentListReq{}
	if len(opt.Status) != 0 {
		req.Filter = &opt.Status
	}
	if len(opt.Category) != 0 {
		req.Category = &opt.Category
	}
	if len(opt.Tag) != 0 {
		req.Tag = &opt.Tag
	}

	listResp, err := c.client.GetTorrentList(context.Background(), req)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "获取列表失败"})
	}
	retn := make([]*client.TorrentItem, len(listResp.Items))
	err = copier.Copy(&retn, &listResp.Items)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "类型转换失败"})
	}
	return retn, nil
}

func (c *QBittorrent) Add(opt *client.AddOptions) error {
	if !c.connected {
		return errors.WithStack(&exceptions.ErrClientNoConnected{Client: Name})
	}
	var err error
	meta := &qbapi.AddTorrentMeta{
		Savepath:         &opt.SavePath,
		Category:         &opt.Category,
		Tags:             opt.Tag,
		SeedingTimeLimit: &client.SeedingTimeMinute,
		Rename:           &opt.Name,
	}
	if len(opt.File) > 0 {
		_, err = c.client.AddNewTorrent(context.Background(), &qbapi.AddNewTorrentReq{
			File: []string{opt.File},
			Meta: meta,
		})
	} else {
		_, err = c.client.AddNewLink(context.Background(), &qbapi.AddNewLinkReq{
			Url:  []string{opt.Url},
			Meta: meta,
		})
	}
	if err != nil {
		log.DebugErr(err)
		return errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "添加下载项失败"})
	}
	return nil
}

func (c *QBittorrent) Delete(opt *client.DeleteOptions) error {
	if !c.connected {
		return errors.WithStack(&exceptions.ErrClientNoConnected{Client: Name})
	}
	_, err := c.client.DeleteTorrents(context.Background(), &qbapi.DeleteTorrentsReq{
		IsDeleteFile: opt.DeleteFile,
		Hash:         opt.Hash,
	})
	if err != nil {
		log.DebugErr(err)
		return errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "删除下载项失败"})
	}
	return nil
}
