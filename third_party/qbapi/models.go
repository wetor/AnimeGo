package qbapi

//doc: https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-4.1)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRsp struct {
}

type GetTorrentListReq struct {
	Filter   *string `json:"filter,omitempty"`
	Category *string `json:"category,omitempty"`
	Tag      *string `json:"tag,omitempty"`
	Sort     *string `json:"sort,omitempty"`
	Reverse  *bool   `json:"reverse,omitempty"`
	Limit    *int    `json:"limit,omitempty"`
	Offset   *int    `json:"offset,omitempty"`
	Hashes   *string `json:"hashes,omitempty"`
}

type TorrentListItem struct {
	AddedOn           int     `json:"added_on"`
	AmountLeft        int     `json:"amount_left"`
	AutoTmm           bool    `json:"auto_tmm"`
	Availability      float64 `json:"availability"`
	Category          string  `json:"category"`
	Completed         int     `json:"completed"`
	CompletionOn      int     `json:"completion_on"`
	ContentPath       string  `json:"content_path"`
	DlLimit           int     `json:"dl_limit"`
	Dlspeed           int     `json:"dlspeed"`
	Downloaded        int     `json:"downloaded"`
	DownloadedSession int     `json:"downloaded_session"`
	Eta               int     `json:"eta"`
	FLPiecePrio       bool    `json:"f_l_piece_prio"`
	ForceStart        bool    `json:"force_start"`
	Hash              string  `json:"hash"`
	LastActivity      int     `json:"last_activity"`
	MagnetUri         string  `json:"magnet_uri"`
	MaxRatio          float64 `json:"max_ratio"`
	MaxSeedingTime    int     `json:"max_seeding_time"`
	Name              string  `json:"name"`
	NumComplete       int     `json:"num_complete"`
	NumIncomplete     int     `json:"num_incomplete"`
	NumLeechs         int     `json:"num_leechs"`
	NumSeeds          int     `json:"num_seeds"`
	Priority          int     `json:"priority"`
	Progress          float64 `json:"progress"`
	Ratio             float64 `json:"ratio"`
	RatioLimit        float64 `json:"ratio_limit"`
	SavePath          string  `json:"save_path"`
	SeedingTime       int     `json:"seeding_time"`
	SeedingTimeLimit  int     `json:"seeding_time_limit"`
	SeenComplete      int     `json:"seen_complete"`
	SeqDl             bool    `json:"seq_dl"`
	Size              int     `json:"size"`
	State             string  `json:"state"`
	SuperSeeding      bool    `json:"super_seeding"`
	Tags              string  `json:"tags"`
	TimeActive        int     `json:"time_active"`
	TotalSize         int     `json:"total_size"`
	Tracker           string  `json:"tracker"`
	UpLimit           int     `json:"up_limit"`
	Uploaded          int     `json:"uploaded"`
	UploadedSession   int     `json:"uploaded_session"`
	Upspeed           int     `json:"upspeed"`
}

type GetTorrentListRsp struct {
	Items []*TorrentListItem
}

type GetApplicationVersionReq struct {
}

type GetApplicationVersionRsp struct {
	Version string
}

type GetAPIVersionReq struct {
}

type GetAPIVersionRsp struct {
	Version string
}

type GetBuildInfoReq struct{}

type BuildInfo struct {
	QT         string `json:"qt"`
	Libtorrent string `json:"libtorrent"`
	Boost      string `json:"boost"`
	Openssl    string `json:"openssl"`
	Bitness    int    `json:"bitness"`
}

type GetBuildInfoRsp struct {
	Info *BuildInfo
}

type ShutdownApplicationReq struct {
}

type ShutdownApplicationRsp struct {
}

type GetApplicationPreferencesReq struct {
}

type GetApplicationPreferencesRsp struct {
	Locale                             string                 `json:"locale"`
	CreateSubfolderEnabled             bool                   `json:"create_subfolder_enabled"`
	StartPausedEnabled                 bool                   `json:"start_paused_enabled"`
	AutoDeleteMode                     int                    `json:"auto_delete_mode"`
	PreallocateAll                     bool                   `json:"preallocate_all"`
	IncompleteFilesExt                 bool                   `json:"incomplete_files_ext"`
	AutoTmmEnabled                     bool                   `json:"auto_tmm_enabled"`                      //True if Automatic Torrent Management is enabled by default
	TorrentChangedTmmEnabled           bool                   `json:"torrent_changed_tmm_enabled"`           //True if torrent should be relocated when its Category changes
	SavePathChangedTmmEnabled          bool                   `json:"save_path_changed_tmm_enabled"`         //True if torrent should be relocated when the default save path changes
	CategoryChangedTmmEnabled          bool                   `json:"category_changed_tmm_enabled"`          //True if torrent should be relocated when its Category's save path changes
	SavePath                           string                 `json:"save_path"`                             //Default save path for torrents, separated by slashes
	TempPathEnabled                    bool                   `json:"temp_path_enabled"`                     //True if folder for incomplete torrents is enabled
	TempPath                           string                 `json:"temp_path"`                             //Path for incomplete torrents, separated by slashes
	ScanDirs                           map[string]interface{} `json:"scan_dirs"`                             //Property: directory to watch for torrent files, value: where torrents loaded from this directory should be downloaded to (see list of possible values below). Slashes are used as path separators; multiple key/value pairs can be specified
	ExportDir                          string                 `json:"export_dir"`                            //Path to directory to copy .torrent files to. Slashes are used as path separators
	ExportDirFin                       string                 `json:"export_dir_fin"`                        //Path to directory to copy .torrent files of completed downloads to. Slashes are used as path separators
	MailNotificationEnabled            bool                   `json:"mail_notification_enabled"`             //True if e-mail notification should be enabled
	MailNotificationSender             string                 `json:"mail_notification_sender"`              //e-mail where notifications should originate from
	MailNotificationEmail              string                 `json:"mail_notification_email"`               //e-mail to send notifications to
	MailNotificationSmtp               string                 `json:"mail_notification_smtp"`                //smtp server for e-mail notifications
	MailNotificationSslEnabled         bool                   `json:"mail_notification_ssl_enabled"`         //True if smtp server requires SSL connection
	MailNotificationAuthEnabled        bool                   `json:"mail_notification_auth_enabled"`        //True if smtp server requires authentication
	MailNotificationUsername           string                 `json:"mail_notification_username"`            //Username for smtp authentication
	MailNotificationPassword           string                 `json:"mail_notification_password"`            //Password for smtp authentication
	AutorunEnabled                     bool                   `json:"autorun_enabled"`                       //True if external program should be run after torrent has finished downloading
	AutorunProgram                     string                 `json:"autorun_program"`                       //Program path/name/arguments to run if autorun_enabled is enabled; path is separated by slashes; you can use %f and %n arguments, which will be expanded by qBittorent as path_to_torrent_file and torrent_name (from the GUI; not the .torrent file name) respectively
	QueueingEnabled                    bool                   `json:"queueing_enabled"`                      //True if torrent queuing is enabled
	MaxActiveDownloads                 int                    `json:"max_active_downloads"`                  //Maximum number of active simultaneous downloads
	MaxActiveTorrents                  int                    `json:"max_active_torrents"`                   //Maximum number of active simultaneous downloads and uploads
	MaxActiveUploads                   int                    `json:"max_active_uploads"`                    //Maximum number of active simultaneous uploads
	DontCountSlowTorrents              bool                   `json:"dont_count_slow_torrents"`              //If true torrents w/o any activity (stalled ones) will not be counted towards max_active_* limits; see dont_count_slow_torrents for more information
	SlowTorrentDlRateThreshold         int                    `json:"slow_torrent_dl_rate_threshold"`        //Download rate in KiB/s for a torrent to be considered "slow"
	SlowTorrentUlRateThreshold         int                    `json:"slow_torrent_ul_rate_threshold"`        //Upload rate in KiB/s for a torrent to be considered "slow"
	SlowTorrentInactiveTimer           int                    `json:"slow_torrent_inactive_timer"`           //Seconds a torrent should be inactive before considered "slow"
	MaxRatioEnabled                    bool                   `json:"max_ratio_enabled"`                     //True if share ratio limit is enabled
	MaxRatio                           float64                `json:"max_ratio"`                             //Get the global share ratio limit
	MaxRatioAct                        int                    `json:"max_ratio_act"`                         //Action performed when a torrent reaches the maximum share ratio. See list of possible values here below.
	ListenPort                         int                    `json:"listen_port"`                           //Port for incoming connections
	Upnp                               bool                   `json:"upnp"`                                  //True if UPnP/NAT-PMP is enabled
	RandomPort                         bool                   `json:"random_port"`                           //True if the port is randomly selected
	DlLimit                            int                    `json:"dl_limit"`                              //Global download speed limit in KiB/s; -1 means no limit is applied
	UpLimit                            int                    `json:"up_limit"`                              //Global upload speed limit in KiB/s; -1 means no limit is applied
	MaxConnec                          int                    `json:"max_connec"`                            //Maximum global number of simultaneous connections
	MaxConnecPerTorrent                int                    `json:"max_connec_per_torrent"`                //Maximum number of simultaneous connections per torrent
	MaxUploads                         int                    `json:"max_uploads"`                           //Maximum number of upload slots
	MaxUploadsPerTorrent               int                    `json:"max_uploads_per_torrent"`               //Maximum number of upload slots per torrent
	StopTrackerTimeout                 int                    `json:"stop_tracker_timeout"`                  //Timeout in seconds for a stopped announce request to trackers
	EnablePieceExtentAffinity          bool                   `json:"enable_piece_extent_affinity"`          //True if the advanced libtorrent option piece_extent_affinity is enabled
	BittorrentProtocol                 int                    `json:"bittorrent_protocol"`                   //Bittorrent Protocol to use (see list of possible values below)
	LimitUtpRate                       bool                   `json:"limit_utp_rate"`                        //True if [du]l_limit should be applied to uTP connections; this option is only available in qBittorent built against libtorrent version 0.16.X and higher
	LimitTcpOverhead                   bool                   `json:"limit_tcp_overhead"`                    //True if [du]l_limit should be applied to estimated TCP overhead (service data: e.g. packet headers)
	LimitLanPeers                      bool                   `json:"limit_lan_peers"`                       //True if [du]l_limit should be applied to peers on the LAN
	AltDlLimit                         int                    `json:"alt_dl_limit"`                          //Alternative global download speed limit in KiB/s
	AltUpLimit                         int                    `json:"alt_up_limit"`                          //Alternative global upload speed limit in KiB/s
	SchedulerEnabled                   bool                   `json:"scheduler_enabled"`                     //True if alternative limits should be applied according to schedule
	ScheduleFromHour                   int                    `json:"schedule_from_hour"`                    //Scheduler starting hour
	ScheduleFromMin                    int                    `json:"schedule_from_min"`                     //Scheduler starting minute
	ScheduleToHour                     int                    `json:"schedule_to_hour"`                      //Scheduler ending hour
	ScheduleToMin                      int                    `json:"schedule_to_min"`                       //Scheduler ending minute
	SchedulerDays                      int                    `json:"scheduler_days"`                        //Scheduler days. See possible values here below
	Dht                                bool                   `json:"dht"`                                   //True if DHT is enabled
	Pex                                bool                   `json:"pex"`                                   //True if PeX is enabled
	Lsd                                bool                   `json:"lsd"`                                   //True if LSD is enabled
	Encryption                         int                    `json:"encryption"`                            //See list of possible values here below
	AnonymousMode                      bool                   `json:"anonymous_mode"`                        //If true anonymous mode will be enabled; read more here; this option is only available in qBittorent built against libtorrent version 0.16.X and higher
	ProxyType                          int                    `json:"proxy_type"`                            //See list of possible values here below
	ProxyIp                            string                 `json:"proxy_ip"`                              //Proxy IP address or domain name
	ProxyPort                          int                    `json:"proxy_port"`                            //Proxy port
	ProxyPeerConnections               bool                   `json:"proxy_peer_connections"`                //True if peer and web seed connections should be proxified; this option will have any effect only in qBittorent built against libtorrent version 0.16.X and higher
	ProxyAuthEnabled                   bool                   `json:"proxy_auth_enabled"`                    //True proxy requires authentication; doesn't apply to SOCKS4 proxies
	ProxyUsername                      string                 `json:"proxy_username"`                        //Username for proxy authentication
	ProxyPassword                      string                 `json:"proxy_password"`                        //Password for proxy authentication
	ProxyTorrentsOnly                  bool                   `json:"proxy_torrents_only"`                   //True if proxy is only used for torrents
	IpFilterEnabled                    bool                   `json:"ip_filter_enabled"`                     //True if external IP filter should be enabled
	IpFilterPath                       string                 `json:"ip_filter_path"`                        //Path to IP filter file (.dat, .p2p, .p2b files are supported); path is separated by slashes
	IpFilterTrackers                   bool                   `json:"ip_filter_trackers"`                    //True if IP filters are applied to trackers
	WebUiDomainList                    string                 `json:"web_ui_domain_list"`                    //Comma-separated list of domains to accept when performing Host header validation
	WebUiAddress                       string                 `json:"web_ui_address"`                        //IP address to use for the WebUI
	WebUiPort                          int                    `json:"web_ui_port"`                           //WebUI port
	WebUiUpnp                          bool                   `json:"web_ui_upnp"`                           //True if UPnP is used for the WebUI port
	WebUiUsername                      string                 `json:"web_ui_username"`                       //WebUI username
	WebUiPassword                      string                 `json:"web_ui_password"`                       //For API ≥ v2.3.0: Plaintext WebUI password, not readable, write-only. For API < v2.3.0: MD5 hash of WebUI password, hash is generated from the following string: username:Web UI Access:plain_text_web_ui_password
	WebUiCsrfProtectionEnabled         bool                   `json:"web_ui_csrf_protection_enabled"`        //True if WebUI CSRF protection is enabled
	WebUiClickjackingProtectionEnabled bool                   `json:"web_ui_clickjacking_protection_enable"` //True if WebUI clickjacking protection is enabled
	WebUiSecureCookieEnabled           bool                   `json:"web_ui_secure_cookie_enabled"`          //True if WebUI cookie Secure flag is enabled
	WebUiMaxAuthFailCount              int                    `json:"web_ui_max_auth_fail_count"`            //Maximum number of authentication failures before WebUI access ban
	WebUiBanDuration                   int                    `json:"web_ui_ban_duration"`                   //WebUI access ban duration in seconds
	WebUiSessionTimeout                int                    `json:"web_ui_session_timeout"`                //Seconds until WebUI is automatically signed off
	WebUiHostHeaderValidationEnabled   bool                   `json:"web_ui_host_header_validation_enabled"` //True if WebUI host header validation is enabled
	BypassLocalAuth                    bool                   `json:"bypass_local_auth"`                     //True if authentication challenge for loopback address (127.0.0.1) should be disabled
	BypassAuthSubnetWhitelistEnabled   bool                   `json:"bypass_auth_subnet_whitelist_enabled"`  //True if webui authentication should be bypassed for clients whose ip resides within (at least) one of the subnets on the whitelist
	BypassAuthSubnetWhitelist          string                 `json:"bypass_auth_subnet_whitelist"`          //(White)list of ipv4/ipv6 subnets for which webui authentication should be bypassed; list entries are separated by commas
	AlternativeWebuiEnabled            bool                   `json:"alternative_webui_enabled"`             //True if an alternative WebUI should be used
	AlternativeWebuiPath               string                 `json:"alternative_webui_path"`                //File path to the alternative WebUI
	UseHttps                           bool                   `json:"use_https"`                             //True if WebUI HTTPS access is enabled
	SslKey                             string                 `json:"ssl_key"`                               //For API < v2.0.1: SSL keyfile contents (this is a not a path)
	SslCert                            string                 `json:"ssl_cert"`                              //For API < v2.0.1: SSL certificate contents (this is a not a path)
	WebUiHttpsKeyPath                  string                 `json:"web_ui_https_key_path"`                 //For API ≥ v2.0.1: Path to SSL keyfile
	WebUiHttpsCertPath                 string                 `json:"web_ui_https_cert_path"`                //For API ≥ v2.0.1: Path to SSL certificate
	DyndnsEnabled                      bool                   `json:"dyndns_enabled"`                        //True if server DNS should be updated dynamically
	DyndnsService                      int                    `json:"dyndns_service"`                        //See list of possible values here below
	DyndnsUsername                     string                 `json:"dyndns_username"`                       //Username for DDNS service
	DyndnsPassword                     string                 `json:"dyndns_password"`                       //Password for DDNS service
	DyndnsDomain                       string                 `json:"dyndns_domain"`                         //Your DDNS domain name
	RssRefreshInterval                 int                    `json:"rss_refresh_interval"`                  //RSS refresh interval
	RssMaxArticlesPerFeed              int                    `json:"rss_max_articles_per_feed"`             //Max stored articles per RSS feed
	RssProcessingEnabled               bool                   `json:"rss_processing_enabled"`                //Enable processing of RSS feeds
	RssAutoDownloadingEnabled          bool                   `json:"rss_auto_downloading_enabled"`          //Enable auto-downloading of torrents from the RSS feeds
	RssDownloadRepackProperEpisodes    bool                   `json:"rss_download_repack_proper_episodes"`   //For API ≥ v2.5.1: Enable downloading of repack/proper Episodes
	RssSmartEpisodeFilters             string                 `json:"rss_smart_episode_filters"`             //For API ≥ v2.5.1: List of RSS Smart Episode Filters
	AddTrackersEnabled                 bool                   `json:"add_trackers_enabled"`                  //Enable automatic adding of trackers to new torrents
	AddTrackers                        string                 `json:"add_trackers"`                          //List of trackers to add to new torrent
	WebUiUseCustomHttpHeadersEnabled   bool                   `json:"web_ui_use_custom_http_headers_enable"` //For API ≥ v2.5.1: Enable custom http headers
	WebUiCustomHttpHeaders             string                 `json:"web_ui_custom_http_headers"`            //For API ≥ v2.5.1: List of custom http headers
	MaxSeedingTimeEnabled              bool                   `json:"max_seeding_time_enabled"`              //True enables max seeding time
	MaxSeedingTime                     int                    `json:"max_seeding_time"`                      //Number of minutes to seed a torrent
	AnnounceIp                         string                 `json:"announce_ip"`                           //TODO
	AnnounceToAllTiers                 bool                   `json:"announce_to_all_tiers"`                 //True always announce to all tiers
	AnnounceToAllTrackers              bool                   `json:"announce_to_all_trackers"`              //True always announce to all trackers in a tier
	AsyncIoThreads                     int                    `json:"async_io_threads"`                      //Number of asynchronous I/O threads
	BannedIps                          string                 `json:"banned_IPs"`                            //List of banned IPs
	CheckingMemoryUse                  int                    `json:"checking_memory_use"`                   //Outstanding memory when checking torrents in MiB
	CurrentInterfaceAddress            string                 `json:"current_interface_address"`             //IP Address to bind to. Empty String means All addresses
	CurrentNetworkInterface            string                 `json:"current_network_interface"`             //Network Interface used
	DiskCache                          int                    `json:"disk_cache"`                            //Disk cache used in MiB
	DiskCacheTtl                       int                    `json:"disk_cache_ttl"`                        //Disk cache expiry interval in seconds
	EmbeddedTrackerPort                int                    `json:"embedded_tracker_port"`                 //Port used for embedded tracker
	EnableCoalesceReadWrite            bool                   `json:"enable_coalesce_read_write"`            //True enables coalesce reads & writes
	EnableEmbeddedTracker              bool                   `json:"enable_embedded_tracker"`               //True enables embedded tracker
	EnableMultiConnectionsFromSameIp   bool                   `json:"enable_multi_connections_from_same_ip"` //True allows multiple connections from the same IP address
	EnableOsCache                      bool                   `json:"enable_os_cache"`                       //True enables os cache
	EnableUploadSuggestions            bool                   `json:"enable_upload_suggestions"`             //True enables sending of upload piece suggestions
	FilePoolSize                       int                    `json:"file_pool_size"`                        //File pool size
	OutgoingPortsMax                   int                    `json:"outgoing_ports_max"`                    //Maximal outgoing port (0: Disabled)
	OutgoingPortsMin                   int                    `json:"outgoing_ports_min"`                    //Minimal outgoing port (0: Disabled)
	RecheckCompletedTorrents           bool                   `json:"recheck_completed_torrents"`            //True rechecks torrents on completion
	ResolvePeerCountries               bool                   `json:"resolve_peer_countries"`                //True resolves peer countries
	SaveResumeDataInterval             int                    `json:"save_resume_data_interval"`             //Save resume data interval in min
	SendBufferLowWatermark             int                    `json:"send_buffer_low_watermark"`             //Send buffer low watermark in KiB
	SendBufferWatermark                int                    `json:"send_buffer_watermark"`                 //Send buffer watermark in KiB
	SendBufferWatermarkFactor          int                    `json:"send_buffer_watermark_factor"`          //Send buffer watermark factor in percent
	SocketBacklogSize                  int                    `json:"socket_backlog_size"`                   //Socket backlog size
	UploadChokingAlgorithm             int                    `json:"upload_choking_algorithm"`              //Upload choking algorithm used (see list of possible values below)
	UploadSlotsBehavior                int                    `json:"upload_slots_behavior"`                 //Upload slots behavior used (see list of possible values below)
	UpnpLeaseDuration                  int                    `json:"upnp_lease_duration"`                   //UPnP lease duration (0: Permanent lease)
	UtpTcpMixedMode                    int                    `json:"utp_tcp_mixed_mode"`
}

type setApplicationPreferencesInnerReq struct {
	Json string `json:"json"`
}

type SetApplicationPreferencesReq struct {
	TorrentContentLayout               *string                `json:"torrent_content_layout,omitempty"`
	Locale                             *string                `json:"locale,omitempty"`
	CreateSubfolderEnabled             *bool                  `json:"create_subfolder_enabled,omitempty"`
	StartPausedEnabled                 *bool                  `json:"start_paused_enabled,omitempty"`
	AutoDeleteMode                     *int                   `json:"auto_delete_mode,omitempty"`
	PreallocateAll                     *bool                  `json:"preallocate_all,omitempty"`
	IncompleteFilesExt                 *bool                  `json:"incomplete_files_ext,omitempty"`
	AutoTmmEnabled                     *bool                  `json:"auto_tmm_enabled,omitempty"`                      //True if Automatic Torrent Management is enabled by default
	TorrentChangedTmmEnabled           *bool                  `json:"torrent_changed_tmm_enabled,omitempty"`           //True if torrent should be relocated when its Category changes
	SavePathChangedTmmEnabled          *bool                  `json:"save_path_changed_tmm_enabled,omitempty"`         //True if torrent should be relocated when the default save path changes
	CategoryChangedTmmEnabled          *bool                  `json:"category_changed_tmm_enabled,omitempty"`          //True if torrent should be relocated when its Category's save path changes
	SavePath                           *string                `json:"save_path,omitempty"`                             //Default save path for torrents, separated by slashes
	TempPathEnabled                    *bool                  `json:"temp_path_enabled,omitempty"`                     //True if folder for incomplete torrents is enabled
	TempPath                           *string                `json:"temp_path,omitempty"`                             //Path for incomplete torrents, separated by slashes
	ScanDirs                           map[string]interface{} `json:"scan_dirs,omitempty"`                             //Property: directory to watch for torrent files, value: where torrents loaded from this directory should be downloaded to (see list of possible values below). Slashes are used as path separators; multiple key/value pairs can be specified
	ExportDir                          *string                `json:"export_dir,omitempty"`                            //Path to directory to copy .torrent files to. Slashes are used as path separators
	ExportDirFin                       *string                `json:"export_dir_fin,omitempty"`                        //Path to directory to copy .torrent files of completed downloads to. Slashes are used as path separators
	MailNotificationEnabled            *bool                  `json:"mail_notification_enabled,omitempty"`             //True if e-mail notification should be enabled
	MailNotificationSender             *string                `json:"mail_notification_sender,omitempty"`              //e-mail where notifications should originate from
	MailNotificationEmail              *string                `json:"mail_notification_email,omitempty"`               //e-mail to send notifications to
	MailNotificationSmtp               *string                `json:"mail_notification_smtp,omitempty"`                //smtp server for e-mail notifications
	MailNotificationSslEnabled         *bool                  `json:"mail_notification_ssl_enabled,omitempty"`         //True if smtp server requires SSL connection
	MailNotificationAuthEnabled        *bool                  `json:"mail_notification_auth_enabled,omitempty"`        //True if smtp server requires authentication
	MailNotificationUsername           *string                `json:"mail_notification_username,omitempty"`            //Username for smtp authentication
	MailNotificationPassword           *string                `json:"mail_notification_password,omitempty"`            //Password for smtp authentication
	AutorunEnabled                     *bool                  `json:"autorun_enabled,omitempty"`                       //True if external program should be run after torrent has finished downloading
	AutorunProgram                     *string                `json:"autorun_program,omitempty"`                       //Program path/name/arguments to run if autorun_enabled is enabled; path is separated by slashes; you can use %f and %n arguments, which will be expanded by qBittorent as path_to_torrent_file and torrent_name (from the GUI; not the .torrent file name) respectively
	QueueingEnabled                    *bool                  `json:"queueing_enabled,omitempty"`                      //True if torrent queuing is enabled
	MaxActiveDownloads                 *int                   `json:"max_active_downloads,omitempty"`                  //Maximum number of active simultaneous downloads
	MaxActiveTorrents                  *int                   `json:"max_active_torrents,omitempty"`                   //Maximum number of active simultaneous downloads and uploads
	MaxActiveUploads                   *int                   `json:"max_active_uploads,omitempty"`                    //Maximum number of active simultaneous uploads
	DontCountSlowTorrents              *bool                  `json:"dont_count_slow_torrents,omitempty"`              //If true torrents w/o any activity (stalled ones) will not be counted towards max_active_* limits; see dont_count_slow_torrents for more information
	SlowTorrentDlRateThreshold         *int                   `json:"slow_torrent_dl_rate_threshold,omitempty"`        //Download rate in KiB/s for a torrent to be considered "slow"
	SlowTorrentUlRateThreshold         *int                   `json:"slow_torrent_ul_rate_threshold,omitempty"`        //Upload rate in KiB/s for a torrent to be considered "slow"
	SlowTorrentInactiveTimer           *int                   `json:"slow_torrent_inactive_timer,omitempty"`           //Seconds a torrent should be inactive before considered "slow"
	MaxRatioEnabled                    *bool                  `json:"max_ratio_enabled,omitempty"`                     //True if share ratio limit is enabled
	MaxRatio                           *float64               `json:"max_ratio,omitempty"`                             //Get the global share ratio limit
	MaxRatioAct                        *int                   `json:"max_ratio_act,omitempty"`                         //Action performed when a torrent reaches the maximum share ratio. See list of possible values here below.
	ListenPort                         *int                   `json:"listen_port,omitempty"`                           //Port for incoming connections
	Upnp                               *bool                  `json:"upnp,omitempty"`                                  //True if UPnP/NAT-PMP is enabled
	RandomPort                         *bool                  `json:"random_port,omitempty"`                           //True if the port is randomly selected
	DlLimit                            *int                   `json:"dl_limit,omitempty"`                              //Global download speed limit in KiB/s; -1 means no limit is applied
	UpLimit                            *int                   `json:"up_limit,omitempty"`                              //Global upload speed limit in KiB/s; -1 means no limit is applied
	MaxConnec                          *int                   `json:"max_connec,omitempty"`                            //Maximum global number of simultaneous connections
	MaxConnecPerTorrent                *int                   `json:"max_connec_per_torrent,omitempty"`                //Maximum number of simultaneous connections per torrent
	MaxUploads                         *int                   `json:"max_uploads,omitempty"`                           //Maximum number of upload slots
	MaxUploadsPerTorrent               *int                   `json:"max_uploads_per_torrent,omitempty"`               //Maximum number of upload slots per torrent
	StopTrackerTimeout                 *int                   `json:"stop_tracker_timeout,omitempty"`                  //Timeout in seconds for a stopped announce request to trackers
	EnablePieceExtentAffinity          *bool                  `json:"enable_piece_extent_affinity,omitempty"`          //True if the advanced libtorrent option piece_extent_affinity is enabled
	BittorrentProtocol                 *int                   `json:"bittorrent_protocol,omitempty"`                   //Bittorrent Protocol to use (see list of possible values below)
	LimitUtpRate                       *bool                  `json:"limit_utp_rate,omitempty"`                        //True if [du]l_limit should be applied to uTP connections; this option is only available in qBittorent built against libtorrent version 0.16.X and higher
	LimitTcpOverhead                   *bool                  `json:"limit_tcp_overhead,omitempty"`                    //True if [du]l_limit should be applied to estimated TCP overhead (service data: e.g. packet headers)
	LimitLanPeers                      *bool                  `json:"limit_lan_peers,omitempty"`                       //True if [du]l_limit should be applied to peers on the LAN
	AltDlLimit                         *int                   `json:"alt_dl_limit,omitempty"`                          //Alternative global download speed limit in KiB/s
	AltUpLimit                         *int                   `json:"alt_up_limit,omitempty"`                          //Alternative global upload speed limit in KiB/s
	SchedulerEnabled                   *bool                  `json:"scheduler_enabled,omitempty"`                     //True if alternative limits should be applied according to schedule
	ScheduleFromHour                   *int                   `json:"schedule_from_hour,omitempty"`                    //Scheduler starting hour
	ScheduleFromMin                    *int                   `json:"schedule_from_min,omitempty"`                     //Scheduler starting minute
	ScheduleToHour                     *int                   `json:"schedule_to_hour,omitempty"`                      //Scheduler ending hour
	ScheduleToMin                      *int                   `json:"schedule_to_min,omitempty"`                       //Scheduler ending minute
	SchedulerDays                      *int                   `json:"scheduler_days,omitempty"`                        //Scheduler days. See possible values here below
	Dht                                *bool                  `json:"dht,omitempty"`                                   //True if DHT is enabled
	Pex                                *bool                  `json:"pex,omitempty"`                                   //True if PeX is enabled
	Lsd                                *bool                  `json:"lsd,omitempty"`                                   //True if LSD is enabled
	Encryption                         *int                   `json:"encryption,omitempty"`                            //See list of possible values here below
	AnonymousMode                      *bool                  `json:"anonymous_mode,omitempty"`                        //If true anonymous mode will be enabled; read more here; this option is only available in qBittorent built against libtorrent version 0.16.X and higher
	ProxyType                          *int                   `json:"proxy_type,omitempty"`                            //See list of possible values here below
	ProxyIp                            *string                `json:"proxy_ip,omitempty"`                              //Proxy IP address or domain name
	ProxyPort                          *int                   `json:"proxy_port,omitempty"`                            //Proxy port
	ProxyPeerConnections               *bool                  `json:"proxy_peer_connections,omitempty"`                //True if peer and web seed connections should be proxified; this option will have any effect only in qBittorent built against libtorrent version 0.16.X and higher
	ProxyAuthEnabled                   *bool                  `json:"proxy_auth_enabled,omitempty"`                    //True proxy requires authentication; doesn't apply to SOCKS4 proxies
	ProxyUsername                      *string                `json:"proxy_username,omitempty"`                        //Username for proxy authentication
	ProxyPassword                      *string                `json:"proxy_password,omitempty"`                        //Password for proxy authentication
	ProxyTorrentsOnly                  *bool                  `json:"proxy_torrents_only,omitempty"`                   //True if proxy is only used for torrents
	IpFilterEnabled                    *bool                  `json:"ip_filter_enabled,omitempty"`                     //True if external IP filter should be enabled
	IpFilterPath                       *string                `json:"ip_filter_path,omitempty"`                        //Path to IP filter file (.dat, .p2p, .p2b files are supported); path is separated by slashes
	IpFilterTrackers                   *bool                  `json:"ip_filter_trackers,omitempty"`                    //True if IP filters are applied to trackers
	WebUiDomainList                    *string                `json:"web_ui_domain_list,omitempty"`                    //Comma-separated list of domains to accept when performing Host header validation
	WebUiAddress                       *string                `json:"web_ui_address,omitempty"`                        //IP address to use for the WebUI
	WebUiPort                          *int                   `json:"web_ui_port,omitempty"`                           //WebUI port
	WebUiUpnp                          *bool                  `json:"web_ui_upnp,omitempty"`                           //True if UPnP is used for the WebUI port
	WebUiUsername                      *string                `json:"web_ui_username,omitempty"`                       //WebUI username
	WebUiPassword                      *string                `json:"web_ui_password,omitempty"`                       //For API ≥ v2.3.0: Plaintext WebUI password, not readable, write-only. For API < v2.3.0: MD5 hash of WebUI password, hash is generated from the following string: username:Web UI Access:plain_text_web_ui_password
	WebUiCsrfProtectionEnabled         *bool                  `json:"web_ui_csrf_protection_enabled,omitempty"`        //True if WebUI CSRF protection is enabled
	WebUiClickjackingProtectionEnabled *bool                  `json:"web_ui_clickjacking_protection_enable,omitempty"` //True if WebUI clickjacking protection is enabled
	WebUiSecureCookieEnabled           *bool                  `json:"web_ui_secure_cookie_enabled,omitempty"`          //True if WebUI cookie Secure flag is enabled
	WebUiMaxAuthFailCount              *int                   `json:"web_ui_max_auth_fail_count,omitempty"`            //Maximum number of authentication failures before WebUI access ban
	WebUiBanDuration                   *int                   `json:"web_ui_ban_duration,omitempty"`                   //WebUI access ban duration in seconds
	WebUiSessionTimeout                *int                   `json:"web_ui_session_timeout,omitempty"`                //Seconds until WebUI is automatically signed off
	WebUiHostHeaderValidationEnabled   *bool                  `json:"web_ui_host_header_validation_enabled,omitempty"` //True if WebUI host header validation is enabled
	BypassLocalAuth                    *bool                  `json:"bypass_local_auth,omitempty"`                     //True if authentication challenge for loopback address (127.0.0.1) should be disabled
	BypassAuthSubnetWhitelistEnabled   *bool                  `json:"bypass_auth_subnet_whitelist_enabled,omitempty"`  //True if webui authentication should be bypassed for clients whose ip resides within (at least) one of the subnets on the whitelist
	BypassAuthSubnetWhitelist          *string                `json:"bypass_auth_subnet_whitelist,omitempty"`          //(White)list of ipv4/ipv6 subnets for which webui authentication should be bypassed; list entries are separated by commas
	AlternativeWebuiEnabled            *bool                  `json:"alternative_webui_enabled,omitempty"`             //True if an alternative WebUI should be used
	AlternativeWebuiPath               *string                `json:"alternative_webui_path,omitempty"`                //File path to the alternative WebUI
	UseHttps                           *bool                  `json:"use_https,omitempty"`                             //True if WebUI HTTPS access is enabled
	SslKey                             *string                `json:"ssl_key,omitempty"`                               //For API < v2.0.1: SSL keyfile contents (this is a not a path)
	SslCert                            *string                `json:"ssl_cert,omitempty"`                              //For API < v2.0.1: SSL certificate contents (this is a not a path)
	WebUiHttpsKeyPath                  *string                `json:"web_ui_https_key_path,omitempty"`                 //For API ≥ v2.0.1: Path to SSL keyfile
	WebUiHttpsCertPath                 *string                `json:"web_ui_https_cert_path,omitempty"`                //For API ≥ v2.0.1: Path to SSL certificate
	DyndnsEnabled                      *bool                  `json:"dyndns_enabled,omitempty"`                        //True if server DNS should be updated dynamically
	DyndnsService                      *int                   `json:"dyndns_service,omitempty"`                        //See list of possible values here below
	DyndnsUsername                     *string                `json:"dyndns_username,omitempty"`                       //Username for DDNS service
	DyndnsPassword                     *string                `json:"dyndns_password,omitempty"`                       //Password for DDNS service
	DyndnsDomain                       *string                `json:"dyndns_domain,omitempty"`                         //Your DDNS domain name
	RssRefreshInterval                 *int                   `json:"rss_refresh_interval,omitempty"`                  //RSS refresh interval
	RssMaxArticlesPerFeed              *int                   `json:"rss_max_articles_per_feed,omitempty"`             //Max stored articles per RSS feed
	RssProcessingEnabled               *bool                  `json:"rss_processing_enabled,omitempty"`                //Enable processing of RSS feeds
	RssAutoDownloadingEnabled          *bool                  `json:"rss_auto_downloading_enabled,omitempty"`          //Enable auto-downloading of torrents from the RSS feeds
	RssDownloadRepackProperEpisodes    *bool                  `json:"rss_download_repack_proper_episodes,omitempty"`   //For API ≥ v2.5.1: Enable downloading of repack/proper Episodes
	RssSmartEpisodeFilters             *string                `json:"rss_smart_episode_filters,omitempty"`             //For API ≥ v2.5.1: List of RSS Smart Episode Filters
	AddTrackersEnabled                 *bool                  `json:"add_trackers_enabled,omitempty"`                  //Enable automatic adding of trackers to new torrents
	AddTrackers                        *string                `json:"add_trackers,omitempty"`                          //List of trackers to add to new torrent
	WebUiUseCustomHttpHeadersEnabled   *bool                  `json:"web_ui_use_custom_http_headers_enable,omitempty"` //For API ≥ v2.5.1: Enable custom http headers
	WebUiCustomHttpHeaders             *string                `json:"web_ui_custom_http_headers,omitempty"`            //For API ≥ v2.5.1: List of custom http headers
	MaxSeedingTimeEnabled              *bool                  `json:"max_seeding_time_enabled,omitempty"`              //True enables max seeding time
	MaxSeedingTime                     *int                   `json:"max_seeding_time,omitempty"`                      //Number of minutes to seed a torrent
	AnnounceIp                         *string                `json:"announce_ip,omitempty"`                           //TODO
	AnnounceToAllTiers                 *bool                  `json:"announce_to_all_tiers,omitempty"`                 //True always announce to all tiers
	AnnounceToAllTrackers              *bool                  `json:"announce_to_all_trackers,omitempty"`              //True always announce to all trackers in a tier
	AsyncIoThreads                     *int                   `json:"async_io_threads,omitempty"`                      //Number of asynchronous I/O threads
	BannedIps                          *string                `json:"banned_IPs,omitempty"`                            //List of banned IPs
	CheckingMemoryUse                  *int                   `json:"checking_memory_use,omitempty"`                   //Outstanding memory when checking torrents in MiB
	CurrentInterfaceAddress            *string                `json:"current_interface_address,omitempty"`             //IP Address to bind to. Empty String means All addresses
	CurrentNetworkInterface            *string                `json:"current_network_interface,omitempty"`             //Network Interface used
	DiskCache                          *int                   `json:"disk_cache,omitempty"`                            //Disk cache used in MiB
	DiskCacheTtl                       *int                   `json:"disk_cache_ttl,omitempty"`                        //Disk cache expiry interval in seconds
	EmbeddedTrackerPort                *int                   `json:"embedded_tracker_port,omitempty"`                 //Port used for embedded tracker
	EnableCoalesceReadWrite            *bool                  `json:"enable_coalesce_read_write,omitempty"`            //True enables coalesce reads & writes
	EnableEmbeddedTracker              *bool                  `json:"enable_embedded_tracker,omitempty"`               //True enables embedded tracker
	EnableMultiConnectionsFromSameIp   *bool                  `json:"enable_multi_connections_from_same_ip,omitempty"` //True allows multiple connections from the same IP address
	EnableOsCache                      *bool                  `json:"enable_os_cache,omitempty"`                       //True enables os cache
	EnableUploadSuggestions            *bool                  `json:"enable_upload_suggestions,omitempty"`             //True enables sending of upload piece suggestions
	FilePoolSize                       *int                   `json:"file_pool_size,omitempty"`                        //File pool size
	OutgoingPortsMax                   *int                   `json:"outgoing_ports_max,omitempty"`                    //Maximal outgoing port (0: Disabled)
	OutgoingPortsMin                   *int                   `json:"outgoing_ports_min,omitempty"`                    //Minimal outgoing port (0: Disabled)
	RecheckCompletedTorrents           *bool                  `json:"recheck_completed_torrents,omitempty"`            //True rechecks torrents on completion
	ResolvePeerCountries               *bool                  `json:"resolve_peer_countries,omitempty"`                //True resolves peer countries
	SaveResumeDataInterval             *int                   `json:"save_resume_data_interval,omitempty"`             //Save resume data interval in min
	SendBufferLowWatermark             *int                   `json:"send_buffer_low_watermark,omitempty"`             //Send buffer low watermark in KiB
	SendBufferWatermark                *int                   `json:"send_buffer_watermark,omitempty"`                 //Send buffer watermark in KiB
	SendBufferWatermarkFactor          *int                   `json:"send_buffer_watermark_factor,omitempty"`          //Send buffer watermark factor in percent
	SocketBacklogSize                  *int                   `json:"socket_backlog_size,omitempty"`                   //Socket backlog size
	UploadChokingAlgorithm             *int                   `json:"upload_choking_algorithm,omitempty"`              //Upload choking algorithm used (see list of possible values below)
	UploadSlotsBehavior                *int                   `json:"upload_slots_behavior,omitempty"`                 //Upload slots behavior used (see list of possible values below)
	UpnpLeaseDuration                  *int                   `json:"upnp_lease_duration,omitempty"`                   //UPnP lease duration (0: Permanent lease)
	UtpTcpMixedMode                    *int                   `json:"utp_tcp_mixed_mode,omitempty"`
}

type SetApplicationPreferencesRsp struct {
}

type GetDefaultSavePathReq struct {
}

type GetDefaultSavePathRsp struct {
	Path string
}

type GetLogReq struct {
	Normal      bool `json:"normal"`
	Info        bool `json:"info"`
	Warning     bool `json:"warning"`
	Critical    bool `json:"critical"`
	LastKnownId int  `json:"last_known_id"`
}

type LogItem struct {
	Id        int    `json:"id"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
	Type      int    `json:"type"`
}

type GetLogRsp struct {
	Items []*LogItem
}

type GetPeerLogReq struct {
	LastKnownId int `json:"last_known_id"`
}

type PeerLogItem struct {
	Id        int    `json:"id"`
	Ip        string `json:"ip"`
	Timestamp int    `json:"timestamp"`
	Blocked   bool   `json:"blocked"`
	Reason    string `json:"reason"`
}

type GetPeerLogRsp struct {
	Items []*PeerLogItem
}

type GetMainDataReq struct {
	Rid int `json:"rid"`
}

type GetMainDataRsp struct {
	Rid               int                         `json:"rid"`                //Response ID
	FullUpdate        bool                        `json:"full_update"`        //Whether the response contains all the data or partial data
	Torrents          map[string]*TorrentListItem `json:"torrents"`           //Property: torrent hash, value: same as torrent list
	TorrentsRemoved   []interface{}               `json:"torrents_removed"`   //List of hashes of torrents removed since last request
	Categories        map[string]interface{}      `json:"categories"`         //Info for categories added since last request
	CategoriesRemoved []interface{}               `json:"categories_removed"` //List of categories removed since last request
	Tags              []interface{}               `json:"tags"`               //List of tags added since last request
	TagsRemoved       []interface{}               `json:"tags_removed"`       //List of tags removed since last request
	ServerState       *GlobalTransferInfo         `json:"server_state"`       //Global transfer info
}

type GetTorrentPeerDataReq struct {
	Hash string `json:"hash"`
	Rid  int    `json:"rid"`
}

type TorrentPeerItem struct {
	Client      string  `json:"client"`
	Connection  string  `json:"connection"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	DlSpeed     int     `json:"dl_speed"`
	Downloaded  int     `json:"downloaded"`
	Files       string  `json:"files"`
	Flags       string  `json:"flags"`
	FlagsDesc   string  `json:"flags_desc"`
	Ip          string  `json:"ip"`
	Port        int     `json:"port"`
	Progress    float64 `json:"progress"`
	Relevance   float64 `json:"relevance"`
	UpSpeed     int     `json:"up_speed"`
	Uploaded    int     `json:"uploaded"`
}

type TorrentPeerData struct {
	FullUpdate bool                        `json:"full_update"`
	Rid        int                         `json:"rid"`
	ShowFlags  bool                        `json:"show_flags"`
	Peers      map[string]*TorrentPeerItem `json:"peers"`
}

type GetTorrentPeerDataRsp struct {
	Data *TorrentPeerData
}

type GetGlobalTransferInfoReq struct {
}

type GlobalTransferInfo struct {
	DlInfoSpeed      int    `json:"dl_info_speed"`     //Global download rate (bytes/s)
	DlInfoData       int    `json:"dl_info_data"`      //Data downloaded this session (bytes)
	UpInfoSpeed      int    `json:"up_info_speed"`     //Global upload rate (bytes/s)
	UpInfoData       int    `json:"up_info_data"`      //Data uploaded this session (bytes)
	DlRateLimit      int    `json:"dl_rate_limit"`     //Download rate limit (bytes/s)
	UpRateLimit      int    `json:"up_rate_limit"`     //Upload rate limit (bytes/s)
	DhtNodes         int    `json:"dht_nodes"`         //DHT nodes connected to
	ConnectionStatus string `json:"connection_status"` //Connection status. See possible values here below
}

type GetGlobalTransferInfoRsp struct {
	Info *GlobalTransferInfo
}

type GetAlternativeSpeedLimitsStateReq struct {
}

type GetAlternativeSpeedLimitsStateRsp struct {
	Enabled bool
}

type ToggleAlternativeSpeedLimitsReq struct {
}

type ToggleAlternativeSpeedLimitsRsp struct {
}

type GetGlobalDownloadLimitReq struct {
}

type GetGlobalDownloadLimitRsp struct {
	Speed int
}

type SetGlobalDownloadLimitReq struct {
	Speed int `json:"limit"`
}

type SetGlobalDownloadLimitRsp struct {
}

type GetGlobalUploadLimitReq struct {
}

type GetGlobalUploadLimitRsp struct {
	Speed int
}

type SetGlobalUploadLimitReq struct {
	Speed int `json:"limit"`
}

type SetGlobalUploadLimitRsp struct {
}

type BanPeersReq struct {
	Peers []string
}

type banPeersReqInner struct {
	Peers string `json:"peers"`
}

type BanPeersRsp struct {
}

type GetTorrentGenericPropertiesReq struct {
	Hash string `json:"hash"`
}

type TorrentGenericProperty struct {
	SavePath               string  `json:"save_path"`                //Torrent save path
	CreationDate           int     `json:"creation_date"`            //Torrent creation date (Unix timestamp)
	PieceSize              int     `json:"piece_size"`               //Torrent piece size (bytes)
	Comment                string  `json:"comment"`                  //Torrent comment
	TotalWasted            int     `json:"total_wasted"`             //Total data wasted for torrent (bytes)
	TotalUploaded          int     `json:"total_uploaded"`           //Total data uploaded for torrent (bytes)
	TotalUploadedSession   int     `json:"total_uploaded_session"`   //Total data uploaded this session (bytes)
	TotalDownloaded        int     `json:"total_downloaded"`         //Total data downloaded for torrent (bytes)
	TotalDownloadedSession int     `json:"total_downloaded_session"` //Total data downloaded this session (bytes)
	UpLimit                int     `json:"up_limit"`                 //Torrent upload limit (bytes/s)
	DlLimit                int     `json:"dl_limit"`                 //Torrent download limit (bytes/s)
	TimeElapsed            int     `json:"time_elapsed"`             //Torrent elapsed time (seconds)
	SeedingTime            int     `json:"seeding_time"`             //Torrent elapsed time while complete (seconds)
	NbConnections          int     `json:"nb_connections"`           //Torrent connection count
	NbConnectionsLimit     int     `json:"nb_connections_limit"`     //Torrent connection count limit
	ShareRatio             float64 `json:"share_ratio"`              //Torrent share ratio
	AdditionDate           int     `json:"addition_date"`            //When this torrent was added (unix timestamp)
	CompletionDate         int     `json:"completion_date"`          //Torrent completion date (unix timestamp)
	CreatedBy              string  `json:"created_by"`               //Torrent creator
	DlSpeedAvg             int     `json:"dl_speed_avg"`             //Torrent average download speed (bytes/second)
	DlSpeed                int     `json:"dl_speed"`                 //Torrent download speed (bytes/second)
	Eta                    int     `json:"eta"`                      //Torrent ETA (seconds)
	LastSeen               int     `json:"last_seen"`                //Last seen complete date (unix timestamp)
	Peers                  int     `json:"peers"`                    //Number of peers connected to
	PeersTotal             int     `json:"peers_total"`              //Number of peers in the swarm
	PiecesHave             int     `json:"pieces_have"`              //Number of pieces owned
	PiecesNum              int     `json:"pieces_num"`               //Number of pieces of the torrent
	Reannounce             int     `json:"reannounce"`               //Number of seconds until the next announce
	Seeds                  int     `json:"seeds"`                    //Number of seeds connected to
	SeedsTotal             int     `json:"seeds_total"`              //Number of seeds in the swarm
	TotalSize              int     `json:"total_size"`               //Torrent total size (bytes)
	UpSpeedAvg             int     `json:"up_speed_avg"`             //Torrent average upload speed (bytes/second)
	UpSpeed                int     `json:"up_speed"`                 //Torrent upload speed (bytes/second)
}

type GetTorrentGenericPropertiesRsp struct {
	Property *TorrentGenericProperty
}

type GetTorrentTrackersReq struct {
	Hash string `json:"hash"`
}

//Note: tier should be integer, but in some trackers, it returns empty string
type TorrentTrackerItem struct {
	Url           string      `json:"url"`            //Tracker url
	Status        int         `json:"status"`         //Tracker status. See the table below for possible values
	Tier          interface{} `json:"tier"`           //Tracker priority tier. Lower tier trackers are tried before higher tiers
	NumPeers      int         `json:"num_peers"`      //Number of peers for current torrent, as reported by the tracker
	NumSeeds      int         `json:"num_seeds"`      //Number of seeds for current torrent, asreported by the tracker
	NumLeeches    int         `json:"num_leeches"`    //Number of leeches for current torrent, as reported by the tracker
	NumDownloaded int         `json:"num_downloaded"` //Number of completed downlods for current torrent, as reported by the tracker
	Msg           string      `json:"msg"`            //Tracker message (there is no way of knowing what this message is - it's up to tracker admins)
}

type GetTorrentTrackersRsp struct {
	Trackers []*TorrentTrackerItem
}

type GetTorrentWebSeedsReq struct {
	Hash string `json:"hash"`
}

type TorrentWebSeedItem struct {
	Url string `json:"url"`
}

type GetTorrentWebSeedsRsp struct {
	WebSeeds []*TorrentWebSeedItem
}

type GetTorrentContentsReq struct {
	Hash  string
	Index []string
}

type getTorrentContentsInnerReq struct {
	Hash    string  `json:"hash"`
	Indexes *string `json:"indexes"`
}

type TorrentContentItem struct {
	Index        int     `json:"index"`        //File index
	Name         string  `json:"name"`         //File name (including relative path)
	Size         int     `json:"size"`         //File size (bytes)
	Progress     float64 `json:"progress"`     //File progress (percentage/100)
	Priority     int     `json:"priority"`     //File priority. See possible values here below
	IsSeed       bool    `json:"is_seed"`      //True if file is seeding/complete
	PieceRange   []int   `json:"piece_range"`  //The first number is the starting piece index and the second number is the ending piece index (inclusive)
	Availability float64 `json:"availability"` //Percentage of file pieces currently available (percentage/100)

}
type GetTorrentContentsRsp struct {
	Contents []*TorrentContentItem
}

type GetTorrentPiecesStatesReq struct {
	Hash string `json:"hash"`
}

type GetTorrentPiecesStatesRsp struct {
	States []int
}

type GetTorrentPiecesHashesReq struct {
	Hash string `json:"hash"`
}

type GetTorrentPiecesHashesRsp struct {
	Hashes []string
}

type PauseTorrentsReq struct {
	Hash []string
}

type pauseTorrentsInnerReq struct {
	Hashes string `json:"hashes"`
}

type PauseTorrentsRsp struct {
}

type ResumeTorrentsReq struct {
	Hash []string
}

type resumeTorrentsInnerReq struct {
	Hashes string `json:"hashes"`
}

type ResumeTorrentsRsp struct {
}

type DeleteTorrentsReq struct {
	IsDeleteAll  bool
	IsDeleteFile bool
	Hash         []string //use only IsDeleteAll set to false
}

type deleteTorrentsInnerReq struct {
	Hashes      string `json:"hashes"`
	DeleteFiles bool   `json:"deleteFiles"`
}

type DeleteTorrentsRsp struct {
}

type RecheckTorrentsReq struct {
	IsRecheckAll bool
	Hash         []string
}

type recheckTorrentsInnerReq struct {
	Hashes string `json:"hashes"`
}

type RecheckTorrentsRsp struct {
}

type ReannounceTorrentsReq struct {
	IsReannounceAll bool
	Hash            []string
}

type reannounceTorrentsInnerReq struct {
	Hashes string `json:"hashes"`
}

type ReannounceTorrentsRsp struct {
}

type AddTorrentMeta struct {
	Savepath           *string  `json:"savepath"`           //Download folder
	Cookie             *string  `json:"cookie"`             //Cookie sent to download the .torrent file
	Category           *string  `json:"category"`           //Category for the torrent
	Tags               string   `json:"tags"`               //Tags for the torrent, split by ','
	SkipChecking       *bool    `json:"skip_checking"`      //Skip hash checking. Possible values are true, false (default)
	Paused             *bool    `json:"paused"`             //Add torrents in the paused state. Possible values are true, false (default)
	RootFolder         *bool    `json:"root_folder"`        //Create the root folder. Possible values are true, false, unset (default)
	Rename             *string  `json:"rename"`             //Rename torrent
	UpLimit            *int     `json:"upLimit"`            //Set torrent upload speed limit. Unit in bytes/second
	DlLimit            *int     `json:"dlLimit"`            //Set torrent download speed limit. Unit in bytes/second
	RatioLimit         *float64 `json:"ratioLimit"`         //Set torrent share ratio limit
	SeedingTimeLimit   *int     `json:"seedingTimeLimit"`   //Set torrent seeding time limit. Unit in seconds
	AutoTMM            *bool    `json:"autoTMM"`            //Whether Automatic Torrent Management should be used
	SequentialDownload *string  `json:"sequentialDownload"` //Enable sequential download. Possible values are true, false (default)
	FirstLastPiecePrio *string  `json:"firstLastPiecePrio"` //Prioritize download first last piece. Possible values are true, false (default)
}

type AddNewTorrentReq struct {
	File []string
	Meta *AddTorrentMeta
}

type AddNewTorrentRsp struct {
}

type AddNewLinkReq struct {
	Url  []string
	Meta *AddTorrentMeta
}

type AddNewLinkRsp struct {
}

type AddTrackersToTorrentReq struct {
	Hash string
	Url  []string
}

type addTrackersToTorrentInnerReq struct {
	Hash string `json:"hash"`
	Urls string `json:"urls"`
}

type AddTrackersToTorrentRsp struct {
}

type EditTrackersReq struct {
	Hash    string `json:"hash"`
	OrigUrl string `json:"origUrl"`
	NewUrl  string `json:"newUrl"`
}

type EditTrackersRsp struct {
}

type RemoveTrackersReq struct {
	Hash string
	Url  []string
}

type removeTrackersInnerReq struct {
	Hash string `json:"hash"`
	Urls string `json:"urls"`
}

type RemoveTrackersRsp struct {
}

type AddPeersReq struct {
	Hash []string
	Peer []string
}

type addPeersInnerReq struct {
	Hashes string `json:"hashes"`
	Peers  string `json:"peers"`
}

type AddPeersRsp struct {
}

type IncreaseTorrentPriorityReq struct {
	IsIncreaseAllTorrent bool
	Hash                 []string
}

type increaseTorrentPriorityInnerReq struct {
	Hashes string `json:"hashes"`
}

type IncreaseTorrentPriorityRsp struct {
}

type DecreaseTorrentPriorityReq struct {
	IsDecreaseAllTorrent bool
	Hash                 []string
}

type decreaseTorrentPriorityInnerReq struct {
	Hashes string `json:"hashes"`
}

type DecreaseTorrentPriorityRsp struct {
}

type MaximalTorrentPriorityReq struct {
	IsMaximalAllTorrent bool
	Hash                []string
}

type maximalTorrentPriorityInnerReq struct {
	Hashes string `json:"hashes"`
}

type MaximalTorrentPriorityRsp struct {
}

type MinimalTorrentPriorityReq struct {
	IsMinimalAllTorrent bool
	Hash                []string
}

type minimalTorrentPriorityInnerReq struct {
	Hashes string `json:"hashes"`
}

type MinimalTorrentPriorityRsp struct {
}

type SetFilePriorityReq struct {
	Hash     string
	Id       []string
	Priority FilePriority
}

type setFilePriorityInnerReq struct {
	Hash     string `json:"hash"`     //The hash of the torrent
	Id       string `json:"id"`       //File ids, separated by |
	Priority int    `json:"priority"` //File priority to set (consult torrent contents API for possible values)
}

type SetFilePriorityRsp struct {
}

type GetTorrentDownloadLimitReq struct {
	IsGetAll bool
	Hash     []string
}

type getTorrentDownloadLimitInnerReq struct {
	Hashes string `json:"hashes"`
}

type GetTorrentDownloadLimitRsp struct {
	SpeedMap map[string]int //{hash:speed, ...}
}

type SetTorrentDownloadLimitReq struct {
	IsSetAll bool
	Hash     []string
	Speed    int
}

type setTorrentDownloadLimitInnerReq struct {
	Hashes string `json:"hashes"`
	Speed  int    `json:"limit"`
}

type SetTorrentDownloadLimitRsp struct {
}

type SetTorrentShareLimitReq struct {
	IsSetAll         bool
	Hash             []string
	SeedingTimeLimit int
	RatioLimit       float64
}

type setTorrentShareLimitInnerReq struct {
	Hashes           string  `json:"hashes"`
	SeedingTimeLimit int     `json:"seedingTimeLimit"`
	RatioLimit       float64 `json:"ratioLimit"`
}

type SetTorrentShareLimitRsp struct {
}

type GetTorrentUploadLimitReq struct {
	IsGetAll bool
	Hash     []string
}

type getTorrentUploadLimitInnerReq struct {
	Hashes string `json:"hashes"`
}

type GetTorrentUploadLimitRsp struct {
	SpeedMap map[string]int
}

type SetTorrentUploadLimitReq struct {
	IsSetAll bool
	Hash     []string
	Speed    int
}

type setTorrentUploadLimitInnerReq struct {
	Hashes string `json:"hashes"`
	Speed  int    `json:"limit"`
}

type SetTorrentUploadLimitRsp struct {
}

type SetTorrentLocationReq struct {
	IsSetAll bool
	Hash     []string
	Location string
}

type setTorrentLocationInnerReq struct {
	Hashes   string `json:"hashes"`
	Location string `json:"location"`
}

type SetTorrentLocationRsp struct {
}

type SetTorrentNameReq struct {
	Hash string `json:"hash"`
	Name string `json:"name"`
}

type SetTorrentNameRsp struct {
}

type SetTorrentCategoryReq struct {
	IsSetAll bool
	Hash     []string
	Category string
}

type setTorrentCategoryInnerReq struct {
	Hashes   string `json:"hashes"`
	Category string `json:"category"`
}

type SetTorrentCategoryRsp struct {
}

type GetAllCategoriesReq struct {
}

type CategoryInfo struct {
	Name     string `json:"name"`
	SavePath string `json:"savePath"`
}

type GetAllCategoriesRsp struct {
	Categories map[string]*CategoryInfo
}

type AddNewCategoryReq struct {
	Category string `json:"category"`
	SavePath string `json:"savePath"`
}

type AddNewCategoryRsp struct {
}

type EditCategoryReq struct {
	Category string `json:"category"`
	SavePath string `json:"savePath"`
}

type EditCategoryRsp struct {
}

type RemoveCategoriesReq struct {
	Category []string
}

type removeCategoriesInnerReq struct {
	Categories string `json:"categories"`
}

type RemoveCategoriesRsp struct {
}

type AddTorrentTagsReq struct {
	IsAddAll bool
	Hash     []string
	Tag      []string
}

type addTorrentTagsInnerReq struct {
	Hashes string `json:"hashes"`
	Tags   string `json:"tags"`
}

type AddTorrentTagsRsp struct {
}

type RemoveTorrentTagsReq struct {
	IsRemoveAll bool
	Hash        []string
	Tag         []string
}

type removeTorrentTagsInnerReq struct {
	Hashes string `json:"hashes"`
	Tags   string `json:"tags"`
}

type RemoveTorrentTagsRsp struct {
}

type GetAllTagsReq struct {
}

type GetAllTagsRsp struct {
	Tags []string
}

type CreateTagsReq struct {
	Tag []string
}

type createTagsInnerReq struct {
	Tags string `json:"tags"`
}

type CreateTagsRsp struct {
}

type DeleteTagsReq struct {
	Tag []string
}

type deleteTagsInnerReq struct {
	Tags string `json:"tags"`
}

type DeleteTagsRsp struct {
}

type SetAutomaticTorrentManagementReq struct {
	IsSetAll bool
	Hash     []string
	Enable   bool
}

type setAutomaticTorrentManagementInnerReq struct {
	Hashes string `json:"hashes"`
	Enable bool   `json:"enable"`
}

type SetAutomaticTorrentManagementRsp struct {
}

type ToggleSequentialDownloadReq struct {
	IsSetAll bool
	Hash     []string
}

type toggleSequentialDownloadInnerReq struct {
	Hashes string `json:"hashes"`
}

type ToggleSequentialDownloadRsp struct {
}

type SetFirstOrLastPiecePriorityReq struct {
	IsSetAll bool
	Hash     []string
}

type setFirstOrLastPiecePriorityInnerReq struct {
	Hashes string `json:"hashes"`
}

type SetFirstOrLastPiecePriorityRsp struct {
}

type SetForceStartReq struct {
	IsSetAll bool
	Hash     []string
	Value    bool
}

type setForceStartInnerReq struct {
	Hashes string `json:"hashes"`
	Value  bool   `json:"value"`
}

type SetForceStartRsp struct {
}

type SetSuperSeedingReq struct {
	IsSetAll bool
	Hash     []string
	Value    bool
}

type setSuperSeedingInnerReq struct {
	Hashes string `json:"hashes"`
	Value  bool   `json:"value"`
}

type SetSuperSeedingRsp struct {
}

type RenameFileReq struct {
	Hash    string `json:"hash"`    //The hash of the torrent
	OldPath string `json:"oldPath"` //The old path of the torrent
	NewPath string `json:"newPath"` //The new path to use for the file
}

type RenameFileRsp struct {
}

type RenameFolderReq struct {
	Hash    string `json:"hash"`    //The hash of the torrent
	OldPath string `json:"oldPath"` //The old path of the torrent
	NewPath string `json:"newPath"` //The new path to use for the file
}

type RenameFolderRsp struct {
}
