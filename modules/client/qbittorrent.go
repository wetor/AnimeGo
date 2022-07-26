package client

import (
	"GoBangumi/models"
	"GoBangumi/utils"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/xxxsen/qbapi"
	"golang.org/x/net/context"
	_ "unsafe"
)

//go:linkname getWithDecoder github.com/xxxsen/qbapi.(*QBAPI).getWithDecoder
func getWithDecoder(qbapi *qbapi.QBAPI, ctx context.Context, path string, req interface{}, rsp interface{}, decoder qbapi.Decoder) error

const (
	QBtStatusAll                = ""
	QBtStatusDownloading        = "downloading"
	QBtStatusSeeding            = "seeding"
	QBtStatusCompleted          = "completed"
	QBtStatusPaused             = "paused"
	QBtStatusResumed            = "resumed"
	QBtStatusActive             = "active"
	QBtStatusInactive           = "inactive"
	QBtStatusStalled            = "stalled"
	QBtStatusStalledUploading   = "stalled_uploading"
	QBtStatusStalledDownloading = "stalled_downloading"
	QBtStatusChecking           = "checking"
	QBtStatusErrored            = "errored"
)

type QBittorrent struct {
	client     *qbapi.QBAPI
	apiVersion string
}

func NewQBittorrent(url, username, password string) Client {
	var opts []qbapi.Option
	opts = append(opts, qbapi.WithAuth(username, password))
	opts = append(opts, qbapi.WithHost(url))
	client, err := qbapi.NewAPI(opts...)
	if err != nil {
		glog.Errorln(err)
		return nil
	}
	if err := client.Login(context.Background()); err != nil {
		glog.Errorln(err)
		return nil
	}
	qbt := &QBittorrent{
		client: client,
	}
	qbt.SetDefaultPreferences()

	pre := qbt.Preferences()
	fmt.Println(pre.CreateSubfolderEnabled)
	glog.V(1).Infof("qBittorrent Version: %s\n", qbt.Version())
	return qbt
}

func (c *QBittorrent) Version() string {
	clientResp, err := c.client.GetApplicationVersion(context.Background(), &qbapi.GetApplicationVersionReq{})
	if err != nil {
		glog.Errorln(err)
		return ""
	}
	apiResp, err := c.client.GetAPIVersion(context.Background(), &qbapi.GetAPIVersionReq{})
	if err != nil {
		glog.Errorln(err)
		return ""
	}
	c.apiVersion = apiResp.Version
	return fmt.Sprintf("Client: %s, API: %s", clientResp.Version, apiResp.Version)
}
func (c *QBittorrent) Preferences() *models.Preferences {
	resp, err := c.client.GetApplicationPreferences(context.Background(), &qbapi.GetApplicationPreferencesReq{})
	if err != nil {
		glog.Errorln(err)
		return nil
	}
	retn := &models.Preferences{}
	utils.ConvertModel(resp, retn)
	return retn
}

func (c *QBittorrent) SetDefaultPreferences() {
	opt := "Subfolder"
	pref := &models.SetApplicationPreferencesReq{
		TorrentContentLayout: &opt,
	}
	js, _ := json.Marshal(pref)
	innerReq := &models.SetApplicationPreferencesInnerReq{
		Json: string(js),
	}
	err := getWithDecoder(c.client, context.Background(), "/api/v2/app/setPreferences", innerReq, nil, json.Unmarshal)
	if err != nil {
		glog.Errorln(err)
		return
	}
}

func (c *QBittorrent) List(opt *models.ClientListOptions) []*models.TorrentItem {
	req := &qbapi.GetTorrentListReq{}
	if opt.Status != string(QBtStatusAll) {
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
		glog.Errorln(err)
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
	if err != nil {
		glog.Errorln(err)
	}
}
func (c *QBittorrent) Add(opt *models.ClientAddOptions) {
	_, err := c.client.AddNewLink(context.Background(), &qbapi.AddNewLinkReq{
		Url: opt.Urls,
		Meta: &qbapi.AddTorrentMeta{
			Savepath:         &opt.SavePath,
			Category:         &opt.Category,
			Tags:             opt.Tag,
			SeedingTimeLimit: &opt.SeedingTime, // 秒
			Rename:           &opt.Rename,
		},
	})
	if err != nil {
		glog.Errorln(err)
	}
}

func (c *QBittorrent) Delete(opt *models.ClientDeleteOptions) {
	_, err := c.client.DeleteTorrents(context.Background(), &qbapi.DeleteTorrentsReq{
		IsDeleteFile: opt.DeleteFile,
		Hash:         opt.Hash,
	})
	if err != nil {
		glog.Errorln(err)
	}
}

func (c *QBittorrent) Get(opt *models.ClientGetOptions) []*models.TorrentContentItem {
	contents, err := c.client.GetTorrentContents(context.Background(), &qbapi.GetTorrentContentsReq{
		Hash: opt.Hash,
	})
	if err != nil {
		glog.Errorln(err)
		return nil
	}
	retn := make([]*models.TorrentContentItem, len(contents.Contents))
	for i, _ := range retn {
		retn[i] = &models.TorrentContentItem{}
		utils.ConvertModel(contents.Contents[i], retn[i])
	}
	return retn
}
