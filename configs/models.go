package configs

type Config struct {
	Version  string `yaml:"version" json:"version" attr:"配置文件版本"`
	Setting  `yaml:"setting" json:"setting" attr:"常规设置"`
	Plugin   `yaml:"plugin" json:"plugin" attr:"插件设置" comment_key:"plugin_help"`
	Advanced `yaml:"advanced" json:"advanced" attr:"高级设置"`
}

type PluginInfo struct {
	Enable bool           `yaml:"enable" json:"enable"`
	Type   string         `yaml:"type" json:"type"`
	File   string         `yaml:"file" json:"file"`
	Args   map[string]any `yaml:"args,omitempty" json:"args,omitempty"`
	Vars   map[string]any `yaml:"vars,omitempty" json:"vars,omitempty"`
}

type Plugin struct {
	Feed     []PluginInfo `yaml:"feed" json:"feed" attr:"订阅" comment:"解析订阅链接"`
	Parser   []PluginInfo `yaml:"parser" json:"parser" attr:"解析器" comment:"解析项目标题"`
	Filter   []PluginInfo `yaml:"filter" json:"filter" attr:"过滤器插件" comment:"用来筛选符合条件的项目进行解析下载"`
	Schedule []PluginInfo `yaml:"schedule" json:"schedule" attr:"定时任务" comment:"定时执行脚本"`
	Rename   []PluginInfo `yaml:"rename" json:"rename" attr:"重命名" comment:"下载完成后重命名规则"`
}

type Setting struct {
	Client struct {
		QBittorrent struct {
			Url          string `yaml:"url" json:"url" attr:"地址" comment:"环境变量ANIMEGO_QBT_URL"`
			Username     string `yaml:"username" json:"username" attr:"用户名" comment:"环境变量ANIMEGO_QBT_USERNAME"`
			Password     string `yaml:"password" json:"password" attr:"密码" comment:"环境变量ANIMEGO_QBT_PASSWORD"`
			DownloadPath string `yaml:"download_path" json:"download_path" attr:"下载文件夹" comment:"环境变量ANIMEGO_QBT_DOWNLOAD_PATH"`
		} `yaml:"qbittorrent" json:"qbittorrent" attr:"qBittorrent客户端"`
	} `yaml:"client" json:"client" attr:"下载客户端设置"`
	DownloadPath string `yaml:"download_path" json:"download_path" attr:"下载文件夹" comment:"环境变量ANIMEGO_DOWNLOAD_PATH. 下载器的下载文件夹"`
	SavePath     string `yaml:"save_path" json:"save_path" attr:"保存文件夹" comment:"环境变量ANIMEGO_SAVE_PATH. 下载完成后，重命名并移动到的文件夹"`
	DataPath     string `yaml:"data_path" json:"data_path" attr:"数据文件夹" comment:"环境变量ANIMEGO_DATA_PATH. 用于保存数据库、插件等数据"`
	Category     string `yaml:"category" json:"category" attr:"分类名" comment:"环境变量ANIMEGO_CATEGORY. 仅qBittorrent有效"`
	Tag          string `yaml:"tag" json:"tag" attr:"标签表达式" comment:"环境变量ANIMEGO_TAG" comment_key:"tag_help"`
	WebApi       struct {
		AccessKey string `yaml:"access_key" json:"access_key" attr:"请求秘钥" comment:"环境变量ANIMEGO_WEB_ACCESS_KEY. 为空则不需要验证"`
		Host      string `yaml:"host" json:"host" attr:"域名" comment:"环境变量ANIMEGO_WEB_HOST"`
		Port      int    `yaml:"port" json:"port" attr:"端口" comment:"环境变量ANIMEGO_WEB_PORT"`
	} `yaml:"webapi" json:"webapi" attr:"WebApi设置"`
	Proxy struct {
		Enable bool   `yaml:"enable" json:"enable" attr:"启用" comment:"环境变量ANIMEGO_PROXY_URL不为空则启用，否则禁用"`
		Url    string `yaml:"url" json:"url" attr:"代理链接" comment:"环境变量ANIMEGO_PROXY_URL"`
	} `yaml:"proxy" json:"proxy" attr:"代理设置" comment:"开启后AnimeGo所有的网络请求都会使用代理"`
	Key struct {
		Themoviedb string `yaml:"themoviedb" json:"themoviedb" attr:"TheMovieDB的APIkey" comment:"环境变量ANIMEGO_THEMOVIEDB_KEY" comment_key:"themoviedb_key"`
	} `yaml:"key" json:"key" attr:"秘钥设置"`
}

type Advanced struct {
	RefreshSecond int `yaml:"refresh_second" json:"refresh_second" attr:"刷新间隔时间" comment_key:"refresh_second_help"`

	AniData struct {
		Mikan struct {
			Redirect string `yaml:"redirect" json:"redirect" attr:"默认mikanani.me"`
			Cookie   string `yaml:"cookie" json:"cookie" attr:"mikan的Cookie" comment:"使用登录后的Cookie可以正常下载mikan的被隐藏番剧. 登录状态的Cookie名为'.AspNetCore.Identity.Application'"`
		} `yaml:"mikan" json:"mikan"`
		Bangumi struct {
			Redirect string `yaml:"redirect" json:"redirect" attr:"默认api.bgm.tv"`
		} `yaml:"bangumi" json:"bangumi"`
		Themoviedb struct {
			Redirect string `yaml:"redirect" json:"redirect" attr:"默认api.themoviedb.org"`
		} `yaml:"themoviedb" json:"themoviedb"`
	} `yaml:"anidata" json:"anidata" attr:"资源网站设置"`

	Request struct {
		TimeoutSecond   int `yaml:"timeout_second" json:"timeout_second" attr:"请求超时时间"`
		RetryNum        int `yaml:"retry_num" json:"retry_num" attr:"额外重试次数"`
		RetryWaitSecond int `yaml:"retry_wait_second" json:"retry_wait_second" attr:"重试间隔等待时间"`
	} `yaml:"request" json:"request" attr:"网络请求设置"`

	Download struct {
		AllowDuplicateDownload bool   `yaml:"allow_duplicate_download" json:"allow_duplicate_download" attr:"允许重复下载"`
		SeedingTimeMinute      int    `yaml:"seeding_time_minute" json:"seeding_time_minute" attr:"做种时间"`
		Rename                 string `yaml:"rename" json:"rename" attr:"重命名方式" comment_key:"rename_help"`
	} `yaml:"download" json:"download" attr:"下载设置"`

	Feed struct {
		DelaySecond int `yaml:"delay_second" json:"delay_second" attr:"订阅解析间隔时间"`
	} `yaml:"feed" json:"feed" attr:"订阅设置"`

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

type Environment struct {
	QbtUrl          *string `env:"QBT_URL" val:"Setting.Client.QBittorrent.Url"`
	QbtUsername     *string `env:"QBT_USERNAME" val:"Setting.Client.QBittorrent.Username"`
	QbtPassword     *string `env:"QBT_PASSWORD" val:"Setting.Client.QBittorrent.Password"`
	QbtDownloadPath *string `env:"QBT_DOWNLOAD_PATH" val:"Setting.Client.QBittorrent.DownloadPath"`

	DownloadPath *string `env:"DOWNLOAD_PATH" val:"Setting.DownloadPath"`
	SavePath     *string `env:"SAVE_PATH" val:"Setting.SavePath"`
	DataPath     *string `env:"DATA_PATH" val:"Setting.DataPath"`
	Category     *string `env:"CATEGORY" val:"Setting.Category"`
	Tag          *string `env:"TAG" val:"Setting.Tag"`

	WebAccessKey *string `env:"WEB_ACCESS_KEY" val:"Setting.WebApi.AccessKey"`
	WebHost      *string `env:"WEB_HOST" val:"Setting.WebApi.Host"`
	WebPort      *int    `env:"WEB_PORT" val:"Setting.WebApi.Port"`

	ProxyUrl *string `env:"PROXY_URL" val:"Setting.Proxy.Url"`

	ThemoviedbKey *string `env:"THEMOVIEDB_KEY" val:"Setting.Key.Themoviedb"`
}
