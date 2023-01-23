package configs

type Config struct {
	Version  string `yaml:"version" json:"version" attr:"配置文件版本"`
	Setting  `yaml:"setting" json:"setting" attr:"常规设置"`
	Advanced `yaml:"advanced" json:"advanced" attr:"高级设置"`
}

type Setting struct {
	Feed struct {
		Mikan struct {
			Name string `yaml:"name" json:"name" attr:"别名"`
			Url  string `yaml:"url" json:"url" attr:"订阅链接" comment:"可空，为空则不使用自动订阅"`
		} `yaml:"mikan" json:"mikan" attr:"Mikan订阅"`
	} `yaml:"feed" json:"feed" attr:"自动订阅设置"`
	Client struct {
		QBittorrent struct {
			Url      string `yaml:"url" json:"url" attr:"地址"`
			Username string `yaml:"username" json:"username" attr:"用户名"`
			Password string `yaml:"password" json:"password" attr:"密码"`
		} `yaml:"qbittorrent" json:"qbittorrent" attr:"qBittorrent客户端"`
	} `yaml:"client" json:"client" attr:"下载客户端设置"`
	DownloadPath string `yaml:"download_path" json:"download_path" attr:"下载文件夹" comment:"下载器的下载文件夹"`
	SavePath     string `yaml:"save_path" json:"save_path" attr:"保存文件夹" comment:"下载完成后，重命名并移动到的文件夹"`
	DataPath     string `yaml:"data_path" json:"data_path" attr:"数据文件夹" comment:"用于保存数据库、插件等数据"`
	Filter       struct {
		JavaScript []string `yaml:"javascript" json:"javascript" attr:"JavaScript插件" comment_key:"filter_javascript"`
	} `yaml:"filter" json:"filter" attr:"过滤器设置" comment:"用来筛选符合条件的项目进行解析下载"`
	Category string `yaml:"category" json:"category" attr:"分类名" comment:"仅qBittorrent有效"`
	Tag      string `yaml:"tag" json:"tag" attr:"标签表达式" comment_key:"tag_help"`
	WebApi   struct {
		AccessKey string `yaml:"access_key" json:"access_key" attr:"请求秘钥" comment:"为空则不需要验证"`
		Host      string `yaml:"host" json:"host" attr:"域名"`
		Port      int    `yaml:"port" json:"port" attr:"端口"`
	} `yaml:"webapi" json:"webapi" attr:"WebApi设置"`
	Proxy struct {
		Enable bool   `yaml:"enable" json:"enable" attr:"启用"`
		Url    string `yaml:"url" json:"url" attr:"代理链接"`
	} `yaml:"proxy" json:"proxy" attr:"代理设置" comment:"开启后AnimeGo所有的网络请求都会使用代理"`
	Key struct {
		Themoviedb string `yaml:"themoviedb" json:"themoviedb" attr:"TheMovieDB的APIkey" comment_key:"themoviedb_key"`
	} `yaml:"key" json:"key" attr:"秘钥设置"`
}

type Advanced struct {
	UpdateDelaySecond int `yaml:"update_delay_second" json:"update_delay_second" attr:"更新状态等待时间"`

	Request struct {
		TimeoutSecond   int `yaml:"timeout_second" json:"timeout_second" attr:"请求超时时间"`
		RetryNum        int `yaml:"retry_num" json:"retry_num" attr:"额外重试次数"`
		RetryWaitSecond int `yaml:"retry_wait_second" json:"retry_wait_second" attr:"重试间隔等待时间"`
	} `yaml:"request" json:"request" attr:"网络请求设置"`

	Download struct {
		AllowDuplicateDownload bool   `yaml:"allow_duplicate_download" json:"allow_duplicate_download" attr:"允许重复下载"`
		SeedingTimeMinute      int    `yaml:"seeding_time_minute" json:"seeding_time_minute" attr:"做种时间"`
		IgnoreSizeMaxKb        int    `yaml:"ignore_size_max_kb" json:"ignore_size_max_kb" attr:"忽略小文件大小"`
		Rename                 string `yaml:"rename" json:"rename" attr:"重命名方式" comment_key:"rename"`
	} `yaml:"download" json:"download" attr:"下载设置"`

	Feed struct {
		UpdateDelayMinute int `yaml:"update_delay_minute" json:"update_delay_minute" attr:"订阅刷新时间"`
		DelaySecond       int `yaml:"delay_second" json:"delay_second" attr:"订阅解析间隔时间"`
		MultiGoroutine    struct {
			Enable       bool `yaml:"enable" json:"enable" attr:"启用"`
			GoroutineMax int  `yaml:"goroutine_max" json:"goroutine_max" attr:"最大协程数量"`
		} `yaml:"multi_goroutine" json:"multi_goroutine" attr:"订阅多协程解析"`
	} `yaml:"feed" json:"feed" attr:"订阅设置"`

	Path struct {
		DbFile   string `yaml:"db_file" json:"db_file" attr:"数据库文件名"`
		LogFile  string `yaml:"log_file" json:"log_file" attr:"日志文件名" comment:"日志会在所在文件夹自动归档"`
		TempPath string `yaml:"temp_path" json:"temp_path" attr:"临时文件夹"`
	} `yaml:"path" json:"path" attr:"其他路径设置"`

	Default struct {
		TMDBFailSkip           bool `yaml:"tmdb_fail_skip" json:"tmdb_fail_skip" attr:"跳过当前项" comment:"tmdb解析季度失败时，跳过当前项。优先级3"`
		TMDBFailUseTitleSeason bool `yaml:"tmdb_fail_use_title_season" json:"tmdb_fail_use_title_season" attr:"文件名解析季度" comment:"tmdb解析季度失败时，从文件名中获取季度信息。优先级2"`
		TMDBFailUseFirstSeason bool `yaml:"tmdb_fail_use_first_season" json:"tmdb_fail_use_first_season" attr:"使用第一季" comment:"tmdb解析季度失败时，默认使用第一季。优先级1"`
	} `yaml:"default" json:"default" attr:"解析季度默认值" comment:"使用tmdb解析季度失败时，同类型默认值按优先级执行。数值越大，优先级越高"`

	Client struct {
		ConnectTimeoutSecond int `yaml:"connect_timeout_second" json:"connect_timeout_second" attr:"连接超时时间"`
		RetryConnectNum      int `yaml:"retry_connect_num" json:"retry_connect_num" attr:"连接失败重试次数"`
		CheckTimeSecond      int `yaml:"check_time_second" json:"check_time_second" attr:"检查连接状态间隔时间"`
	} `yaml:"client" json:"client" attr:"下载客户端设置"`

	Cache struct {
		MikanCacheHour      int `yaml:"mikan_cache_hour" json:"mikan_cache_hour" attr:"Mikan缓存时间" comment:"默认7*24小时(7天)。主要为mikan-id与bangumi-id的映射关系"`
		BangumiCacheHour    int `yaml:"bangumi_cache_hour" json:"bangumi_cache_hour" attr:"Bangumi缓存时间" comment:"默认3*24小时(3天)。主要为bangumi-id与详细信息的映射"`
		ThemoviedbCacheHour int `yaml:"themoviedb_cache_hour" json:"themoviedb_cache_hour" attr:"Themoviedb缓存时间" comment:"默认14*24小时(14天)。主要为tmdb-id与季度信息的映射"`
	} `yaml:"cache" json:"cache" attr:"缓存设置"`
}
