package configs

import (
	"os"

	encoder "github.com/wetor/AnimeGo/third_party/yaml-encoder"
)

var (
	defaultConfig = &Config{}
	configComment = make(map[string]string)
	isInit        = false
)

func defaultSettingComment() {
	configComment["filter_javascript"] = `js插件的文件名列表，依次执行。路径相对于data_path
插件名可以忽略'.js'后缀；插件名也可以使用上层文件夹名，会自动加载文件夹内部的 'main.js' 或 'plugin.js'
如设置为 'plugin/test'，会依次尝试加载 'plugin/test/main.js', 'plugin/test/plugin.js', 'plugin/test.js'
`

	configComment["tag_help"] = `仅qBittorrent有效，可用通配符列表：
  {year} int 番剧更新年
  {quarter} int 番剧季度月号，取值为[4, 7, 10, 1]分别对应[春, 夏, 秋, 冬]季番剧
  {quarter_index} int 番剧季度序号，取值为[1, 2, 3, 4]分别对应春(4月)、夏(7月)、秋(10月)、冬(1月)季番剧
  {quarter_name} string 番剧季度名，取值为[春, 夏, 秋, 冬]
  {ep} int 番剧当前剧集序号，从1开始
  {week} int 番剧更新星期数，取值为[1, 2, 3, 4, 5, 6, 7]
  {week_name} string 番剧更新星期名，取值为[星期一, 星期二, 星期三, 星期四, 星期五, 星期六, 星期日]`

	configComment["themoviedb_key"] = `可以自行申请链接（需注册）：https://www.themoviedb.org/settings/api?language=zh-CN
以下为wetor的个人APIkey，仅用于AnimeGo使用`
}

func defaultSetting() {
	defaultConfig.Setting.Feed.Mikan.Name = "Mikan"
	defaultConfig.Setting.Feed.Mikan.Url = ""

	defaultConfig.Setting.Client.QBittorrent.Url = "http://localhost:8080"
	defaultConfig.Setting.Client.QBittorrent.Username = "admin"
	defaultConfig.Setting.Client.QBittorrent.Password = "adminadmin"

	defaultConfig.Setting.DataPath = "./data"
	defaultConfig.Setting.SavePath = "./download/anime"
	defaultConfig.Setting.DownloadPath = "./download/incomplete"

	defaultConfig.Setting.Filter.JavaScript = []string{
		"plugin/filter/default.js",
	}

	defaultConfig.Setting.Category = "AnimeGo"
	defaultConfig.Setting.Tag = "{year}年{quarter}月新番"

	defaultConfig.Setting.WebApi.Host = "localhost"
	defaultConfig.Setting.WebApi.Port = 7991
	defaultConfig.Setting.WebApi.AccessKey = "animego123"

	defaultConfig.Setting.Proxy.Enable = false
	defaultConfig.Setting.Proxy.Url = "http://127.0.0.1:7890"

	defaultConfig.Setting.Key.Themoviedb = "d3d8430aefee6c19520d0f7da145daf5"
}

func defaultAdvancedComment() {
	configComment["update_delay_second"] = `每隔这一段时间，都会更新一次下载进度、清空下载队列(添加到下载项)
等待过程是异步的，等待期间不影响操作
在下载项较多、等待时间过少时会出现请求超时，所以有个最小等待时间为2秒的限制
默认为10，最小值为2`

	configComment["rename"] = `下载状态顺序为: 创建下载项->下载->下载完成->做种->做种完成
可选值为: ['link', 'link_delete', 'move', 'wait_move']
  link: 使用硬链接方式，下载完成后触发。不影响做种
  link_delete: 使用硬链接方式，下载完成后触发。不影响做种，做种完成后删除原文件
  move: 使用移动方式，下载完成后触发。无法做种
  wait_move: 使用移动方式，做种完成后触发`
}

func defaultAdvanced() {
	defaultConfig.Advanced.UpdateDelaySecond = 10

	defaultConfig.Advanced.Request.TimeoutSecond = 5
	defaultConfig.Advanced.Request.RetryNum = 3
	defaultConfig.Advanced.Request.RetryWaitSecond = 5

	defaultConfig.Advanced.Download.AllowDuplicateDownload = false
	defaultConfig.Advanced.Download.SeedingTimeMinute = 0
	defaultConfig.Advanced.Download.IgnoreSizeMaxKb = 1024
	defaultConfig.Advanced.Download.Rename = "wait_move"

	defaultConfig.Advanced.Feed.UpdateDelayMinute = 20
	defaultConfig.Advanced.Feed.DelaySecond = 5
	defaultConfig.Advanced.Feed.MultiGoroutine.Enable = false
	defaultConfig.Advanced.Feed.MultiGoroutine.GoroutineMax = 4

	defaultConfig.Advanced.Path.DbFile = "cache/bolt.db"
	defaultConfig.Advanced.Path.LogFile = "log/animego.log"
	defaultConfig.Advanced.Path.TempPath = "temp"

	defaultConfig.Advanced.Default.TMDBFailSkip = false
	defaultConfig.Advanced.Default.TMDBFailUseTitleSeason = true
	defaultConfig.Advanced.Default.TMDBFailUseFirstSeason = true

	defaultConfig.Advanced.Client.ConnectTimeoutSecond = 5
	defaultConfig.Advanced.Client.RetryConnectNum = 10
	defaultConfig.Advanced.Client.CheckTimeSecond = 30

	defaultConfig.Advanced.Cache.MikanCacheHour = 7 * 24
	defaultConfig.Advanced.Cache.BangumiCacheHour = 3 * 24
	defaultConfig.Advanced.Cache.ThemoviedbCacheHour = 14 * 24
}

func defaultAll() {
	if !isInit {
		defaultConfig.Version = os.Getenv("ANIMEGO_CONFIG_VERSION")
		defaultSettingComment()
		defaultSetting()
		defaultAdvancedComment()
		defaultAdvanced()
		isInit = true
	}
}

func DefaultConfig() *Config {
	defaultAll()
	return defaultConfig
}

func Config2Bytes(config *Config) []byte {
	defaultAll()
	yaml := encoder.NewEncoder(config,
		encoder.WithComments(encoder.CommentsOnHead),
		encoder.WithCommentsMap(configComment),
	)
	content, err := yaml.Encode()
	if err != nil {
		panic(err)
	}
	return content
}

func DefaultDoc() []byte {
	defaultAll()
	yaml := encoder.NewEncoder(defaultConfig,
		encoder.WithComments(encoder.CommentsOnHead),
		encoder.WithCommentsMap(configComment),
	)
	content, err := yaml.EncodeDoc()
	if err != nil {
		panic(err)
	}
	return content
}

func DefaultFile(filename string) error {
	// 所有者可读可写，其他用户只读
	err := os.WriteFile(filename, Config2Bytes(defaultConfig), 0644)
	if err != nil {
		return err
	}
	return nil
}
