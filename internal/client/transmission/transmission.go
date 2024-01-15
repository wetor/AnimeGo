package transmission

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/google/wire"
	"github.com/hekmon/transmissionrpc/v3"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	Name = "Transmission"
)

var Set = wire.NewSet(
	NewTransmission,
	wire.Bind(new(api.Client), new(*Transmission)),
)

type Transmission struct {
	connectFunc func() bool
	retryChan   chan int
	retryNum    int // 重试次数

	connected bool
	client    *transmissionrpc.Client
	endpoint  *url.URL

	*models.ClientOptions
}

func NewTransmission(opts *models.ClientOptions) *Transmission {
	c := &Transmission{
		retryChan:     make(chan int, 1),
		retryNum:      1,
		connected:     false,
		ClientOptions: opts,
	}
	u, _ := url.Parse(c.Url)
	c.endpoint, _ = url.Parse(fmt.Sprintf("%s://%s:%s@%s/transmission/rpc",
		u.Scheme, c.Username, c.Password, u.Host))
	c.connectFunc = func() bool {
		var err error
		c.client, err = transmissionrpc.New(c.endpoint, nil)
		if err != nil {
			log.DebugErr(err)
			log.Warnf("初始化 %s 客户端第%d次，失败", Name, c.retryNum)
			return false
		}
		ok, _, miniVersion, err := c.client.RPCVersion(c.Ctx)
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
	c.retryChan <- constant.ChanRetryConnect
	return c
}

func (c *Transmission) Name() string {
	return Name
}

func (c *Transmission) Config() *models.Config {
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
func (c *Transmission) State(state string) constant.TorrentState {
	switch state {
	case TorrentStatusCheckWait, TorrentStatusCheck,
		TorrentStatusDownloadWait, TorrentStatusSeedWait:
		// 若进度为100，则下载完成
		return constant.StateWaiting
	case TorrentStatusDownload:
		return constant.StateDownloading
	case TorrentStatusSeed:
		// 已下载完成
		return constant.StateSeeding
	case TorrentStatusStopped:
		return constant.StatePausing
	case TorrentStatusComplete:
		// 已下载完成
		return constant.StateComplete
	case TorrentStatusIsolated:
		return constant.StateError
	default:
		return constant.StateUnknown
	}
}

func (c *Transmission) Connected() bool {
	return c.connected
}

func (c *Transmission) clientVersion() string {
	ok, version, _, err := c.client.RPCVersion(c.Ctx)
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
func (c *Transmission) Start() {
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

func (c *Transmission) List(opt *models.ListOptions) ([]*models.TorrentItem, error) {
	if !c.connected {
		return nil, errors.WithStack(&exceptions.ErrClientNoConnected{Client: Name})
	}
	torrents, err := c.client.TorrentGetAll(c.Ctx)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "获取列表失败"})
	}
	items := make([]*models.TorrentItem, 0, len(torrents))

	for _, torrent := range torrents {
		if len(opt.Category) > 0 && *torrent.Group != opt.Category {
			continue
		}
		if len(opt.Tag) > 0 && !slices.Contains[string](torrent.Labels, opt.Tag) {
			continue
		}
		status := torrent.Status.String()
		if int(*torrent.PercentDone) == 1 && torrent.Status.String() == TorrentStatusStopped {
			// TODO: 下载进度100% 且停止状态，即完成
			status = TorrentStatusComplete
		}
		if len(opt.Status) > 0 && status != opt.Status {
			continue
		}
		items = append(items, &models.TorrentItem{
			Hash:     *torrent.HashString,
			State:    status,
			Progress: *torrent.PercentDone,
		})
	}
	return items, nil
}

func (c *Transmission) Add(opt *models.AddOptions) error {
	if !c.connected {
		return errors.WithStack(&exceptions.ErrClientNoConnected{Client: Name})
	}
	payload := transmissionrpc.TorrentAddPayload{
		DownloadDir: &opt.SavePath,
		Labels:      []string{opt.Tag},
	}
	if len(opt.File) > 0 {
		filename, err := transmissionrpc.File2Base64(opt.File)
		if err != nil {
			return err
		}
		payload.MetaInfo = &filename
	} else {
		payload.Filename = &opt.Url
	}
	torrent, err := c.client.TorrentAdd(c.Ctx, payload)
	if err != nil {
		log.DebugErr(err)
		return errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "添加下载项失败"})
	}
	mode := IdleModeGlobal
	if c.SeedingTimeMinute > 0 {
		mode = IdleModeSingle // 空闲指定时间后停止做种
	} else if c.SeedingTimeMinute < 0 {
		mode = IdleModeUnlimited // 始终做种
	}
	seedTime := time.Duration(c.SeedingTimeMinute) * time.Minute

	err = c.client.TorrentSet(c.Ctx, transmissionrpc.TorrentSetPayload{
		IDs:           []int64{*torrent.ID},
		Group:         &opt.Category,
		SeedIdleMode:  &mode,
		SeedIdleLimit: &seedTime,
	})
	if err != nil {
		log.DebugErr(err)
		// 设置状态失败，删除下载
		_ = c.client.TorrentRemove(c.Ctx, transmissionrpc.TorrentRemovePayload{
			IDs:             []int64{*torrent.ID},
			DeleteLocalData: true,
		})
		return errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "添加下载项失败"})
	}
	return nil
}

func (c *Transmission) Delete(opt *models.DeleteOptions) error {
	if !c.connected {
		return errors.WithStack(&exceptions.ErrClientNoConnected{Client: Name})
	}
	torrents, err := c.client.TorrentGetHashes(c.Ctx, []string{"id"}, opt.Hash)
	if err != nil {
		log.DebugErr(err)
		return errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "删除下载项失败"})
	}
	ids := make([]int64, 0, len(torrents))
	for _, torrent := range torrents {
		ids = append(ids, *torrent.ID)
	}
	err = c.client.TorrentRemove(c.Ctx, transmissionrpc.TorrentRemovePayload{
		IDs:             ids,
		DeleteLocalData: opt.DeleteFile,
	})
	if err != nil {
		log.DebugErr(err)
		return errors.WithStack(&exceptions.ErrClient{Client: Name, Message: "删除下载项失败"})
	}
	return nil
}
