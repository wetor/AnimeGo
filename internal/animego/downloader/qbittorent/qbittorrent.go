package qbittorent

import (
	"AnimeGo/internal/models"
	"AnimeGo/internal/store"
	"AnimeGo/internal/utils"
	"AnimeGo/third_party/qbapi"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

const (
	QbtError              = "error"              // Some error occurred, applies to paused torrents
	QbtMissingFiles       = "missingFiles"       // Torrent data files is missing
	QbtUploading          = "uploading"          // Torrent is being seeded and data is being transferred
	QbtPausedUP           = "pausedUP"           // Torrent is paused and has finished downloading
	QbtQueuedUP           = "queuedUP"           // Queuing is enabled and torrent is queued for upload
	QbtStalledUP          = "stalledUP"          // Torrent is being seeded, but no connection were made
	QbtCheckingUP         = "checkingUP"         // Torrent has finished downloading and is being checked
	QbtForcedUP           = "forcedUP"           // Torrent is forced to uploading and ignore queue limit
	QbtAllocating         = "allocating"         // Torrent is allocating disk space for download
	QbtDownloading        = "downloading"        // Torrent is being downloaded and data is being transferred
	QbtMetaDL             = "metaDL"             // Torrent has just started downloading and is fetching metadata
	QbtPausedDL           = "pausedDL"           // Torrent is paused and has NOT finished downloading
	QbtQueuedDL           = "queuedDL"           // Queuing is enabled and torrent is queued for download
	QbtStalledDL          = "stalledDL"          // Torrent is being downloaded, but no connection were made
	QbtCheckingDL         = "checkingDL"         // Same as checkingUP, but torrent has NOT finished downloading
	QbtForcedDL           = "forcedDL"           // Torrent is forced to downloading to ignore queue limit
	QbtCheckingResumeData = "checkingResumeData" // Checking resume data on qBt startup
	QbtMoving             = "moving"             // Torrent is moving to another location
	QbtUnknown            = "unknown"            // Unknown status

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
	qbt.option = append(qbt.option, qbapi.WithTimeout(time.Duration(store.Config.Advanced.ClientConf.ConnectTimeoutSecond)*time.Second))
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

func (c *QBittorrent) Start(ctx context.Context) {
	c.connectFunc = func() bool {
		var err error
		c.client, err = qbapi.NewAPI(c.option...)
		if err != nil {
			zap.S().Debug(err)
			zap.S().Warnf("初始化QBittorrent客户端第%d次，失败", c.retryNum)
			return false
		}
		if err = c.client.Login(ctx); err != nil {
			zap.S().Debug(err)
			zap.S().Warnf("连接QBittorrent第%d次，失败", c.retryNum)
			return false
		}
		return true
	}
	store.WG.Add(2)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.S().Error(err)
			}
		}()
		defer store.WG.Done()
		for {
			select {
			case <-ctx.Done():
				zap.S().Info("正常退出")
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
						zap.S().Info("连接QBittorrent成功")
					}
				}
			}
		}
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.S().Error(err)
			}
		}()
		defer store.WG.Done()
		for {
			select {
			case <-ctx.Done():
				zap.S().Info("正常退出")
				return
			default:
				if c.retryNum == 0 {
					c.retryChan <- ChanRetryConnect
					// 检查是否在线，时间长
					utils.Sleep(store.Config.CheckTimeSecond, ctx)
				} else if c.retryNum <= store.Config.RetryConnectNum {
					c.retryChan <- ChanRetryConnect
					// 失败重试，时间短
					utils.Sleep(store.Config.ConnectTimeoutSecond, ctx)
				} else {
					// 超过重试次数，不在频繁重试
					c.retryNum = 0
				}
			}
		}
	}()

}

// checkError
//  @Description: 检查错误，返回是否需要结束流程
//  @receiver *QBittorrent
//  @param err error
//  @return bool
//
func (c *QBittorrent) checkError(err error, fun string) bool {
	if err == nil {
		return false
	}
	if qerror, ok := err.(*qbapi.QError); ok && qerror.Code() == -10004 {
		zap.S().Debugf("[%s] 请求失败，等待客户端响应...", fun)
		c.retryNum = 1
		c.connected = false
		c.retryChan <- ChanRetryConnect
	} else {
		zap.S().Debugf("[%s] %v", fun, err)
		zap.S().Warn("请求QBittorrent接口失败")
	}
	return true
}

func (c *QBittorrent) Version() string {
	if !c.connected {
		return ""
	}
	clientResp, err := c.client.GetApplicationVersion(context.Background(), &qbapi.GetApplicationVersionReq{})
	if c.checkError(err, "Version 1") {
		return ""
	}
	apiResp, err := c.client.GetAPIVersion(context.Background(), &qbapi.GetAPIVersionReq{})
	if c.checkError(err, "Version 2") {
		return ""
	}
	return fmt.Sprintf("Client: %s, API: %s", clientResp.Version, apiResp.Version)
}

func (c *QBittorrent) Preferences() *models.Preferences {
	if !c.connected {
		return nil
	}
	resp, err := c.client.GetApplicationPreferences(context.Background(), &qbapi.GetApplicationPreferencesReq{})
	if c.checkError(err, "Preferences") {
		return nil
	}
	retn := &models.Preferences{}
	utils.ConvertModel(resp, retn)
	return retn
}

func (c *QBittorrent) SetDefaultPreferences() {
	if !c.connected {
		return
	}
	opt := "NoSubfolder"
	pref := &qbapi.SetApplicationPreferencesReq{
		TorrentContentLayout: &opt,
	}
	_, err := c.client.SetApplicationPreferences(context.Background(), pref)
	if c.checkError(err, "SetDefaultPreferences") {
		return
	}
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
	if c.checkError(err, "List") {
		return nil
	}
	retn := make([]*models.TorrentItem, len(listResp.Items))
	for i, _ := range retn {
		retn[i] = &models.TorrentItem{}
		utils.ConvertModel(listResp.Items[i], retn[i])
	}
	return retn
}

func (c *QBittorrent) Rename(opt *models.ClientRenameOptions) {
	if !c.connected {
		return
	}
	_, err := c.client.RenameFile(context.Background(), &qbapi.RenameFileReq{
		Hash:    opt.Hash,
		OldPath: opt.OldPath,
		NewPath: opt.NewPath,
	})
	if c.checkError(err, "Rename") {
		return
	}
}

func (c *QBittorrent) Add(opt *models.ClientAddOptions) {
	if !c.connected {
		return
	}
	_, err := c.client.AddNewLink(context.Background(), &qbapi.AddNewLinkReq{
		Url: opt.Urls,
		Meta: &qbapi.AddTorrentMeta{
			Savepath:         &opt.SavePath,
			Category:         &opt.Category,
			Tags:             opt.Tag,
			SeedingTimeLimit: &opt.SeedingTime,
			Rename:           &opt.Rename,
		},
	})
	if c.checkError(err, "Add") {
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
	if c.checkError(err, "Delete") {
		return
	}
}

func (c *QBittorrent) Get(opt *models.ClientGetOptions) *models.TorrentItem {
	if !c.connected {
		return nil
	}
	resp, err := c.client.GetTorrentGenericProperties(context.Background(), &qbapi.GetTorrentGenericPropertiesReq{
		Hash: opt.Hash,
	})
	if c.checkError(err, "Get") {
		return nil
	}

	retn := &models.TorrentItem{}
	utils.ConvertModel(resp.Property, retn)
	return retn
}

func (c *QBittorrent) GetContent(opt *models.ClientGetOptions) []*models.TorrentContentItem {
	if !c.connected {
		return nil
	}
	contents, err := c.client.GetTorrentContents(context.Background(), &qbapi.GetTorrentContentsReq{
		Hash: opt.Hash,
	})
	if c.checkError(err, "GetContent") {
		return nil
	}
	retn := make([]*models.TorrentContentItem, len(contents.Contents))
	for i, _ := range retn {
		retn[i] = &models.TorrentContentItem{}
		utils.ConvertModel(contents.Contents[i], retn[i])
	}
	return retn
}
