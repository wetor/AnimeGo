package transmission

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/hekmon/transmissionrpc/v3"
	"github.com/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/client"
	"github.com/wetor/AnimeGo/pkg/exceptions"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	Name             = "Transmission"
	ChanRetryConnect = 1 // 重连消息
)

type Transmission struct {
	connectFunc func() bool
	retryChan   chan int
	retryNum    int // 重试次数

	connected bool
	client    *transmissionrpc.Client

	config   *Conf
	endpoint *url.URL

	WG  *sync.WaitGroup
	ctx context.Context
}

func NewTransmission(opts *Options) *Transmission {
	opts.Default()
	c := &Transmission{
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
		WG:  opts.WG,
		ctx: opts.Ctx,
	}
	u, _ := url.Parse(c.config.Url)
	c.endpoint, _ = url.Parse(fmt.Sprintf("%s://%s:%s@%s/transmission/rpc",
		u.Scheme, c.config.Username, c.config.Password, u.Host))
	return c
}

func (c *Transmission) Name() string {
	return Name
}

func (c *Transmission) Config() *client.Config {
	return &client.Config{
		ApiUrl:       c.config.Url,
		DownloadPath: c.config.DownloadPath,
	}
}

// State
//
//	@Description: 下载器状态转换
//	@param state string
//	@return client.TorrentState
func (c *Transmission) State(state string) client.TorrentState {
	switch state {
	case TorrentStatusCheckWait, TorrentStatusCheck,
		TorrentStatusDownloadWait, TorrentStatusSeedWait:
		// 若进度为100，则下载完成
		return client.StateWaiting
	case TorrentStatusDownload:
		return client.StateDownloading
	case TorrentStatusSeed:
		// 已下载完成
		return client.StateSeeding
	case TorrentStatusStopped:
		return client.StatePausing
	case TorrentStatusComplete:
		// 已下载完成
		return client.StateComplete
	case TorrentStatusIsolated:
		return client.StateError
	default:
		return client.StateUnknown
	}
}

func (c *Transmission) Connected() bool {
	return c.connected
}

func (c *Transmission) clientVersion() string {
	ok, version, _, err := c.client.RPCVersion(c.ctx)
	if err != nil || !ok {
		return ""
	}
	return strconv.Itoa(int(version))
}

// Start
//
//	@Description: 启动下载器协程
//	@Description: 客户端在线监听、登录重试
//	@Description: 客户端处理下载消息，获取下载进度
//	@receiver *QBittorrent

// @param ctx context.Context
func (c *Transmission) Start() {
	c.connectFunc = func() bool {
		var err error
		c.client, err = transmissionrpc.New(c.endpoint, nil)
		if err != nil {
			log.DebugErr(err)
			log.Warnf("初始化 %s 客户端第%d次，失败", Name, c.retryNum)
			return false
		}
		ok, _, miniVersion, err := c.client.RPCVersion(c.ctx)
		if err != nil {
			log.DebugErr(err)
			log.Warnf("连接 %s 第%d次，失败", Name, c.retryNum)
			return false
		}
		if !ok {
			log.Warnf("连接 %s 失败。最小支持RPC版本 %s，当前RPC版本 %s", Name,
				c.retryNum, miniVersion, transmissionrpc.RPCVersion)
			return false
		}
		return true
	}
	c.WG.Add(1)
	go func() {
		defer c.WG.Done()
		for {
			select {
			case <-c.ctx.Done():
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
			case <-c.ctx.Done():
				log.Debugf("正常退出 %s check listen", Name)
				return
			default:
				if c.retryNum == 0 {
					c.retryChan <- ChanRetryConnect
					// 检查是否在线，时间长
					utils.Sleep(c.config.CheckTimeSecond, c.ctx)
				} else if c.retryNum <= c.config.RetryConnectNum {
					c.retryChan <- ChanRetryConnect
					// 失败重试，时间短
					utils.Sleep(c.config.ConnectTimeoutSecond, c.ctx)
				} else {
					// 超过重试次数，不在频繁重试
					c.retryNum = 0
				}
			}
		}
	}()
}

func (c *Transmission) List(opt *client.ListOptions) ([]*client.TorrentItem, error) {
	if !c.connected {
		return nil, errors.WithStack(&exceptions.ErrClientNoConnected{Client: Name})
	}
	torrents, err := c.client.TorrentGetAll(c.ctx)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "获取列表失败"})
	}
	items := make([]*client.TorrentItem, 0, len(torrents))

	for _, torrent := range torrents {
		if *torrent.Group == opt.Category {
			item := &client.TorrentItem{
				Hash:     *torrent.HashString,
				State:    torrent.Status.String(),
				Progress: *torrent.PercentDone,
			}
			if int(item.Progress) == 1 && item.State == TorrentStatusStopped {
				// 下载进度100% 且停止状态，即完成
				item.State = TorrentStatusComplete
			}
			items = append(items, item)
		}
	}
	return items, nil
}

func (c *Transmission) Add(opt *client.AddOptions) error {
	if !c.connected {
		return errors.WithStack(&exceptions.ErrClientNoConnected{Client: Name})
	}
	var err error
	var filename string
	if len(opt.File) > 0 {
		filename = opt.File
	} else {
		filename = opt.Url
	}
	payload := transmissionrpc.TorrentAddPayload{
		Filename:    &filename,
		DownloadDir: &opt.SavePath,
		Labels:      []string{opt.Tag},
		//Name:             &opt.Name,
		//Category:         &opt.Category,
		//Tags:             opt.Tag,
		//SeedingTimeLimit: &opt.SeedingTime,
	}
	torrent, err := c.client.TorrentAdd(c.ctx, payload)
	if err != nil {
		log.DebugErr(err)
		return errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "添加下载项失败"})
	}
	mode := IdleModeGlobal
	if opt.SeedingTime > 0 {
		mode = IdleModeSingle // 空闲指定时间后停止做种
	} else if opt.SeedingTime < 0 {
		mode = IdleModeUnlimited // 始终做种
	}
	seedTime := time.Duration(opt.SeedingTime) * time.Minute

	err = c.client.TorrentSet(c.ctx, transmissionrpc.TorrentSetPayload{
		IDs:           []int64{*torrent.ID},
		Group:         &opt.Category,
		SeedIdleMode:  &mode,
		SeedIdleLimit: &seedTime,
	})
	if err != nil {
		log.DebugErr(err)
		// 设置状态失败，删除下载
		_ = c.client.TorrentRemove(c.ctx, transmissionrpc.TorrentRemovePayload{
			IDs:             []int64{*torrent.ID},
			DeleteLocalData: true,
		})
		return errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "添加下载项失败"})
	}
	return nil
}

func (c *Transmission) Delete(opt *client.DeleteOptions) error {
	if !c.connected {
		return errors.WithStack(&exceptions.ErrClientNoConnected{Client: Name})
	}
	torrents, err := c.client.TorrentGetHashes(c.ctx, []string{"id"}, opt.Hash)
	if err != nil {
		log.DebugErr(err)
		return errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "删除下载项失败"})
	}
	ids := make([]int64, 0, len(torrents))
	for _, torrent := range torrents {
		ids = append(ids, *torrent.ID)
	}
	err = c.client.TorrentRemove(c.ctx, transmissionrpc.TorrentRemovePayload{
		IDs:             ids,
		DeleteLocalData: opt.DeleteFile,
	})
	if err != nil {
		log.DebugErr(err)
		return errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "删除下载项失败"})
	}
	return nil
}
