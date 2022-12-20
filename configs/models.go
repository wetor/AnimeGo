package configs

type Config struct {
	Version  string `yaml:"version" json:"version" comment:"配置文件版本"`
	Setting  `yaml:"setting" json:"setting" comment:"常规设置"`
	Advanced `yaml:"advanced" json:"advanced" comment:"高级设置"`
}
type Setting struct {
	Feed struct {
		Mikan struct {
			Name string `yaml:"name" json:"name"`
			Url  string `yaml:"url" json:"url" comment:"Mikan订阅链接，为空则不使用自动订阅"`
		} `yaml:"mikan" json:"mikan" comment:"Mikan Project(mikanani.me)订阅信息"`
	} `yaml:"feed" json:"feed" comment:"自动订阅设置"`
	Client struct {
		QBittorrent struct {
			Url      string `yaml:"url" json:"url"`
			Username string `yaml:"username" json:"username"`
			Password string `yaml:"password" json:"password"`
		} `yaml:"qbittorrent" json:"qbittorrent" comment:"qBittorrent客户端信息"`
	} `yaml:"client" json:"client" comment:"下载客户端设置"`
	DownloadPath string `yaml:"download_path" json:"download_path" comment:"下载到的根目录"`
	SavePath     string `yaml:"save_path" json:"save_path" comment:"下载完成后移动到的根目录"`
	DataPath     string `yaml:"data_path" json:"data_path" comment:"数据存储根目录，用于保存数据库、插件等数据"`
	Filter       struct {
		JavaScript []string `yaml:"javascript" json:"javascript" comment_key:"filter_javascript"`
	} `yaml:"filter" json:"filter" comment:"过滤器设置，用来筛选符合条件的项目进行解析下载"`
	Category string `yaml:"category" json:"category" comment:"分类名，仅qBittorrent有效"`
	TagSrc   string `yaml:"tag" json:"tag" comment_key:"tag_help"`
	WebApi   struct {
		AccessKey string `yaml:"access_key" json:"access_key" comment:"请求秘钥，为空则不需要验证"`
		Host      string `yaml:"host" json:"host"`
		Port      int    `yaml:"port" json:"port"`
	} `yaml:"webapi" json:"webapi" comment:"WebApi设置"`
	Proxy struct {
		Enable bool   `yaml:"enable" json:"enable" comment:"开启后AnimeGo所有的网络请求都会使用代理"`
		Url    string `yaml:"url" json:"url" comment:"支持http、https和socks5代理"`
	} `yaml:"proxy" json:"proxy" comment:"代理设置"`
	Key struct {
		Themoviedb string `yaml:"themoviedb" json:"themoviedb" comment_key:"themoviedb_key"`
	} `yaml:"key" json:"key" comment:"秘钥设置"`
}

type Advanced struct {
	UpdateDelaySecond int `yaml:"update_delay_second" json:"update_delay_second"`

	Request struct {
		TimeoutSecond   int `yaml:"timeout_second" json:"timeout_second" comment:"请求超时时间"`
		RetryNum        int `yaml:"retry_num" json:"retry_num" comment:"额外重试次数"`
		RetryWaitSecond int `yaml:"retry_wait_second" json:"retry_wait_second" comment:"重试间隔等待时间"`
	} `yaml:"request" json:"request" comment:"网络请求设置"`

	Download struct {
		AllowDuplicateDownload bool   `yaml:"allow_duplicate_download" json:"allow_duplicate_download" comment:"允许重复下载同剧集不同资源"`
		SeedingTimeMinute      int    `yaml:"seeding_time_minute" json:"seeding_time_minute" comment:"做种时间"`
		IgnoreSizeMaxKb        int    `yaml:"ignore_size_max_kb" json:"ignore_size_max_kb" comment:"忽略小文件大小"`
		Rename                 string `yaml:"rename" json:"rename" comment_key:"rename"`
	} `yaml:"download" json:"download" comment:"下载设置"`

	Feed struct {
		UpdateDelayMinute int `yaml:"update_delay_minute" json:"update_delay_minute" comment:"订阅刷新时间"`
		DelaySecond       int `yaml:"delay_second" json:"delay_second" comment:"订阅解析间隔时间，防止高频请求"`
		MultiGoroutine    struct {
			Enable       bool `yaml:"enable" json:"enable" comment:"多协程解析是否启用"`
			GoroutineMax int  `yaml:"goroutine_max" json:"goroutine_max" comment:"多协程解析最大协程数量"`
		} `yaml:"multi_goroutine" json:"multi_goroutine" comment:"订阅多协程解析"`
	} `yaml:"feed" json:"feed" comment:"订阅设置"`

	Path struct {
		DbFile   string `yaml:"db_file" json:"db_file" comment:"数据库保存文件名"`
		LogFile  string `yaml:"log_file" json:"log_file" comment:"日志保存文件名，日志会在所在文件夹自动归档"`
		TempPath string `yaml:"temp_path" json:"temp_path" comment:"临时文件保存文件夹"`
	} `yaml:"path" json:"path" comment:"其他路径设置，路径相对于data_path"`

	Default struct {
		TMDBFailSkip           bool `yaml:"tmdb_fail_skip" json:"tmdb_fail_skip" comment:"tmdb解析失败时，跳过此条目。优先级3"`
		TMDBFailUseTitleSeason bool `yaml:"tmdb_fail_use_title_season" json:"tmdb_fail_use_title_season" comment:"tmdb解析失败时，从文件名中获取季度信息。优先级2"`
		TMDBFailUseFirstSeason bool `yaml:"tmdb_fail_use_first_season" json:"tmdb_fail_use_first_season" comment:"tmdb解析失败时，默认使用第一季。优先级1"`
	} `yaml:"default" json:"default" comment:"默认值开关设置，同类型默认值按优先级执行。数值越大，优先级越高"`

	Client struct {
		ConnectTimeoutSecond int `yaml:"connect_timeout_second" json:"connect_timeout_second" comment:"连接超时时间"`
		RetryConnectNum      int `yaml:"retry_connect_num" json:"retry_connect_num" comment:"连接失败重试次数"`
		CheckTimeSecond      int `yaml:"check_time_second" json:"check_time_second" comment:"检查连接状态间隔时间，每次检查都会进行重试连接"`
	} `yaml:"client" json:"client" comment:"下载客户端设置"`

	Cache struct {
		MikanCacheHour      int `yaml:"mikan_cache_hour" json:"mikan_cache_hour" comment:"Mikan数据缓存时间，默认7*24小时(7天)。主要为mikan-id与bangumi-id的映射关系"`
		BangumiCacheHour    int `yaml:"bangumi_cache_hour" json:"bangumi_cache_hour" comment:"Bangumi数据缓存时间，默认3*24小时(3天)。主要为bangumi-id与详细信息的映射"`
		ThemoviedbCacheHour int `yaml:"themoviedb_cache_hour" json:"themoviedb_cache_hour" comment:"Themoviedb数据缓存时间，默认14*24小时(14天)。主要为tmdb-id与季度信息的映射"`
	} `yaml:"cache" json:"cache" comment:"缓存设置"`
}
