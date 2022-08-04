package qbapi

const (
	ErrOK         = 0
	ErrParams     = -10000
	ErrUnmarsal   = -10001
	ErrMarsal     = -10002
	ErrInternal   = -10003
	ErrNetwork    = -10004
	ErrStatusCode = -10005
	ErrLogin      = -10006
	ErrUnknown    = -10007
	ErrFile       = -10008
)

type FilePriority int

const (
	FilePriorityDoNotDownload = 0
	FilePriorityNormal        = 1
	FilePriorityHigh          = 6
	FilePriorityMaximal       = 7
)

const (
	//login
	apiLogin = "/api/v2/auth/login"
	//
	apiGetAPPVersion      = "/api/v2/app/version"
	apiGetAPIVersion      = "/api/v2/app/webapiVersion"
	apiGetBuildInfo       = "/api/v2/app/buildInfo"
	apiShutdownAPP        = "/api/v2/app/shutdown"
	apiGetAPPPerf         = "/api/v2/app/preferences"
	apiSetAPPPref         = "/api/v2/app/setPreferences"
	apiGetDefaultSavePath = "/api/v2/app/defaultSavePath"
	//Log
	apiGetLog     = "/api/v2/log/main"
	apiGetPeerLog = "/api/v2/log/peers"
	//sync
	apiGetMainData        = "/api/v2/sync/maindata"
	apiGetTorrentPeerData = "/api/v2/sync/torrentPeers"
	//transfer info
	apiGetGlobalTransferInfo  = "/api/v2/transfer/info"
	apiGetAltSpeedLimitState  = "/api/v2/transfer/speedLimitsMode"
	apiToggleAltSpeedLimits   = "/api/v2/transfer/toggleSpeedLimitsMode"
	apiGetGlobalDownloadLimit = "/api/v2/transfer/downloadLimit"
	apiSetGlobalDownloadLimit = "/api/v2/transfer/setDownloadLimit"
	apiGetGlobalUploadLimit   = "/api/v2/transfer/uploadLimit"
	apiSetGlobalUploadLimit   = "/api/v2/transfer/setUploadLimit"
	apiBanPeers               = "/api/v2/transfer/banPeers"
	//torrent management
	apiGetTorrentList                = "/api/v2/torrents/info"
	apiGetTorrentGenericProp         = "/api/v2/torrents/properties"
	apiGetTorrentTrackers            = "/api/v2/torrents/trackers"
	apiGetTorrentWebSeeds            = "/api/v2/torrents/webseeds"
	apiGetTorrentContents            = "/api/v2/torrents/files"
	apiGetTorrentPiecesStates        = "/api/v2/torrents/pieceStates"
	apiGetTorrentPiecesHashes        = "/api/v2/torrents/pieceHashes"
	apiPauseTorrents                 = "/api/v2/torrents/pause"
	apiResumeTorrents                = "/api/v2/torrents/resume"
	apiDeleteTorrents                = "/api/v2/torrents/delete"
	apiRecheckTorrents               = "/api/v2/torrents/recheck"
	apiReannounceTorrents            = "/api/v2/torrents/reannounce"
	apiAddNewTorrent                 = "/api/v2/torrents/add" //form upload
	apiAddTrackersToTorrent          = "/api/v2/torrents/addTrackers"
	apiEditTrackers                  = "/api/v2/torrents/editTracker"
	apiRemoveTrackers                = "/api/v2/torrents/removeTrackers"
	apiAddPeers                      = "/api/v2/torrents/addPeers"
	apiIncreaseTorrentPriority       = "/api/v2/torrents/increasePrio"
	apiDecreaseTorrentPriority       = "/api/v2/torrents/decreasePrio"
	apiMaximalTorrentPriority        = "/api/v2/torrents/topPrio"
	apiMinimalTorrentPriority        = "/api/v2/torrents/bottomPrio"
	apiSetFilePriority               = "/api/v2/torrents/filePrio"
	apiGetTorrentDownloadLimit       = "/api/v2/torrents/downloadLimit"
	apiSetTorrentDownloadLimit       = "/api/v2/torrents/setDownloadLimit"
	apiSetTorrentShareLimit          = "/api/v2/torrents/setShareLimits"
	apiGetTorrentUploadLimit         = "/api/v2/torrents/uploadLimit"
	apiSetTorrentUploadLimit         = "/api/v2/torrents/setUploadLimit"
	apiSetTorrentLocation            = "/api/v2/torrents/setLocation"
	apiSetTorrentName                = "/api/v2/torrents/rename"
	apiSetTorrentCategory            = "/api/v2/torrents/setCategory"
	apiGetAllCategories              = "/api/v2/torrents/categories"
	apiAddNewCategory                = "/api/v2/torrents/createCategory"
	apiEditCategory                  = "/api/v2/torrents/editCategory"
	apiRemoveCategories              = "/api/v2/torrents/removeCategories"
	apiAddTorrentTags                = "/api/v2/torrents/addTags"
	apiRemoveTorrentTags             = "/api/v2/torrents/removeTags"
	apiGetAllTags                    = "/api/v2/torrents/tags"
	apiCreateTags                    = "/api/v2/torrents/createTags"
	apiDeleteTags                    = "/api/v2/torrents/deleteTags"
	apiSetAutomaticTorrentManagement = "/api/v2/torrents/setAutoManagement"
	apiToggleSequentialDownload      = "/api/v2/torrents/toggleSequentialDownload"
	apiSetFirstOrLastPiecePriority   = "/api/v2/torrents/toggleFirstLastPiecePrio"
	apiSetForceStart                 = "/api/v2/torrents/setForceStart"
	apiSetSuperSeeding               = "/api/v2/torrents/setSuperSeeding"
	apiRenameFile                    = "/api/v2/torrents/renameFile"
	apiRenameFolder                  = "/api/v2/torrents/renameFolder"
	//RSS (experimental)

)
