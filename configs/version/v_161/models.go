package v_161

type Config struct {
	Version  string `yaml:"version" json:"version"`
	Setting  `yaml:"setting" json:"setting"`
	Plugin   `yaml:"plugin" json:"plugin"`
	Advanced `yaml:"advanced" json:"advanced"`
}

type PluginInfo struct {
	Enable bool           `yaml:"enable" json:"enable"`
	Type   string         `yaml:"type" json:"type"`
	File   string         `yaml:"file" json:"file"`
	Args   map[string]any `yaml:"args,omitempty" json:"args,omitempty"`
	Vars   map[string]any `yaml:"vars,omitempty" json:"vars,omitempty"`
}

type Plugin struct {
	Feed     []PluginInfo `yaml:"feed" json:"feed"`
	Parser   []PluginInfo `yaml:"parser" json:"parser"`
	Filter   []PluginInfo `yaml:"filter" json:"filter"`
	Schedule []PluginInfo `yaml:"schedule" json:"schedule"`
	Rename   []PluginInfo `yaml:"rename" json:"rename"`
}

type Setting struct {
	Client struct {
		QBittorrent struct {
			Url          string `yaml:"url" json:"url"`
			Username     string `yaml:"username" json:"username"`
			Password     string `yaml:"password" json:"password"`
			DownloadPath string `yaml:"download_path" json:"download_path"`
		} `yaml:"qbittorrent" json:"qbittorrent"`
	} `yaml:"client" json:"client"`
	DownloadPath string `yaml:"download_path" json:"download_path"`
	SavePath     string `yaml:"save_path" json:"save_path"`
	DataPath     string `yaml:"data_path" json:"data_path"`
	Category     string `yaml:"category" json:"category"`
	Tag          string `yaml:"tag" json:"tag"`
	WebApi       struct {
		AccessKey string `yaml:"access_key" json:"access_key"`
		Host      string `yaml:"host" json:"host"`
		Port      int    `yaml:"port" json:"port"`
	} `yaml:"webapi" json:"webapi"`
	Proxy struct {
		Enable bool   `yaml:"enable" json:"enable"`
		Url    string `yaml:"url" json:"url"`
	} `yaml:"proxy" json:"proxy"`
	Key struct {
		Themoviedb string `yaml:"themoviedb" json:"themoviedb"`
	} `yaml:"key" json:"key"`
}

type Advanced struct {
	RefreshSecond int `yaml:"refresh_second" json:"refresh_second"`

	AniData struct {
		Mikan struct {
			Redirect string `yaml:"redirect" json:"redirect"`
			Cookie   string `yaml:"cookie" json:"cookie"`
		} `yaml:"mikan" json:"mikan"`
		Bangumi struct {
			Redirect string `yaml:"redirect" json:"redirect"`
		} `yaml:"bangumi" json:"bangumi"`
		Themoviedb struct {
			Redirect string `yaml:"redirect" json:"redirect"`
		} `yaml:"themoviedb" json:"themoviedb"`
	} `yaml:"anidata" json:"anidata"`

	Request struct {
		TimeoutSecond   int `yaml:"timeout_second" json:"timeout_second"`
		RetryNum        int `yaml:"retry_num" json:"retry_num"`
		RetryWaitSecond int `yaml:"retry_wait_second" json:"retry_wait_second"`
	} `yaml:"request" json:"request"`

	Download struct {
		AllowDuplicateDownload bool   `yaml:"allow_duplicate_download" json:"allow_duplicate_download"`
		SeedingTimeMinute      int    `yaml:"seeding_time_minute" json:"seeding_time_minute"`
		Rename                 string `yaml:"rename" json:"rename"`
	} `yaml:"download" json:"download"`

	Feed struct {
		DelaySecond int `yaml:"delay_second" json:"delay_second"`
	} `yaml:"feed" json:"feed"`

	Default struct {
		TMDBFailSkip           bool `yaml:"tmdb_fail_skip" json:"tmdb_fail_skip"`
		TMDBFailUseTitleSeason bool `yaml:"tmdb_fail_use_title_season" json:"tmdb_fail_use_title_season"`
		TMDBFailUseFirstSeason bool `yaml:"tmdb_fail_use_first_season" json:"tmdb_fail_use_first_season"`
	} `yaml:"default" json:"default"`

	Client struct {
		ConnectTimeoutSecond int `yaml:"connect_timeout_second" json:"connect_timeout_second"`
		RetryConnectNum      int `yaml:"retry_connect_num" json:"retry_connect_num"`
		CheckTimeSecond      int `yaml:"check_time_second" json:"check_time_second"`
	} `yaml:"client" json:"client"`

	Cache struct {
		MikanCacheHour      int `yaml:"mikan_cache_hour" json:"mikan_cache_hour"`
		BangumiCacheHour    int `yaml:"bangumi_cache_hour" json:"bangumi_cache_hour"`
		ThemoviedbCacheHour int `yaml:"themoviedb_cache_hour" json:"themoviedb_cache_hour"`
	} `yaml:"cache" json:"cache"`
}
