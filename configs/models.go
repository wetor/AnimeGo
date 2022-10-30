package configs

type Config struct {
	Version   string `yaml:"version"`
	*Setting  `yaml:"setting"`
	*Advanced `yaml:"advanced"`
}
type Setting struct {
	Feed struct {
		Mikan struct {
			Name string
			Url  string
		}
	}
	Client struct {
		QBittorrent struct {
			Url      string
			Username string
			Password string
		}
	}
	Key struct {
		Themoviedb string
	}
	DataPath string `yaml:"data_path"`
	SavePath string `yaml:"save_path"`
	Category string `yaml:"category"` // 分类
	TagSrc   string `yaml:"tag"`      // 标签
	Proxy    struct {
		Enable bool   `yaml:"enable"`
		Url    string `json:"url"`
	} `yaml:"proxy"`
	Filter struct {
		JavaScript string `yaml:"javascript"` // 脚本名
	} `yaml:"filter"`
	WebApi struct {
		AccessKey string `yaml:"access_key"`
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
	} `yaml:"webapi"`
}

type Advanced struct {
	UpdateDelaySecond int `yaml:"update_delay_second"`

	Request struct {
		TimeoutSecond   int `yaml:"timeout_second"`
		RetryWaitSecond int `yaml:"retry_wait_second"`
		RetryNum        int `yaml:"retry_num"`
	} `yaml:"request"`

	Download struct {
		QueueDelaySecond       int  `yaml:"queue_delay_second"`
		QueueMaxNum            int  `yaml:"queue_max_num"`
		AllowDuplicateDownload bool `yaml:"allow_duplicate_download"`
		SeedingTime            int  `yaml:"seeding_time_minute"`
		IgnoreSizeMaxKb        int  `yaml:"ignore_size_max_kb"`
	} `yaml:"download"`

	Feed struct {
		UpdateDelayMinute int `yaml:"update_delay_minute"`
		Delay             int `yaml:"delay_second"`
		MultiGoroutine    struct {
			Enable       bool `yaml:"enable"`
			GoroutineMax int  `yaml:"goroutine_max"`
		} `yaml:"multi_goroutine"`
	} `yaml:"feed"`

	Path struct {
		DbFile   string `yaml:"db_file"`
		LogFile  string `yaml:"log_file"`
		TempPath string `yaml:"temp_path"`
	} `yaml:"path"`

	Default struct {
		TMDBFailSkip           bool `yaml:"tmdb_fail_skip"`
		TMDBFailUseTitleSeason bool `yaml:"tmdb_fail_use_title_season"`
		TMDBFailUseFirstSeason bool `yaml:"tmdb_fail_use_first_season"`
	} `yaml:"default"`

	Client struct {
		ConnectTimeoutSecond int `yaml:"connect_timeout_second"`
		RetryConnectNum      int `yaml:"retry_connect_num"`
		CheckTimeSecond      int `yaml:"check_time_second"`
	} `yaml:"client"`
}
