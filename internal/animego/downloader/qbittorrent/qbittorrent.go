package qbittorrent

import (
	"context"
	"time"

	"github.com/jinzhu/copier"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/third_party/qbapi"
)

const (
	ChanRetryConnect = 1 // 重连消息
)

type QBittorrent struct {
	option      []qbapi.Option
	connectFunc func() bool
	retryChan   chan int
	retryNum    int // 重试次数

	connected bool
	client    *qbapi.QBAPI
}

func NewQBittorrent(url, username, password string) *QBittorrent {
	qbt := &QBittorrent{
		retryChan: make(chan int, 1),
		retryNum:  0,
	}
	qbt.option = make([]qbapi.Option, 0, 3)

	qbt.option = append(qbt.option, qbapi.WithAuth(username, password))
	qbt.option = append(qbt.option, qbapi.WithHost(url))
	qbt.option = append(qbt.option, qbapi.WithTimeout(time.Duration(downloader.ConnectTimeoutSecond)*time.Second))
	qbt.retryNum = 1
	qbt.connected = false
	qbt.retryChan <- ChanRetryConnect
	return qbt
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
			log.Debugf("", errors.NewAniErrorD(err))
			log.Warnf("初始化QBittorrent客户端第%d次，失败", c.retryNum)
			return false
		}
		if err = c.client.Login(ctx); err != nil {
			log.Debugf("", errors.NewAniErrorD(err))
			log.Warnf("连接QBittorrent第%d次，失败", c.retryNum)
			return false
		}
		// c.Init()
		return true
	}
	downloader.WG.Add(1)
	go func() {
		defer downloader.WG.Done()
		for {
			exit := false
			func() {
				defer errors.HandleError(func(err error) {
					log.Errorf("", err)
				})
				select {
				case <-ctx.Done():
					log.Debugf("正常退出 qbittorrent 1")
					exit = true
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
							log.Infof("连接QBittorrent成功")
						}
					}
				}
			}()
			if exit {
				return
			}
		}
	}()
	downloader.WG.Add(1)
	go func() {
		defer downloader.WG.Done()
		for {
			exit := false
			func() {
				defer errors.HandleError(func(err error) {
					log.Errorf("", err)
				})
				select {
				case <-ctx.Done():
					log.Debugf("正常退出 qbittorrent 2")
					exit = true
					return
				default:
					if c.retryNum == 0 {
						c.retryChan <- ChanRetryConnect
						// 检查是否在线，时间长
						utils.Sleep(downloader.CheckTimeSecond, ctx)
					} else if c.retryNum <= downloader.RetryConnectNum {
						c.retryChan <- ChanRetryConnect
						// 失败重试，时间短
						utils.Sleep(downloader.ConnectTimeoutSecond, ctx)
					} else {
						// 超过重试次数，不在频繁重试
						c.retryNum = 0
					}
				}
			}()
			if exit {
				return
			}
		}
	}()

}

// checkError
//
//	@Description: 检查错误，返回是否需要结束流程
//	@receiver *QBittorrent
//	@param err error
//	@return bool
func (c *QBittorrent) checkError(err error) bool {
	if err == nil {
		return false
	}
	if qerror, ok := err.(*qbapi.QError); ok && qerror.Code() == -10004 {
		log.Debugf("", errors.NewAniErrorSkipf(2, "请求失败，等待客户端响应").SetData(err))
		c.retryNum = 1
		c.connected = false
		c.retryChan <- ChanRetryConnect
	} else {
		log.Debugf("", errors.NewAniErrorSkipf(2, "").SetData(err))
		log.Warnf("请求QBittorrent接口失败")
	}
	return true
}

func (c *QBittorrent) List(opt *models.ClientListOptions) []*models.TorrentItem {
	if !c.connected {
		return nil
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
	if c.checkError(err) {
		return nil
	}
	retn := make([]*models.TorrentItem, len(listResp.Items))
	err = copier.Copy(&retn, &listResp.Items)
	if c.checkError(err) {
		return nil
	}
	return retn
}

func (c *QBittorrent) Add(opt *models.ClientAddOptions) {
	if !c.connected {
		return
	}
	var err error
	meta := &qbapi.AddTorrentMeta{
		Savepath:         &opt.SavePath,
		Category:         &opt.Category,
		Tags:             opt.Tag,
		SeedingTimeLimit: &opt.SeedingTime,
		Rename:           &opt.Rename,
	}
	if len(opt.File) > 0 {
		_, err = c.client.AddNewTorrent(context.Background(), &qbapi.AddNewTorrentReq{
			File: []string{opt.File},
		})
	} else {
		_, err = c.client.AddNewLink(context.Background(), &qbapi.AddNewLinkReq{
			Url:  []string{opt.Url},
			Meta: meta,
		})
	}

	if c.checkError(err) {
		return
	}
}

func (c *QBittorrent) Delete(opt *models.ClientDeleteOptions) {
	if !c.connected {
		return
	}
	_, err := c.client.DeleteTorrents(context.Background(), &qbapi.DeleteTorrentsReq{
		IsDeleteFile: opt.DeleteFile,
		Hash:         opt.Hash,
	})
	if c.checkError(err) {
		return
	}
}

func (c *QBittorrent) Get(opt *models.ClientGetOptions) *models.TorrentItem {
	if !c.connected {
		return nil
	}
	req := &qbapi.GetTorrentListReq{
		Hashes: &opt.Hash,
	}
	listResp, err := c.client.GetTorrentList(context.Background(), req)
	if c.checkError(err) || len(listResp.Items) == 0 {
		return nil
	}
	retn := &models.TorrentItem{}
	err = copier.Copy(retn, listResp.Items[0])
	if c.checkError(err) {
		return nil
	}
	return retn
}

func (c *QBittorrent) GetContent(opt *models.ClientGetOptions) []*models.TorrentContentItem {
	if !c.connected {
		return nil
	}
	contents, err := c.client.GetTorrentContents(context.Background(), &qbapi.GetTorrentContentsReq{
		Hash: opt.Hash,
	})
	if c.checkError(err) {
		return nil
	}
	if opt.Item == nil {
		opt.Item = c.Get(opt)
	}
	retn := make([]*models.TorrentContentItem, len(contents.Contents))
	err = copier.Copy(&retn, &contents.Contents)
	if c.checkError(err) {
		return nil
	}
	return retn
}
