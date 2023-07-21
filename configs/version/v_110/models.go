package v_110

type Config struct {
	Version  string `yaml:"version" json:"version"`
	Setting  `yaml:"setting" json:"setting"`
	Advanced `yaml:"advanced" json:"advanced"`
}

type Setting struct {
	Feed struct {
		Mikan struct {
			Name string `yaml:"name" json:"name"`
			Url  string `yaml:"url" json:"url"`
		} `yaml:"mikan" json:"mikan"`
	} `yaml:"feed" json:"feed"`
	Client struct {
		QBittorrent struct {
			Url      string `yaml:"url" json:"url"`
			Username string `yaml:"username" json:"username"`
			Password string `yaml:"password" json:"password"`
		} `yaml:"qbittorrent" json:"qbittorrent"`
	} `yaml:"client" json:"client"`
	DownloadPath string `yaml:"download_path" json:"download_path"`
	SavePath     string `yaml:"save_path" json:"save_path"`
	DataPath     string `yaml:"data_path" json:"data_path"`
	Filter       struct {
		JavaScript []string `yaml:"javascript" json:"javascript"`
	} `yaml:"filter" json:"filter"`
	Category string `yaml:"category" json:"category"`
	Tag      string `yaml:"tag" json:"tag"`
	WebApi   struct {
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
	UpdateDelaySecond int `yaml:"update_delay_second" json:"update_delay_second"`

	Request struct {
		TimeoutSecond   int `yaml:"timeout_second" json:"timeout_second"`
		RetryNum        int `yaml:"retry_num" json:"retry_num"`
		RetryWaitSecond int `yaml:"retry_wait_second" json:"retry_wait_second"`
	} `yaml:"request" json:"request"`

	Download struct {
		AllowDuplicateDownload bool   `yaml:"allow_duplicate_download" json:"allow_duplicate_download"`
		SeedingTimeMinute      int    `yaml:"seeding_time_minute" json:"seeding_time_minute"`
		IgnoreSizeMaxKb        int    `yaml:"ignore_size_max_kb" json:"ignore_size_max_kb"`
		Rename                 string `yaml:"rename" json:"rename"`
	} `yaml:"download" json:"download"`

	Feed struct {
		UpdateDelayMinute int `yaml:"update_delay_minute" json:"update_delay_minute"`
		DelaySecond       int `yaml:"delay_second" json:"delay_second"`
		MultiGoroutine    struct {
			Enable       bool `yaml:"enable" json:"enable"`
			GoroutineMax int  `yaml:"goroutine_max" json:"goroutine_max"`
		} `yaml:"multi_goroutine" json:"multi_goroutine"`
	} `yaml:"feed" json:"feed"`

	Path struct {
		DbFile   string `yaml:"db_file" json:"db_file"`
		LogFile  string `yaml:"log_file" json:"log_file"`
		TempPath string `yaml:"temp_path" json:"temp_path"`
	} `yaml:"path" json:"path"`

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
