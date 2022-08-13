package qbittorent

import (
	"GoBangumi/internal/models"
	"GoBangumi/store"
	"GoBangumi/third_party/qbapi"
	"GoBangumi/utils"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"time"
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
)

type QBittorrent struct {
	client     *qbapi.QBAPI
	apiVersion string
}

func NewQBittorrent(url, username, password string) *QBittorrent {
	var opts []qbapi.Option
	var err error
	var client *qbapi.QBAPI
	opts = append(opts, qbapi.WithAuth(username, password))
	opts = append(opts, qbapi.WithHost(url))
	opts = append(opts, qbapi.WithTimeout(time.Duration(store.Config.Advanced.ClientConf.ConnectTimeoutSecond)*time.Second))
	connectClient := func() bool {
		client, err = qbapi.NewAPI(opts...)
		if err != nil {
			zap.S().Warn(err)
			return false
		}
		if err = client.Login(context.Background()); err != nil {
			zap.S().Warn(err)
			return false
		}
		return true
	}

	retryNum := store.Config.Advanced.ClientConf.RetryConnectNum
	if retryNum == 0 {
		// 无限重试
		for i := 1; ; i++ {
			if connectClient() {
				break
			}
			zap.S().Infof("第%d次连接客户端失败...重新尝试连接", i)
		}
	} else {
		// 重试指定次数
		for i := 1; i <= retryNum; i++ {
			if connectClient() {
				break
			}
			zap.S().Infof("第%d次连接客户端失败...剩余%d次尝试连接", i, retryNum-i)
		}
	}

	qbt := &QBittorrent{
		client: client,
	}
	qbt.SetDefaultPreferences()
	zap.S().Infof("qBittorrent Version: %s", qbt.Version())
	return qbt
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
		// TODO: 添加下载任务后一段时间会无法获取列表
		// context deadline exceeded (Client.Timeout exceeded while awaiting headers)
		zap.S().Debugf("[%s] 请求失败，等待客户端响应...", fun)
	} else {
		zap.S().Warnf("[%s] %v", fun, err)
	}
	return true
}

func (c *QBittorrent) Version() string {
	clientResp, err := c.client.GetApplicationVersion(context.Background(), &qbapi.GetApplicationVersionReq{})
	if c.checkError(err, "Version 1") {
		return ""
	}
	apiResp, err := c.client.GetAPIVersion(context.Background(), &qbapi.GetAPIVersionReq{})
	if c.checkError(err, "Version 2") {
		return ""
	}
	c.apiVersion = apiResp.Version
	return fmt.Sprintf("Client: %s, API: %s", clientResp.Version, apiResp.Version)
}

func (c *QBittorrent) Preferences() *models.Preferences {
	resp, err := c.client.GetApplicationPreferences(context.Background(), &qbapi.GetApplicationPreferencesReq{})
	if c.checkError(err, "Preferences") {
		return nil
	}
	retn := &models.Preferences{}
	utils.ConvertModel(resp, retn)
	return retn
}

func (c *QBittorrent) SetDefaultPreferences() {
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
	_, err := c.client.DeleteTorrents(context.Background(), &qbapi.DeleteTorrentsReq{
		IsDeleteFile: opt.DeleteFile,
		Hash:         opt.Hash,
	})
	if c.checkError(err, "Delete") {
		return
	}
}

func (c *QBittorrent) Get(opt *models.ClientGetOptions) *models.TorrentItem {
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
