package qbittorrent

import (
	"context"
	"sync"
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
	Name             = "qBittorrent"
	ChanRetryConnect = 1 // 重连消息
)

type QBittorrent struct {
	option      []qbapi.Option
	connectFunc func() bool
	retryChan   chan int
	retryNum    int // 重试次数

	connected bool
	client    *qbapi.QBAPI

	config *Conf
	WG     *sync.WaitGroup
}

func NewQBittorrent(opts *Options) *QBittorrent {
	opts.Default()
	qbt := &QBittorrent{
		retryChan: make(chan int, 1),
		config: &Conf{
			Url:                  opts.Url,
			Username:             opts.Username,
			Password:             opts.Password,
			DownloadPath:         opts.DownloadPath,
			ConnectTimeoutSecond: opts.ConnectTimeoutSecond,
			CheckTimeSecond:      opts.CheckTimeSecond,
			RetryConnectNum:      opts.RetryConnectNum,
		},
		WG: opts.WG,
	}
	qbt.option = make([]qbapi.Option, 0, 3)

	qbt.option = append(qbt.option, qbapi.WithAuth(opts.Username, opts.Password))
	qbt.option = append(qbt.option, qbapi.WithHost(opts.Url))
	qbt.option = append(qbt.option, qbapi.WithTimeout(time.Duration(opts.ConnectTimeoutSecond)*time.Second))
	qbt.retryNum = 1
	qbt.connected = false
	qbt.retryChan <- ChanRetryConnect
	return qbt
}

func (c *QBittorrent) Name() string {
	return Name
}

func (c *QBittorrent) Config() *client.Config {
	return &client.Config{
		ApiUrl:       c.config.Url,
		DownloadPath: c.config.DownloadPath,
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
//	@param ctx context.Context
func (c *QBittorrent) Start(ctx context.Context) {
	c.connectFunc = func() bool {
		var err error
		c.client, err = qbapi.NewAPI(c.option...)
		if err != nil {
			log.DebugErr(err)
			log.Warnf("初始化 %s 客户端第%d次，失败", Name, c.retryNum)
			return false
		}
		if err = c.client.Login(ctx); err != nil {
			log.DebugErr(err)
			log.Warnf("连接 %s 第%d次，失败", Name, c.retryNum)
			return false
		}
		// c.Init()
		return true
	}
	c.WG.Add(1)
	go func() {
		defer c.WG.Done()
		for {
			select {
			case <-ctx.Done():
				log.Debugf("正常退出 %s reconnect listen", Name)
				return
			case msg := <-c.retryChan:
				c.connected = true
				if msg == ChanRetryConnect && (c.client == nil || len(c.clientVersion()) == 0) {
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
			case <-ctx.Done():
				log.Debugf("正常退出 %s check listen", Name)
				return
			default:
				if c.retryNum == 0 {
					c.retryChan <- ChanRetryConnect
					// 检查是否在线，时间长
					utils.Sleep(c.config.CheckTimeSecond, ctx)
				} else if c.retryNum <= c.config.RetryConnectNum {
					c.retryChan <- ChanRetryConnect
					// 失败重试，时间短
					utils.Sleep(c.config.ConnectTimeoutSecond, ctx)
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
		SeedingTimeLimit: &opt.SeedingTime,
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
