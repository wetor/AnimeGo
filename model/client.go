package model

type TorrentItem struct {
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

type Preferences struct {
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
