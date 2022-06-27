package config

type AdvancedConf struct {
	*GoBangumiConf  `yaml:"gobangumi"`
	*BangumiConf    `yaml:"bangumi"`
	*ThemoviedbConf `yaml:"themoviedb"`
	*MikanConf      `yaml:"mikan"`
}

type GoBangumiConf struct {
	RssDelay       int `yaml:"rss_delay_second"`
	MultiGoroutine struct {
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

func (a *AdvancedConf) GoBangumi() *GoBangumiConf {
	return a.GoBangumiConf
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
