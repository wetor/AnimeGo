package configs

type Config struct {
	Feed struct {
		Rss map[string]*Rss
	}

	Client    map[string]*Client
	Key       map[string]string
	ProxyConf struct {
		Enable bool
		Url    string
	} `yaml:"proxy"`
	*Setting  `yaml:"setting"`
	*Advanced `yaml:"advanced"`
}
type Client struct {
	Url      string
	Username string
	Password string
}
type Rss struct {
	Name string
	Url  string
}

type Advanced struct {
	*ClientConf     `yaml:"client"`
	*MainConf       `yaml:"main"`
	*BangumiConf    `yaml:"bangumi"`
	*ThemoviedbConf `yaml:"themoviedb"`
	*MikanConf      `yaml:"mikan"`
}

type ClientConf struct {
	ConnectTimeoutSecond int `yaml:"connect_timeout_second"`
	RetryConnectNum      int `yaml:"retry_connect_num"`
}

type MainConf struct {
	UpdateDelaySecond        int  `yaml:"update_delay_second"`
	DownloadQueueDelaySecond int  `yaml:"download_queue_delay_second"`
	DownloadQueueMaxNum      int  `yaml:"download_queue_max_num"`
	AllowDuplicateDownload   bool `yaml:"allow_duplicate_download"`
	SeedingTime              int  `yaml:"seeding_time_minute"`
	IgnoreSizeMaxKb          int  `yaml:"ignore_size_max_kb"`
	FeedUpdateDelayMinute    int  `yaml:"feed_update_delay_minute"`
	FeedDelay                int  `yaml:"feed_delay_second"`
	MultiGoroutine           struct {
		Enable       bool `yaml:"enable"`
		GoroutineMax int  `yaml:"goroutine_max"`
	} `yaml:"multi_goroutine"`
}

type BangumiConf struct {
	Host            string `yaml:"host"`
	MatchEpRange    int    `yaml:"match_ep_range"`
	MatchEpDays     int    `yaml:"match_ep_days"`
	CacheInfoExpire int64  `yaml:"cache_info_expire_second"`
	CacheEpExpire   int64  `yaml:"cache_ep_expire_second"`
}

type ThemoviedbConf struct {
	Host              string `yaml:"host"`
	MatchSeasonDays   int    `yaml:"match_season_days"`
	CacheIdExpire     int64  `yaml:"cache_id_expire_second"`
	CacheSeasonExpire int64  `yaml:"cache_season_expire_second"`
}

type MikanConf struct {
	Host               string `yaml:"host"`
	CacheIdExpire      int64  `yaml:"cache_id_expire_second"`
	CacheBangumiExpire int64  `yaml:"cache_bangumi_expire_second"`
}

type Setting struct {
	DataPath  string `yaml:"data_path"`
	CachePath string `yaml:"cache_path"`
	SavePath  string `yaml:"save_path"`
	Category  string `yaml:"category"` // 分类
	TagSrc    string `yaml:"tag"`      // 标签
}
