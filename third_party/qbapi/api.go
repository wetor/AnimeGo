package qbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

type QBAPI struct {
	c      *Config
	client *http.Client
}

func NewAPI(opts ...Option) (*QBAPI, error) {
	c := &Config{
		Timeout: 5 * time.Second,
	}
	for _, opt := range opts {
		opt(c)
	}
	if strings.HasSuffix(c.Host, "/") {
		c.Host = strings.TrimRight(c.Host, "/")
	}
	if len(c.Host) == 0 || len(c.Username) == 0 || len(c.Password) == 0 {
		return nil, NewMsgError(ErrParams, "params err")
	}
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: c.Timeout,
		Jar:     jar,
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	return &QBAPI{c: c, client: client}, nil
}

func (q *QBAPI) get(ctx context.Context, path string, req map[string]string) (*http.Response, error) {
	uri := q.buildURI(path)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, NewError(ErrInternal, err)
	}
	if len(req) != 0 {
		query := httpReq.URL.Query()
		for k, v := range req {
			query.Add(k, v)
		}
		httpReq.URL.RawQuery = query.Encode()
	}

	resp, err := q.client.Do(httpReq)
	if err != nil {
		return nil, NewError(ErrNetwork, err)
	}
	return resp, nil
}

func (q *QBAPI) post(ctx context.Context, path string, values map[string]string) (*http.Response, error) {
	uri := q.buildURI(path)

	var reader io.Reader
	if len(values) != 0 {
		form := url.Values{}
		for k, v := range values {
			form.Add(k, v)
		}
		reader = bytes.NewReader([]byte(form.Encode()))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, reader)
	if err != nil {
		return nil, NewError(ErrInternal, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, NewError(ErrNetwork, err)
	}
	return resp, nil
}

func (q *QBAPI) postMultiPart(ctx context.Context, path string, buffer *bytes.Buffer, part *multipart.Writer) (*http.Response, error) {
	uri := q.buildURI(path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", part.FormDataContentType())
	rsp, err := q.client.Do(req)
	if err != nil {
		return nil, NewError(ErrNetwork, err)
	}
	return rsp, nil
}

func (q *QBAPI) struct2map(req interface{}) (map[string]string, error) {
	return ToMap(req, "json")
}

func (q *QBAPI) getWithDecoder(ctx context.Context, path string, req interface{}, rsp interface{}, decoder Decoder) error {
	mp, err := q.struct2map(req)
	if err != nil {
		return NewError(ErrMarsal, err)
	}
	httpRsp, err := q.get(ctx, path, mp)
	if err != nil {
		return err
	}
	defer httpRsp.Body.Close()
	if httpRsp.StatusCode != http.StatusOK {
		return NewError(ErrStatusCode, NewStatusCodeErr(httpRsp.StatusCode))
	}
	if rsp == nil {
		return nil
	}
	data, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return NewError(ErrNetwork, err)
	}
	if err := decoder(data, rsp); err != nil {
		return NewError(ErrUnmarsal, err)
	}
	return nil
}

func (q *QBAPI) buildURI(path string) string {
	return q.c.Host + path
}

func (q *QBAPI) postWithDecoder(ctx context.Context, path string, req interface{}, rsp interface{}, decoder Decoder) error {
	mp, err := q.struct2map(req)
	if err != nil {
		return NewError(ErrMarsal, err)
	}
	httpRsp, err := q.post(ctx, path, mp)
	if err != nil {
		return err
	}
	defer httpRsp.Body.Close()
	if httpRsp.StatusCode != http.StatusOK {
		return NewError(ErrStatusCode, NewStatusCodeErr(httpRsp.StatusCode))
	}
	if rsp == nil {
		return nil
	}

	data, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return NewError(ErrNetwork, err)
	}
	if err := decoder(data, rsp); err != nil {
		return NewError(ErrUnmarsal, err)
	}
	return nil
}

//Login /api/v2/auth/login
func (q *QBAPI) Login(ctx context.Context) error {
	req := &LoginReq{Username: q.c.Username, Password: q.c.Password}
	var rsp string
	err := q.postWithDecoder(ctx, apiLogin, req, &rsp, StrDec)
	if err != nil {
		return err
	}
	if !strings.Contains(strings.ToLower(rsp), "ok") {
		return NewError(ErrLogin, fmt.Errorf("login fail:%s", rsp))
	}
	return nil
}

//GetApplicationVersion /api/v2/app/version
func (q *QBAPI) GetApplicationVersion(ctx context.Context, req *GetApplicationVersionReq) (*GetApplicationVersionRsp, error) {
	var version string
	err := q.getWithDecoder(ctx, apiGetAPPVersion, nil, &version, StrDec)
	if err != nil {
		return nil, err
	}
	return &GetApplicationVersionRsp{version}, nil
}

//GetAPIVersion /api/v2/app/webapiVersion
func (q *QBAPI) GetAPIVersion(ctx context.Context, req *GetAPIVersionReq) (*GetAPIVersionRsp, error) {
	var version string
	err := q.getWithDecoder(ctx, apiGetAPIVersion, nil, &version, StrDec)
	if err != nil {
		return nil, err
	}
	return &GetAPIVersionRsp{Version: version}, nil
}

//GetBuildInfo /api/v2/app/buildInfo
func (q *QBAPI) GetBuildInfo(ctx context.Context, req *GetBuildInfoReq) (*GetBuildInfoRsp, error) {
	rsp := &GetBuildInfoRsp{Info: &BuildInfo{}}
	if err := q.getWithDecoder(ctx, apiGetBuildInfo, req, rsp.Info, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//ShutDownAPPlication /api/v2/app/shutdown
func (q *QBAPI) ShutDownAPPlication(ctx context.Context, req *ShutdownApplicationReq) (*ShutdownApplicationRsp, error) {
	err := q.postWithDecoder(ctx, apiShutdownAPP, nil, nil, JsonDec)
	if err != nil {
		return nil, err
	}
	return &ShutdownApplicationRsp{}, nil
}

//GetApplicationPreferences /api/v2/app/preferences
func (q *QBAPI) GetApplicationPreferences(ctx context.Context, req *GetApplicationPreferencesReq) (*GetApplicationPreferencesRsp, error) {
	rsp := &GetApplicationPreferencesRsp{}
	err := q.getWithDecoder(ctx, apiGetAPPPerf, req, rsp, JsonDec)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

//SetApplicationPreferences /api/v2/app/setPreferences
func (q *QBAPI) SetApplicationPreferences(ctx context.Context, req *SetApplicationPreferencesReq) (*SetApplicationPreferencesRsp, error) {
	js, err := json.Marshal(req)
	if err != nil {
		return nil, NewError(ErrMarsal, err)
	}
	innerReq := &setApplicationPreferencesInnerReq{
		Json: string(js),
	}
	if err := q.postWithDecoder(ctx, apiSetAPPPref, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetApplicationPreferencesRsp{}, nil
}

//GetDefaultSavePath /api/v2/app/defaultSavePath
func (q *QBAPI) GetDefaultSavePath(ctx context.Context, req *GetDefaultSavePathReq) (*GetDefaultSavePathRsp, error) {
	var path string
	if err := q.getWithDecoder(ctx, apiGetDefaultSavePath, nil, &path, StrDec); err != nil {
		return nil, err
	}
	return &GetDefaultSavePathRsp{Path: path}, nil
}

//GetLog /api/v2/log/main
func (q *QBAPI) GetLog(ctx context.Context, req *GetLogReq) (*GetLogRsp, error) {
	rsp := &GetLogRsp{Items: make([]*LogItem, 0)}
	if err := q.getWithDecoder(ctx, apiGetLog, req, &rsp.Items, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetPeerLog /api/v2/log/peers
func (q *QBAPI) GetPeerLog(ctx context.Context, req *GetPeerLogReq) (*GetPeerLogRsp, error) {
	rsp := &GetPeerLogRsp{Items: make([]*PeerLogItem, 0)}
	if err := q.getWithDecoder(ctx, apiGetPeerLog, req, &rsp.Items, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetMainData /api/v2/sync/maindata
func (q *QBAPI) GetMainData(ctx context.Context, req *GetMainDataReq) (*GetMainDataRsp, error) {
	rsp := &GetMainDataRsp{}
	if err := q.getWithDecoder(ctx, apiGetMainData, req, &rsp, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetTorrentPeerData /api/v2/sync/torrentPeers
func (q *QBAPI) GetTorrentPeerData(ctx context.Context, req *GetTorrentPeerDataReq) (*GetTorrentPeerDataRsp, error) {
	rsp := &GetTorrentPeerDataRsp{Data: &TorrentPeerData{}}
	if err := q.getWithDecoder(ctx, apiGetTorrentPeerData, req, rsp.Data, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetGlobalTransferInfo /api/v2/transfer/info
func (q *QBAPI) GetGlobalTransferInfo(ctx context.Context, req *GetGlobalTransferInfoReq) (*GetGlobalTransferInfoRsp, error) {
	rsp := &GetGlobalTransferInfoRsp{Info: &GlobalTransferInfo{}}
	if err := q.getWithDecoder(ctx, apiGetGlobalTransferInfo, req, rsp.Info, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetAlternativeSpeedLimitsState /api/v2/transfer/speedLimitsMode
func (q *QBAPI) GetAlternativeSpeedLimitsState(ctx context.Context, req *GetAlternativeSpeedLimitsStateReq) (*GetAlternativeSpeedLimitsStateRsp, error) {
	var intEnabled int
	rsp := &GetAlternativeSpeedLimitsStateRsp{Enabled: true}
	if err := q.getWithDecoder(ctx, apiGetAltSpeedLimitState, req, &intEnabled, IntDec); err != nil {
		return nil, err
	}
	if intEnabled == 0 {
		rsp.Enabled = false
	}
	return rsp, nil
}

//ToggleAlternativeSpeedLimits /api/v2/transfer/toggleSpeedLimitsMode
func (q *QBAPI) ToggleAlternativeSpeedLimits(ctx context.Context, req *ToggleAlternativeSpeedLimitsReq) (*ToggleAlternativeSpeedLimitsRsp, error) {
	rsp := &ToggleAlternativeSpeedLimitsRsp{}
	if err := q.postWithDecoder(ctx, apiToggleAltSpeedLimits, req, nil, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetGlobalDownloadLimit /api/v2/transfer/downloadLimit
func (q *QBAPI) GetGlobalDownloadLimit(ctx context.Context, req *GetGlobalDownloadLimitReq) (*GetGlobalDownloadLimitRsp, error) {
	rsp := &GetGlobalDownloadLimitRsp{}
	if err := q.getWithDecoder(ctx, apiGetGlobalDownloadLimit, req, &rsp.Speed, IntDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//SetGlobalDownloadLimit /api/v2/transfer/setDownloadLimit
func (q *QBAPI) SetGlobalDownloadLimit(ctx context.Context, req *SetGlobalDownloadLimitReq) (*SetGlobalDownloadLimitRsp, error) {
	if err := q.postWithDecoder(ctx, apiSetGlobalDownloadLimit, req, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetGlobalDownloadLimitRsp{}, nil
}

//GetGlobalUploadLimit /api/v2/transfer/uploadLimit
func (q *QBAPI) GetGlobalUploadLimit(ctx context.Context, req *GetGlobalUploadLimitReq) (*GetGlobalUploadLimitRsp, error) {
	rsp := &GetGlobalUploadLimitRsp{}
	if err := q.getWithDecoder(ctx, apiGetGlobalUploadLimit, req, &rsp.Speed, IntDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//SetGlobalUploadLimit /api/v2/transfer/setUploadLimit
func (q *QBAPI) SetGlobalUploadLimit(ctx context.Context, req *SetGlobalUploadLimitReq) (*SetGlobalUploadLimitRsp, error) {
	if err := q.postWithDecoder(ctx, apiSetGlobalUploadLimit, req, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetGlobalUploadLimitRsp{}, nil
}

//BanPeers /api/v2/transfer/banPeers
func (q *QBAPI) BanPeers(ctx context.Context, req *BanPeersReq) (*BanPeersRsp, error) {
	for _, item := range req.Peers {
		if !strings.Contains(item, ":") {
			return nil, NewError(ErrParams, fmt.Errorf("invalid peer:%s", item))
		}
	}
	innerReq := &banPeersReqInner{Peers: strings.Join(req.Peers, "|")}
	if err := q.postWithDecoder(ctx, apiBanPeers, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &BanPeersRsp{}, nil
}

//GetTorrentList /api/v2/torrents/info
func (q *QBAPI) GetTorrentList(ctx context.Context, req *GetTorrentListReq) (*GetTorrentListRsp, error) {
	rsp := &GetTorrentListRsp{Items: make([]*TorrentListItem, 0)}
	if err := q.getWithDecoder(ctx, apiGetTorrentList, req, &rsp.Items, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetTorrentGenericProperties /api/v2/torrents/properties
func (q *QBAPI) GetTorrentGenericProperties(ctx context.Context, req *GetTorrentGenericPropertiesReq) (*GetTorrentGenericPropertiesRsp, error) {
	rsp := &GetTorrentGenericPropertiesRsp{Property: &TorrentGenericProperty{}}
	if err := q.getWithDecoder(ctx, apiGetTorrentGenericProp, req, rsp.Property, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetTorrentTrackers /api/v2/torrents/trackers
func (q *QBAPI) GetTorrentTrackers(ctx context.Context, req *GetTorrentTrackersReq) (*GetTorrentTrackersRsp, error) {
	rsp := &GetTorrentTrackersRsp{Trackers: make([]*TorrentTrackerItem, 0)}
	if err := q.getWithDecoder(ctx, apiGetTorrentTrackers, req, &rsp.Trackers, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetTorrentWebSeeds /api/v2/torrents/webseeds
func (q *QBAPI) GetTorrentWebSeeds(ctx context.Context, req *GetTorrentWebSeedsReq) (*GetTorrentWebSeedsRsp, error) {
	rsp := &GetTorrentWebSeedsRsp{WebSeeds: make([]*TorrentWebSeedItem, 0)}
	if err := q.getWithDecoder(ctx, apiGetTorrentWebSeeds, req, &rsp.WebSeeds, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetTorrentContents /api/v2/torrents/files
func (q *QBAPI) GetTorrentContents(ctx context.Context, req *GetTorrentContentsReq) (*GetTorrentContentsRsp, error) {
	rsp := &GetTorrentContentsRsp{Contents: make([]*TorrentContentItem, 0)}

	innerReq := &getTorrentContentsInnerReq{Hash: req.Hash}
	if len(req.Index) > 0 {
		indexes := strings.Join(req.Index, "|")
		innerReq.Indexes = &indexes
	}
	if err := q.getWithDecoder(ctx, apiGetTorrentContents, innerReq, &rsp.Contents, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetTorrentPiecesStates /api/v2/torrents/pieceStates
func (q *QBAPI) GetTorrentPiecesStates(ctx context.Context, req *GetTorrentPiecesStatesReq) (*GetTorrentPiecesStatesRsp, error) {
	rsp := &GetTorrentPiecesStatesRsp{States: make([]int, 0)}
	if err := q.getWithDecoder(ctx, apiGetTorrentPiecesStates, req, &rsp.States, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetTorrentPiecesHashes /api/v2/torrents/pieceHashes
func (q *QBAPI) GetTorrentPiecesHashes(ctx context.Context, req *GetTorrentPiecesHashesReq) (*GetTorrentPiecesHashesRsp, error) {
	rsp := &GetTorrentPiecesHashesRsp{Hashes: make([]string, 0)}
	if err := q.getWithDecoder(ctx, apiGetTorrentPiecesHashes, req, &rsp.Hashes, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//PauseTorrents /api/v2/torrents/pause
func (q *QBAPI) PauseTorrents(ctx context.Context, req *PauseTorrentsReq) (*PauseTorrentsRsp, error) {
	if len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("non hashes found"))
	}
	hashes := strings.Join(req.Hash, "|")
	if err := q.getWithDecoder(ctx, apiPauseTorrents, &pauseTorrentsInnerReq{Hashes: hashes}, nil, JsonDec); err != nil {
		return nil, err
	}
	return &PauseTorrentsRsp{}, nil
}

//ResumeTorrents /api/v2/torrents/resume
func (q *QBAPI) ResumeTorrents(ctx context.Context, req *ResumeTorrentsReq) (*ResumeTorrentsRsp, error) {
	if len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("non hashes found"))
	}
	hashes := strings.Join(req.Hash, "|")
	if err := q.getWithDecoder(ctx, apiResumeTorrents, &resumeTorrentsInnerReq{Hashes: hashes}, nil, JsonDec); err != nil {
		return nil, err
	}
	return &ResumeTorrentsRsp{}, nil
}

//DeleteTorrents /api/v2/torrents/delete
func (q *QBAPI) DeleteTorrents(ctx context.Context, req *DeleteTorrentsReq) (*DeleteTorrentsRsp, error) {
	if !req.IsDeleteAll && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("non hashes found"))
	}
	innerReq := &deleteTorrentsInnerReq{
		DeleteFiles: req.IsDeleteFile,
	}
	if req.IsDeleteAll {
		innerReq.Hashes = "all"
	} else {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.getWithDecoder(ctx, apiDeleteTorrents, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &DeleteTorrentsRsp{}, nil
}

//RecheckTorrents /api/v2/torrents/recheck
func (q *QBAPI) RecheckTorrents(ctx context.Context, req *RecheckTorrentsReq) (*RecheckTorrentsRsp, error) {
	if !req.IsRecheckAll && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("non hashes found"))
	}
	innerReq := &recheckTorrentsInnerReq{}
	if req.IsRecheckAll {
		innerReq.Hashes = "all"
	} else {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.getWithDecoder(ctx, apiRecheckTorrents, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &RecheckTorrentsRsp{}, nil
}

//ReannounceTorrents /api/v2/torrents/reannounce
func (q *QBAPI) ReannounceTorrents(ctx context.Context, req *ReannounceTorrentsReq) (*ReannounceTorrentsRsp, error) {
	if !req.IsReannounceAll && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("non hashes found"))
	}
	innerReq := &reannounceTorrentsInnerReq{}
	if req.IsReannounceAll {
		innerReq.Hashes = "all"
	} else {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.getWithDecoder(ctx, apiReannounceTorrents, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &ReannounceTorrentsRsp{}, nil
}

func (q *QBAPI) writeProperties(writer *multipart.Writer, meta *AddTorrentMeta) error {
	if meta == nil {
		return nil
	}
	cb := func(key string, value interface{}) error {
		w, err := writer.CreateFormField(key)
		if err != nil {
			return err
		}
		v := fmt.Sprintf("%v", value)
		if _, err := w.Write([]byte(v)); err != nil {
			return err
		}
		return nil
	}
	if err := IterStruct(meta, "json", cb); err != nil {
		return err
	}
	return nil
}

//AddNewTorrent /api/v2/torrents/add
func (q *QBAPI) AddNewTorrent(ctx context.Context, req *AddNewTorrentReq) (*AddNewTorrentRsp, error) {
	if len(req.File) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("params err"))
	}

	buffer := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buffer)
	for _, torrent := range req.File {
		w, err := writer.CreateFormFile("torrents", filepath.Base(torrent))
		if err != nil {
			return nil, NewError(ErrFile, err)
		}
		data, err := ioutil.ReadFile(torrent)
		if err != nil {
			return nil, NewError(ErrFile, err)
		}
		if _, err := w.Write(data); err != nil {
			return nil, NewError(ErrFile, err)
		}
	}
	if err := q.writeProperties(writer, req.Meta); err != nil {
		return nil, NewError(ErrUnknown, err)
	}
	httpRsp, err := q.postMultiPart(ctx, apiAddNewTorrent, buffer, writer)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()
	if httpRsp.StatusCode != http.StatusOK {
		return nil, NewError(ErrStatusCode, NewStatusCodeErr(httpRsp.StatusCode))
	}
	return &AddNewTorrentRsp{}, nil
}

//AddNewLink /api/v2/torrents/add
func (q *QBAPI) AddNewLink(ctx context.Context, req *AddNewLinkReq) (*AddNewLinkRsp, error) {
	if len(req.Url) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("params err"))
	}
	buffer := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buffer)
	w, err := writer.CreateFormField("urls")
	if err != nil {
		return nil, NewError(ErrInternal, err)
	}
	if _, err := w.Write([]byte(strings.Join(req.Url, "\r\n"))); err != nil {
		return nil, NewError(ErrInternal, err)
	}
	if err := q.writeProperties(writer, req.Meta); err != nil {
		return nil, NewError(ErrUnknown, err)
	}
	httpRsp, err := q.postMultiPart(ctx, apiAddNewTorrent, buffer, writer)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()
	if httpRsp.StatusCode != http.StatusOK {
		return nil, NewError(ErrStatusCode, NewStatusCodeErr(httpRsp.StatusCode))
	}
	return &AddNewLinkRsp{}, nil
}

//AddTrackersToTorrent /api/v2/torrents/addTrackers
func (q *QBAPI) AddTrackersToTorrent(ctx context.Context, req *AddTrackersToTorrentReq) (*AddTrackersToTorrentRsp, error) {
	if len(req.Url) == 0 || len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &addTrackersToTorrentInnerReq{Urls: strings.Join(req.Url, "\n"), Hash: req.Hash}
	err := q.postWithDecoder(ctx, apiAddTrackersToTorrent, innerReq, nil, JsonDec)
	if err == nil {
		return &AddTrackersToTorrentRsp{}, nil
	}
	return nil, err
}

/*
400	newUrl is not a valid URL
404	Torrent hash was not found
409	newUrl already exists for the torrent
409	origUrl was not found
200	All other scenarios
*/
//EditTrackers /api/v2/torrents/editTracker
func (q *QBAPI) EditTrackers(ctx context.Context, req *EditTrackersReq) (*EditTrackersRsp, error) {
	if len(req.Hash) == 0 || len(req.NewUrl) == 0 || len(req.OrigUrl) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	if err := q.postWithDecoder(ctx, apiEditTrackers, req, nil, JsonDec); err != nil {
		return nil, err
	}
	return &EditTrackersRsp{}, nil
}

/*
404	Torrent hash was not found
409	All urls were not found
200	All other scenarios
*/
//RemoveTrackers /api/v2/torrents/removeTrackers
func (q *QBAPI) RemoveTrackers(ctx context.Context, req *RemoveTrackersReq) (*RemoveTrackersRsp, error) {
	if len(req.Hash) == 0 || len(req.Url) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &removeTrackersInnerReq{
		Hash: req.Hash,
		Urls: strings.Join(req.Url, "|"),
	}
	if err := q.postWithDecoder(ctx, apiRemoveTrackers, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &RemoveTrackersRsp{}, nil
}

/*
400	None of the supplied peers are valid
200	All other scenarios
*/
//AddPeers /api/v2/torrents/addPeers
func (q *QBAPI) AddPeers(ctx context.Context, req *AddPeersReq) (*AddPeersRsp, error) {
	if len(req.Hash) == 0 || len(req.Peer) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &addPeersInnerReq{
		Hashes: strings.Join(req.Hash, "|"),
		Peers:  strings.Join(req.Peer, "|"),
	}
	if err := q.postWithDecoder(ctx, apiAddPeers, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &AddPeersRsp{}, nil
}

/*
409	Torrent queueing is not enabled
200	All other scenarios
*/
//IncreaseTorrentPriority /api/v2/torrents/increasePrio
func (q *QBAPI) IncreaseTorrentPriority(ctx context.Context, req *IncreaseTorrentPriorityReq) (*IncreaseTorrentPriorityRsp, error) {
	if !req.IsIncreaseAllTorrent && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := increaseTorrentPriorityInnerReq{
		Hashes: "all",
	}
	if !req.IsIncreaseAllTorrent {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiIncreaseTorrentPriority, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &IncreaseTorrentPriorityRsp{}, nil
}

/*
409	Torrent queueing is not enabled
200	All other scenarios
*/
//DecreaseTorrentPriority /api/v2/torrents/decreasePrio
func (q *QBAPI) DecreaseTorrentPriority(ctx context.Context, req *DecreaseTorrentPriorityReq) (*DecreaseTorrentPriorityRsp, error) {
	if !req.IsDecreaseAllTorrent && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := decreaseTorrentPriorityInnerReq{
		Hashes: "all",
	}
	if !req.IsDecreaseAllTorrent {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiDecreaseTorrentPriority, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &DecreaseTorrentPriorityRsp{}, nil
}

//MaximalTorrentPriority /api/v2/torrents/topPrio
func (q *QBAPI) MaximalTorrentPriority(ctx context.Context, req *MaximalTorrentPriorityReq) (*MaximalTorrentPriorityRsp, error) {
	if !req.IsMaximalAllTorrent && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := maximalTorrentPriorityInnerReq{
		Hashes: "all",
	}
	if !req.IsMaximalAllTorrent {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiMaximalTorrentPriority, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &MaximalTorrentPriorityRsp{}, nil
}

//MinimalTorrentPriority /api/v2/torrents/bottomPrio
func (q *QBAPI) MinimalTorrentPriority(ctx context.Context, req *MinimalTorrentPriorityReq) (*MinimalTorrentPriorityRsp, error) {
	if !req.IsMinimalAllTorrent && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := minimalTorrentPriorityInnerReq{
		Hashes: "all",
	}
	if !req.IsMinimalAllTorrent {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiMinimalTorrentPriority, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &MinimalTorrentPriorityRsp{}, nil
}

/*
400	Priority is invalid
400	At least one file id is not a valid integer
404	Torrent hash was not found
409	Torrent metadata hasn't downloaded yet
409	At least one file id was not found
200	All other scenarios
*/
//SetFilePriority /api/v2/torrents/filePrio
func (q *QBAPI) SetFilePriority(ctx context.Context, req *SetFilePriorityReq) (*SetFilePriorityRsp, error) {
	if len(req.Hash) == 0 || len(req.Id) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &setFilePriorityInnerReq{
		Hash:     req.Hash,
		Id:       strings.Join(req.Id, "|"),
		Priority: int(req.Priority),
	}
	if err := q.postWithDecoder(ctx, apiSetFilePriority, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetFilePriorityRsp{}, nil
}

//GetTorrentDownloadLimit /api/v2/torrents/downloadLimit
func (q *QBAPI) GetTorrentDownloadLimit(ctx context.Context, req *GetTorrentDownloadLimitReq) (*GetTorrentDownloadLimitRsp, error) {
	if !req.IsGetAll && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &getTorrentDownloadLimitInnerReq{Hashes: "all"}
	if !req.IsGetAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	rsp := &GetTorrentDownloadLimitRsp{SpeedMap: make(map[string]int)}
	if err := q.postWithDecoder(ctx, apiGetTorrentDownloadLimit, innerReq, &rsp.SpeedMap, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//SetTorrentDownloadLimit /api/v2/torrents/setDownloadLimit
func (q *QBAPI) SetTorrentDownloadLimit(ctx context.Context, req *SetTorrentDownloadLimitReq) (*SetTorrentDownloadLimitRsp, error) {
	if !req.IsSetAll && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &setTorrentDownloadLimitInnerReq{
		Hashes: "all",
		Speed:  req.Speed,
	}
	if !req.IsSetAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiSetTorrentDownloadLimit, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetTorrentDownloadLimitRsp{}, nil
}

//SetTorrentShareLimit /api/v2/torrents/setShareLimits
func (q *QBAPI) SetTorrentShareLimit(ctx context.Context, req *SetTorrentShareLimitReq) (*SetTorrentShareLimitRsp, error) {
	if (!req.IsSetAll && len(req.Hash) == 0) || req.RatioLimit == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &setTorrentShareLimitInnerReq{
		Hashes:           "all",
		SeedingTimeLimit: req.SeedingTimeLimit,
		RatioLimit:       req.RatioLimit,
	}
	if !req.IsSetAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiSetTorrentShareLimit, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetTorrentShareLimitRsp{}, nil
}

//GetTorrentUploadLimit /api/v2/torrents/uploadLimit
func (q *QBAPI) GetTorrentUploadLimit(ctx context.Context, req *GetTorrentUploadLimitReq) (*GetTorrentUploadLimitRsp, error) {
	if !req.IsGetAll && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &getTorrentUploadLimitInnerReq{
		Hashes: "all",
	}
	if !req.IsGetAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	rsp := &GetTorrentUploadLimitRsp{SpeedMap: make(map[string]int)}
	if err := q.postWithDecoder(ctx, apiGetTorrentUploadLimit, innerReq, &rsp.SpeedMap, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//SetTorrentUploadLimit /api/v2/torrents/setUploadLimit
func (q *QBAPI) SetTorrentUploadLimit(ctx context.Context, req *SetTorrentUploadLimitReq) (*SetTorrentUploadLimitRsp, error) {
	if !req.IsSetAll && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &setTorrentUploadLimitInnerReq{
		Hashes: "all",
		Speed:  req.Speed,
	}
	if !req.IsSetAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiSetTorrentUploadLimit, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetTorrentUploadLimitRsp{}, nil
}

/*
400	Save path is empty
403	User does not have write access to directory
409	Unable to create save path directory
200	All other scenarios
*/
//SetTorrentLocation /api/v2/torrents/setLocation
func (q *QBAPI) SetTorrentLocation(ctx context.Context, req *SetTorrentLocationReq) (*SetTorrentLocationRsp, error) {
	if !req.IsSetAll && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &setTorrentLocationInnerReq{
		Hashes:   "all",
		Location: req.Location,
	}
	if !req.IsSetAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiSetTorrentLocation, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetTorrentLocationRsp{}, nil
}

/*
404	Torrent hash is invalid
409	Torrent name is empty
200	All other scenarios
*/
//SetTorrentName /api/v2/torrents/rename
func (q *QBAPI) SetTorrentName(ctx context.Context, req *SetTorrentNameReq) (*SetTorrentNameRsp, error) {
	if len(req.Hash) == 0 || len(req.Name) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	if err := q.postWithDecoder(ctx, apiSetTorrentName, req, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetTorrentNameRsp{}, nil
}

/*
409	Category name does not exist
200	All other scenarios
*/
//SetTorrentCategory /api/v2/torrents/setCategory
func (q *QBAPI) SetTorrentCategory(ctx context.Context, req *SetTorrentCategoryReq) (*SetTorrentCategoryRsp, error) {
	if len(req.Category) == 0 || (!req.IsSetAll && len(req.Hash) == 0) {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &setTorrentCategoryInnerReq{
		Hashes:   "all",
		Category: req.Category,
	}
	if !req.IsSetAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiSetTorrentCategory, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetTorrentCategoryRsp{}, nil
}

//GetAllCategories /api/v2/torrents/categories
func (q *QBAPI) GetAllCategories(ctx context.Context, req *GetAllCategoriesReq) (*GetAllCategoriesRsp, error) {
	rsp := &GetAllCategoriesRsp{Categories: make(map[string]*CategoryInfo)}
	if err := q.getWithDecoder(ctx, apiGetAllCategories, req, &rsp.Categories, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

/*
400	Category name is empty
409	Category name is invalid
200	All other scenarios
*/
//AddNewCategory /api/v2/torrents/createCategory
func (q *QBAPI) AddNewCategory(ctx context.Context, req *AddNewCategoryReq) (*AddNewCategoryRsp, error) {
	if len(req.Category) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	if err := q.postWithDecoder(ctx, apiAddNewCategory, req, nil, JsonDec); err != nil {
		return nil, err
	}
	return &AddNewCategoryRsp{}, nil
}

/*
400	Category name is empty
409	Category editing failed
200	All other scenarios
*/
//EditCategory /api/v2/torrents/editCategory
func (q *QBAPI) EditCategory(ctx context.Context, req *EditCategoryReq) (*EditCategoryRsp, error) {
	if len(req.Category) == 0 || len(req.SavePath) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	if err := q.postWithDecoder(ctx, apiEditCategory, req, nil, JsonDec); err != nil {
		return nil, err
	}
	return &EditCategoryRsp{}, nil
}

//RemoveCategories /api/v2/torrents/removeCategories
func (q *QBAPI) RemoveCategories(ctx context.Context, req *RemoveCategoriesReq) (*RemoveCategoriesRsp, error) {
	if len(req.Category) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &removeCategoriesInnerReq{
		Categories: strings.Join(req.Category, "\n"),
	}
	if err := q.postWithDecoder(ctx, apiRemoveCategories, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &RemoveCategoriesRsp{}, nil
}

//AddTorrentTags /api/v2/torrents/addTags
func (q *QBAPI) AddTorrentTags(ctx context.Context, req *AddTorrentTagsReq) (*AddTorrentTagsRsp, error) {
	if len(req.Tag) == 0 || (!req.IsAddAll && len(req.Hash) == 0) {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &addTorrentTagsInnerReq{
		Tags:   strings.Join(req.Tag, ","),
		Hashes: "all",
	}
	if !req.IsAddAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiAddTorrentTags, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &AddTorrentTagsRsp{}, nil
}

//RemoveTorrentTags /api/v2/torrents/removeTags
func (q *QBAPI) RemoveTorrentTags(ctx context.Context, req *RemoveTorrentTagsReq) (*RemoveTorrentTagsRsp, error) {
	if len(req.Tag) == 0 || (!req.IsRemoveAll && len(req.Hash) == 0) {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &removeTorrentTagsInnerReq{
		Tags:   strings.Join(req.Tag, ","),
		Hashes: "all",
	}
	if !req.IsRemoveAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiRemoveTorrentTags, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &RemoveTorrentTagsRsp{}, nil
}

//GetAllTags /api/v2/torrents/tags
func (q *QBAPI) GetAllTags(ctx context.Context, req *GetAllTagsReq) (*GetAllTagsRsp, error) {
	rsp := &GetAllTagsRsp{Tags: make([]string, 0)}
	if err := q.getWithDecoder(ctx, apiGetAllTags, req, &rsp.Tags, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//CreateTags /api/v2/torrents/createTags
func (q *QBAPI) CreateTags(ctx context.Context, req *CreateTagsReq) (*CreateTagsRsp, error) {
	if len(req.Tag) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := createTagsInnerReq{
		Tags: strings.Join(req.Tag, ","),
	}
	if err := q.postWithDecoder(ctx, apiCreateTags, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &CreateTagsRsp{}, nil
}

//DeleteTags /api/v2/torrents/deleteTags
func (q *QBAPI) DeleteTags(ctx context.Context, req *DeleteTagsReq) (*DeleteTagsRsp, error) {
	if len(req.Tag) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := deleteTagsInnerReq{
		Tags: strings.Join(req.Tag, ","),
	}
	if err := q.postWithDecoder(ctx, apiDeleteTags, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &DeleteTagsRsp{}, nil

}

//SetAutomaticTorrentManagement /api/v2/torrents/setAutoManagement
func (q *QBAPI) SetAutomaticTorrentManagement(ctx context.Context, req *SetAutomaticTorrentManagementReq) (*SetAutomaticTorrentManagementRsp, error) {
	if !req.IsSetAll && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &setAutomaticTorrentManagementInnerReq{
		Hashes: "all",
		Enable: req.Enable,
	}
	if !req.IsSetAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiSetAutomaticTorrentManagement, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetAutomaticTorrentManagementRsp{}, nil
}

//ToggleSequentialDownload /api/v2/torrents/toggleSequentialDownload
func (q *QBAPI) ToggleSequentialDownload(ctx context.Context, req *ToggleSequentialDownloadReq) (*ToggleSequentialDownloadRsp, error) {
	if !req.IsSetAll && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &toggleSequentialDownloadInnerReq{
		Hashes: "all",
	}
	if !req.IsSetAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiToggleSequentialDownload, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &ToggleSequentialDownloadRsp{}, nil
}

//SetFirstOrLastPiecePriority /api/v2/torrents/toggleFirstLastPiecePrio
func (q *QBAPI) SetFirstOrLastPiecePriority(ctx context.Context, req *SetFirstOrLastPiecePriorityReq) (*SetFirstOrLastPiecePriorityRsp, error) {
	if !req.IsSetAll && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &setFirstOrLastPiecePriorityInnerReq{
		Hashes: "all",
	}
	if !req.IsSetAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiSetFirstOrLastPiecePriority, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetFirstOrLastPiecePriorityRsp{}, nil
}

//SetForceStart /api/v2/torrents/setForceStart
func (q *QBAPI) SetForceStart(ctx context.Context, req *SetForceStartReq) (*SetForceStartRsp, error) {
	if !req.IsSetAll && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &setForceStartInnerReq{
		Hashes: "all",
		Value:  req.Value,
	}
	if !req.IsSetAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiSetForceStart, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetForceStartRsp{}, nil
}

//SetSuperSeeding /api/v2/torrents/setSuperSeeding
func (q *QBAPI) SetSuperSeeding(ctx context.Context, req *SetSuperSeedingReq) (*SetSuperSeedingRsp, error) {
	if !req.IsSetAll && len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	innerReq := &setSuperSeedingInnerReq{
		Hashes: "all",
		Value:  req.Value,
	}
	if !req.IsSetAll {
		innerReq.Hashes = strings.Join(req.Hash, "|")
	}
	if err := q.postWithDecoder(ctx, apiSetSuperSeeding, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetSuperSeedingRsp{}, nil
}

/*
400	Missing newPath parameter
409	Invalid newPath or oldPath, or newPath already in use
200	All other scenarios
*/
//RenameFile /api/v2/torrents/renameFile
func (q *QBAPI) RenameFile(ctx context.Context, req *RenameFileReq) (*RenameFileRsp, error) {
	if len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	if err := q.postWithDecoder(ctx, apiRenameFile, req, nil, JsonDec); err != nil {
		return nil, err
	}
	return &RenameFileRsp{}, nil
}

//RenameFolder /api/v2/torrents/renameFolder
func (q *QBAPI) RenameFolder(ctx context.Context, req *RenameFolderReq) (*RenameFolderRsp, error) {
	if len(req.Hash) == 0 {
		return nil, NewError(ErrParams, fmt.Errorf("invalid params"))
	}
	if err := q.postWithDecoder(ctx, apiRenameFolder, req, nil, JsonDec); err != nil {
		return nil, err
	}
	return &RenameFolderRsp{}, nil
}
