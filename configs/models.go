package configs

type Config struct {
	Feed struct {
		Rss map[string]*Rss
	}

	Client    map[string]*Client
	Key       map[string]string
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
	*ClientConf `yaml:"client"`
	*MainConf   `yaml:"main"`
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
	HttpTimeoutSecond int `yaml:"http_timeout_second"`
	HttpRetryNum      int `yaml:"http_retry_num"`
}

type Setting struct {
	DataPath  string `yaml:"data_path"`
	DbFile    string `yaml:"db_file"`
	SavePath  string `yaml:"save_path"`
	Category  string `yaml:"category"` // 分类
	TagSrc    string `yaml:"tag"`      // 标签
	ProxyConf struct {
		Enable bool
		Url    string
	} `yaml:"proxy"`
	*Filter `yaml:"filter"`
}

type Filter struct {
	JavaScript string `yaml:"javascript"` // 脚本名
}
