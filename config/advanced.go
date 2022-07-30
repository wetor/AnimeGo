package config

type AdvancedConf struct {
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
	RssDelay                 int  `yaml:"rss_delay_second"`
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

func (a *AdvancedConf) Client() *ClientConf {
	return a.ClientConf
}
func (a *AdvancedConf) Main() *MainConf {
	return a.MainConf
}
func (a *AdvancedConf) Bangumi() *BangumiConf {
	return a.BangumiConf
}
func (a *AdvancedConf) Themoviedb() *ThemoviedbConf {
	return a.ThemoviedbConf
}
func (a *AdvancedConf) Mikan() *MikanConf {
	return a.MikanConf
}
