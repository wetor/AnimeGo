package qbapi

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type cfg struct {
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Host        string   `json:"host"`
	TorrentFile []string `json:"torrent"`
	Link        []string `json:"link"`
	ValidHash   string   `json:"valid_hash"`
}

var testCfg = getCfg()
var testApi = getAPI()

func getCfg() *cfg {
	data, err := ioutil.ReadFile(".vscode/cfg.json")
	if err != nil {
		panic(err)
	}
	cf := &cfg{}
	err = json.Unmarshal(data, cf)
	if err != nil {
		panic(err)
	}
	return cf
}

func getAPI() *QBAPI {
	cf := testCfg
	var opts []Option
	opts = append(opts, WithAuth(cf.Username, cf.Password))
	opts = append(opts, WithHost(cf.Host))
	api, err := NewAPI(opts...)
	if err != nil {
		panic(err)
	}
	if err := api.Login(context.Background()); err != nil {
		panic(err)
	}
	return api
}

func TestGetTorrentList(t *testing.T) {
	limit := 10
	rsp, err := testApi.GetTorrentList(context.Background(), &GetTorrentListReq{
		Limit: &limit,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range rsp.Items {
		t.Logf("data:%+v", *item)
	}
}

func TestGetApplicationVersion(t *testing.T) {
	rsp, err := testApi.GetApplicationVersion(context.Background(), &GetApplicationVersionReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("version:%+v", rsp)
}

func TestShutdownApplication(t *testing.T) {
	_, err := testApi.ShutDownAPPlication(context.Background(), &ShutdownApplicationReq{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetApplicationPref(t *testing.T) {
	rsp, err := testApi.GetApplicationPreferences(context.Background(), &GetApplicationPreferencesReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestGetDefaultSavePath(t *testing.T) {
	rsp, err := testApi.GetDefaultSavePath(context.Background(), &GetDefaultSavePathReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestGetLog(t *testing.T) {
	rsp, err := testApi.GetLog(context.Background(), &GetLogReq{
		Normal:   true,
		Info:     true,
		Warning:  true,
		Critical: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range rsp.Items {
		t.Logf("data:%+v", *item)
	}
}

func TestGetMainData(t *testing.T) {
	rsp, err := testApi.GetMainData(context.Background(), &GetMainDataReq{
		Rid: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestGetPeerHash(t *testing.T) {
	_, err := testApi.GetTorrentPeerData(context.Background(), &GetTorrentPeerDataReq{
		Hash: "1b175f0992fe932de8de33139698c1fd26988096",
		Rid:  0,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAlternativeSpeedLimitsState(t *testing.T) {
	rsp, err := testApi.GetAlternativeSpeedLimitsState(context.Background(), &GetAlternativeSpeedLimitsStateReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestToggleAlternativeSpeedLimits(t *testing.T) {
	rsp, err := testApi.ToggleAlternativeSpeedLimits(context.Background(), &ToggleAlternativeSpeedLimitsReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestGetGlobalDownloadLimit(t *testing.T) {
	rsp, err := testApi.GetGlobalDownloadLimit(context.Background(), &GetGlobalDownloadLimitReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestSetGlobalDownloadLimit(t *testing.T) {
	_, err := testApi.SetGlobalDownloadLimit(context.Background(), &SetGlobalDownloadLimitReq{
		Speed: 50 * 1024 * 1024, //50Mb/s
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetGlobalUploadLimit(t *testing.T) {
	rsp, err := testApi.GetGlobalUploadLimit(context.Background(), &GetGlobalUploadLimitReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestSetGlobalUploadLimit(t *testing.T) {
	sp := 1.2 * 1024 * 1024
	_, err := testApi.SetGlobalUploadLimit(context.Background(), &SetGlobalUploadLimitReq{
		Speed: int(sp),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBanPeers(t *testing.T) {
	_, err := testApi.BanPeers(context.Background(), &BanPeersReq{
		[]string{
			"54.111.178.247:18635",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddTorrent(t *testing.T) {
	_, err := testApi.AddNewTorrent(context.Background(), &AddNewTorrentReq{
		File: testCfg.TorrentFile,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddLink(t *testing.T) {
	pause := true
	dlLimit := 512 * 1024
	_, err := testApi.AddNewLink(context.Background(), &AddNewLinkReq{
		Url: testCfg.Link,
		Meta: &AddTorrentMeta{
			Paused:  &pause,
			DlLimit: &dlLimit,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetBuildInfo(t *testing.T) {
	rsp, err := testApi.GetBuildInfo(context.Background(), &GetBuildInfoReq{})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp)
	assert.NotNil(t, rsp.Info)
	assert.NotEqual(t, "", rsp.Info.QT)
}

func TestGetTorrentGenericProperties(t *testing.T) {
	_, err := testApi.GetTorrentGenericProperties(context.Background(), &GetTorrentGenericPropertiesReq{Hash: "123"})
	assert.NoError(t, err)
	rsp, err := testApi.GetTorrentGenericProperties(context.Background(), &GetTorrentGenericPropertiesReq{Hash: testCfg.ValidHash})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp.Property)
}

func TestGetTorrentTrackers(t *testing.T) {
	rsp, err := testApi.GetTorrentTrackers(context.Background(), &GetTorrentTrackersReq{Hash: testCfg.ValidHash})
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(rsp.Trackers))
	t.Logf("data:%+v", rsp.Trackers)
}

func TestGetTorrentWebSeeds(t *testing.T) {
	rsp, err := testApi.GetTorrentWebSeeds(context.Background(), &GetTorrentWebSeedsReq{
		Hash: testCfg.ValidHash,
	})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp.WebSeeds)
}

func TestGetTorrentContents(t *testing.T) {
	rsp, err := testApi.GetTorrentContents(context.Background(), &GetTorrentContentsReq{Hash: testCfg.ValidHash})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp.Contents)
}

func TestGetTorrentPiecesStates(t *testing.T) {
	rsp, err := testApi.GetTorrentPiecesStates(context.Background(), &GetTorrentPiecesStatesReq{Hash: testCfg.ValidHash})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp.States)
}

func TestGetTorrentPiecesHashes(t *testing.T) {
	rsp, err := testApi.GetTorrentPiecesHashes(context.Background(), &GetTorrentPiecesHashesReq{Hash: testCfg.ValidHash})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp.Hashes)
}

func TestPauseTorrents(t *testing.T) {
	_, err := testApi.PauseTorrents(context.Background(), &PauseTorrentsReq{Hash: []string{testCfg.ValidHash}})
	assert.NoError(t, err)
}

func TestResumeTorrents(t *testing.T) {
	_, err := testApi.ResumeTorrents(context.Background(), &ResumeTorrentsReq{Hash: []string{testCfg.ValidHash}})
	assert.NoError(t, err)
}

func TestDeleteTorrents(t *testing.T) {
	_, err := testApi.DeleteTorrents(context.Background(), &DeleteTorrentsReq{Hash: []string{testCfg.ValidHash}})
	assert.NoError(t, err)
}

func TestRecheckTorrents(t *testing.T) {
	_, err := testApi.RecheckTorrents(context.Background(), &RecheckTorrentsReq{IsRecheckAll: true})
	assert.NoError(t, err)
}

func TestReannounceTorrents(t *testing.T) {
	_, err := testApi.ReannounceTorrents(context.Background(), &ReannounceTorrentsReq{IsReannounceAll: true})
	assert.NoError(t, err)
}

func TestAddTrackersToTorrent(t *testing.T) {
	_, err := testApi.AddTrackersToTorrent(context.Background(), &AddTrackersToTorrentReq{
		Hash: testCfg.ValidHash,
		Url: []string{
			"http://local.xxx.com/announce?a=x&b=2",
			"http://tracker.loadbt.com:6969/announce",
		},
	})
	assert.NoError(t, err)
}

func TestEditTrackers(t *testing.T) {
	_, err := testApi.EditTrackers(context.Background(), &EditTrackersReq{
		Hash:    testCfg.ValidHash,
		OrigUrl: "http://local.xxx.com/announce?a=x&b=2",
		NewUrl:  "http://local.ddd.com/announce?a=x&b=2",
	})
	assert.NoError(t, err)
}

func TestRemoveTrackers(t *testing.T) {
	_, err := testApi.RemoveTrackers(context.Background(), &RemoveTrackersReq{
		Hash: testCfg.ValidHash,
		Url: []string{
			"http://local.ddd.com/announce?a=x&b=2",
			"http://tracker.loadbt.com:6969/announce",
		},
	})
	assert.NoError(t, err)
}

func TestAddPeers(t *testing.T) {
	_, err := testApi.AddPeers(context.Background(), &AddPeersReq{
		Hash: []string{testCfg.ValidHash},
		Peer: []string{"192.168.50.220:8000"},
	})
	assert.NoError(t, err)
}

func TestIncreaseTorrentPriority(t *testing.T) {
	_, err := testApi.IncreaseTorrentPriority(context.Background(), &IncreaseTorrentPriorityReq{
		Hash:                 []string{testCfg.ValidHash},
		IsIncreaseAllTorrent: false,
	})
	assert.NoError(t, err)
}

func TestDecreaseTorrentPriority(t *testing.T) {
	_, err := testApi.DecreaseTorrentPriority(context.Background(), &DecreaseTorrentPriorityReq{
		Hash:                 []string{testCfg.ValidHash},
		IsDecreaseAllTorrent: false,
	})
	assert.NoError(t, err)
}

func TestMaximalTorrentPriority(t *testing.T) {
	_, err := testApi.MaximalTorrentPriority(context.Background(), &MaximalTorrentPriorityReq{
		Hash:                []string{testCfg.ValidHash},
		IsMaximalAllTorrent: false,
	})
	assert.NoError(t, err)
}

func TestMinimalTorrentPriority(t *testing.T) {
	_, err := testApi.MinimalTorrentPriority(context.Background(), &MinimalTorrentPriorityReq{
		Hash:                []string{testCfg.ValidHash},
		IsMinimalAllTorrent: false,
	})
	assert.NoError(t, err)
}

func TestSetFilePriority(t *testing.T) {
	_, err := testApi.SetFilePriority(context.Background(), &SetFilePriorityReq{
		Hash:     testCfg.ValidHash,
		Id:       []string{"4"},
		Priority: FilePriorityMaximal,
	})
	assert.NoError(t, err)
}

func TestGetTorrentDownloadLimit(t *testing.T) {
	rsp, err := testApi.GetTorrentDownloadLimit(context.Background(), &GetTorrentDownloadLimitReq{
		Hash: []string{testCfg.ValidHash},
	})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp.SpeedMap)
}

func TestSetTorrentDownloadLimit(t *testing.T) {
	_, err := testApi.SetTorrentDownloadLimit(context.Background(), &SetTorrentDownloadLimitReq{
		Hash:  []string{testCfg.ValidHash},
		Speed: 2 * 1024 * 1024,
	})
	assert.NoError(t, err)
}

func TestSetTorrentShareLimit(t *testing.T) {
	_, err := testApi.SetTorrentShareLimit(context.Background(), &SetTorrentShareLimitReq{
		Hash:             []string{testCfg.ValidHash},
		SeedingTimeLimit: 86400 * 7,
		RatioLimit:       10,
	})
	assert.NoError(t, err)
}

func TestGetTorrentUploadLimit(t *testing.T) {
	rsp, err := testApi.GetTorrentUploadLimit(context.Background(), &GetTorrentUploadLimitReq{
		Hash: []string{testCfg.ValidHash},
	})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp.SpeedMap)
}

func TestSetTorrentUploadLimit(t *testing.T) {
	_, err := testApi.SetTorrentUploadLimit(context.Background(), &SetTorrentUploadLimitReq{
		Hash:  []string{testCfg.ValidHash},
		Speed: 300 * 1024,
	})
	assert.NoError(t, err)
}

func TestSetTorrentLocation(t *testing.T) {
	_, err := testApi.SetTorrentLocation(context.Background(), &SetTorrentLocationReq{
		Hash:     []string{testCfg.ValidHash},
		Location: "abc",
	})
	assert.NoError(t, err)
}

func TestSetTorrentName(t *testing.T) {
	_, err := testApi.SetTorrentName(context.Background(), &SetTorrentNameReq{
		Hash: testCfg.ValidHash,
		Name: "abc",
	})
	assert.NoError(t, err)
}

func TestSetTorrentCategory(t *testing.T) {
	_, err := testApi.SetTorrentCategory(context.Background(), &SetTorrentCategoryReq{
		Hash:     []string{testCfg.ValidHash},
		Category: "电影",
	})
	assert.NoError(t, err)
}

func TestGetAllCategories(t *testing.T) {
	rsp, err := testApi.GetAllCategories(context.Background(), &GetAllCategoriesReq{})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp.Categories)
}

func TestAddNewCategory(t *testing.T) {
	_, err := testApi.AddNewCategory(context.Background(), &AddNewCategoryReq{
		Category: "abc",
		SavePath: "",
	})
	assert.NoError(t, err)
}

func TestEditCategory(t *testing.T) {
	_, err := testApi.EditCategory(context.Background(), &EditCategoryReq{
		Category: "abc",
		SavePath: "/123",
	})
	assert.NoError(t, err)
}

func TestRemoveCategories(t *testing.T) {
	_, err := testApi.RemoveCategories(context.Background(), &RemoveCategoriesReq{
		Category: []string{"abc"},
	})
	assert.NoError(t, err)
}

func TestAddTorrentTags(t *testing.T) {
	_, err := testApi.AddTorrentTags(context.Background(), &AddTorrentTagsReq{
		Hash: []string{testCfg.ValidHash},
		Tag:  []string{"1", "2", "a", "b"},
	})
	assert.NoError(t, err)
}

func TestRemoveTorrentTags(t *testing.T) {
	_, err := testApi.RemoveTorrentTags(context.Background(), &RemoveTorrentTagsReq{
		Hash: []string{testCfg.ValidHash},
		Tag:  []string{"1", "a"},
	})
	assert.NoError(t, err)
}

func TestGetAllTags(t *testing.T) {
	rsp, err := testApi.GetAllTags(context.Background(), &GetAllTagsReq{})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp.Tags)
}

func TestCreateTags(t *testing.T) {
	_, err := testApi.CreateTags(context.Background(), &CreateTagsReq{
		Tag: []string{"aa", "bvb"},
	})
	assert.NoError(t, err)
}

func TestDeleteTags(t *testing.T) {
	_, err := testApi.DeleteTags(context.Background(), &DeleteTagsReq{
		Tag: []string{"aa", "bvb"},
	})
	assert.NoError(t, err)
}

func TestSetAutomaticTorrentManagement(t *testing.T) {
	_, err := testApi.SetAutomaticTorrentManagement(context.Background(), &SetAutomaticTorrentManagementReq{
		Hash: []string{testCfg.ValidHash},
	})
	assert.NoError(t, err)
}

func TestToggleSequentialDownload(t *testing.T) {
	_, err := testApi.ToggleSequentialDownload(context.Background(), &ToggleSequentialDownloadReq{
		Hash: []string{testCfg.ValidHash},
	})
	assert.NoError(t, err)
}

func TestSetFirstOrLastPiecePriority(t *testing.T) {
	_, err := testApi.SetFirstOrLastPiecePriority(context.Background(), &SetFirstOrLastPiecePriorityReq{
		Hash: []string{testCfg.ValidHash},
	})
	assert.NoError(t, err)
}

func TestSetForceStart(t *testing.T) {
	_, err := testApi.SetForceStart(context.Background(), &SetForceStartReq{
		Hash:  []string{testCfg.ValidHash},
		Value: true,
	})
	assert.NoError(t, err)
}

func TestSetSuperSeeding(t *testing.T) {
	_, err := testApi.SetSuperSeeding(context.Background(), &SetSuperSeedingReq{
		Hash:  []string{testCfg.ValidHash},
		Value: true,
	})
	assert.NoError(t, err)
}

func TestRenameFile(t *testing.T) {
	_, err := testApi.RenameFile(context.Background(), &RenameFileReq{
		Hash:    testCfg.ValidHash,
		OldPath: "abc",
		NewPath: "abcd",
	})
	assert.NoError(t, err)
}

func TestRenameFolder(t *testing.T) {
	_, err := testApi.RenameFolder(context.Background(), &RenameFolderReq{
		Hash:    testCfg.ValidHash,
		OldPath: "/downloads/",
		NewPath: "downloads2",
	})
	assert.NoError(t, err)
}

func TestSetApplicationPreferences(t *testing.T) {
	var dl int = 100 * 1024 * 1024
	_, err := testApi.SetApplicationPreferences(context.Background(), &SetApplicationPreferencesReq{
		DlLimit: &dl,
	})
	assert.NoError(t, err)
}
