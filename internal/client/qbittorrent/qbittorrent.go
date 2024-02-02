package qbittorrent

import (
	"time"

	"github.com/google/wire"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/third_party/qbapi"
)

const (
	Name = "qBittorrent"
)

var Set = wire.NewSet(
	NewQBittorrent,
	wire.Bind(new(api.Client), new(*QBittorrent)),
)

type QBittorrent struct {
	option      []qbapi.Option
	connectFunc func() bool
	retryChan   chan int
	retryNum    int // 重试次数

	connected bool
	client    *qbapi.QBAPI

	*models.ClientOptions
}

func NewQBittorrent(opts *models.ClientOptions) *QBittorrent {
	c := &QBittorrent{
		retryChan:     make(chan int, 1),
		retryNum:      1,
		connected:     false,
		ClientOptions: opts,
	}
	c.option = make([]qbapi.Option, 0, 3)
	c.option = append(c.option, qbapi.WithAuth(c.Username, c.Password))
	c.option = append(c.option, qbapi.WithHost(c.Url))
	c.option = append(c.option, qbapi.WithTimeout(time.Duration(c.ConnectTimeoutSecond)*time.Second))

	c.connectFunc = func() bool {
		var err error
		c.client, err = qbapi.NewAPI(c.option...)
		if err != nil {
			log.DebugErr(err)
			log.Warnf("初始化 %s 客户端第%d次，失败", Name, c.retryNum)
			return false
		}
		if err = c.client.Login(c.Ctx); err != nil {
			log.DebugErr(err)
			log.Warnf("连接 %s 第%d次，失败", Name, c.retryNum)
			return false
		}
		return true
	}
	c.retryChan <- constant.ChanRetryConnect
	return c
}

func (c *QBittorrent) Name() string {
	return Name
}

func (c *QBittorrent) Config() *models.Config {
	return &models.Config{
		ApiUrl:       c.Url,
		DownloadPath: c.DownloadPath,
	}
}

// State
//
//	@Description: 下载器状态转换
//	@param state string
//	@return client.TorrentState
func (c *QBittorrent) State(state string) constant.TorrentState {
	switch state {
	case QbtAllocating, QbtMetaDL, QbtStalledDL,
		QbtCheckingDL, QbtCheckingResumeData, QbtQueuedDL,
		QbtForcedUP, QbtQueuedUP:
		// 若进度为100，则下载完成
		return constant.StateWaiting
	case QbtDownloading, QbtForcedDL:
		return constant.StateDownloading
	case QbtMoving:
		return constant.StateMoving
	case QbtUploading, QbtStalledUP:
		// 已下载完成
		return constant.StateSeeding
	case QbtPausedDL:
		return constant.StatePausing
	case QbtPausedUP, QbtCheckingUP:
		// 已下载完成
		return constant.StateComplete
	case QbtError, QbtMissingFiles:
		return constant.StateError
	case QbtUnknown:
		return constant.StateUnknown
	default:
		return constant.StateUnknown
	}
}

func (c *QBittorrent) Connected() bool {
	return c.connected
}

func (c *QBittorrent) clientVersion() string {
	clientResp, err := c.client.GetApplicationVersion(c.Ctx, &qbapi.GetApplicationVersionReq{})
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
	c.WG.Add(1)
	go func() {
		defer c.WG.Done()
		for {
			select {
			case <-c.Ctx.Done():
				log.Debugf("正常退出 %s reconnect listen", Name)
				return
			case msg := <-c.retryChan:
				c.connected = true
				if msg == constant.ChanRetryConnect && (c.client == nil || len(c.clientVersion()) == 0) {
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
	c.WG.Add(1)
	go func() {
		defer c.WG.Done()
		for {
			select {
			case <-c.Ctx.Done():
				log.Debugf("正常退出 %s check listen", Name)
				return
			default:
				if c.retryNum == 0 {
					c.retryChan <- constant.ChanRetryConnect
					// 检查是否在线，时间长
					utils.Sleep(c.CheckTimeSecond, c.Ctx)
				} else if c.retryNum <= c.RetryConnectNum {
					c.retryChan <- constant.ChanRetryConnect
					// 失败重试，时间短
					utils.Sleep(c.ConnectTimeoutSecond, c.Ctx)
				} else {
					// 超过重试次数，不在频繁重试
					c.retryNum = 0
				}
			}
		}
	}()
}

func (c *QBittorrent) List(opt *models.ListOptions) ([]*models.TorrentItem, error) {
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

	listResp, err := c.client.GetTorrentList(c.Ctx, req)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "获取列表失败"})
	}
	retn := make([]*models.TorrentItem, len(listResp.Items))
	err = copier.Copy(&retn, &listResp.Items)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "类型转换失败"})
	}
	return retn, nil
}

func (c *QBittorrent) Add(opt *models.AddOptions) error {
	if !c.connected {
		return errors.WithStack(&exceptions.ErrClientNoConnected{Client: Name})
	}
	var err error
	meta := &qbapi.AddTorrentMeta{
		Savepath:         &opt.SavePath,
		Category:         &opt.Category,
		Tags:             opt.Tag,
		SeedingTimeLimit: &c.SeedingTimeMinute,
		Rename:           &opt.Name,
	}
	if len(opt.File) > 0 {
		_, err = c.client.AddNewTorrent(c.Ctx, &qbapi.AddNewTorrentReq{
			File: []string{opt.File},
			Meta: meta,
		})
	} else {
		_, err = c.client.AddNewLink(c.Ctx, &qbapi.AddNewLinkReq{
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

func (c *QBittorrent) Pause(opt *models.PauseOptions) error {
	if !c.connected {
		return errors.WithStack(&exceptions.ErrClientNoConnected{Client: Name})
	}
	var err error
	if opt.Pause {
		_, err = c.client.PauseTorrents(c.Ctx, &qbapi.PauseTorrentsReq{
			Hash: opt.Hash,
		})
		if err != nil {
			log.DebugErr(err)
			return errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "暂停下载项失败"})
		}
	} else {
		_, err = c.client.ResumeTorrents(c.Ctx, &qbapi.ResumeTorrentsReq{
			Hash: opt.Hash,
		})
		if err != nil {
			log.DebugErr(err)
			return errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "恢复下载项失败"})
		}
	}
	return nil
}

func (c *QBittorrent) Delete(opt *models.DeleteOptions) error {
	if !c.connected {
		return errors.WithStack(&exceptions.ErrClientNoConnected{Client: Name})
	}
	var err error
	_, err = c.client.DeleteTorrents(c.Ctx, &qbapi.DeleteTorrentsReq{
		IsDeleteFile: opt.DeleteFile,
		Hash:         opt.Hash,
	})
	if err != nil {
		log.DebugErr(err)
		return errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "删除下载项失败"})
	}
	return nil
}
