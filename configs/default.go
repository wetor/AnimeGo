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

	defaultConfig.Setting.Client.QBittorrent.Url = "http://localhost:8080"
	defaultConfig.Setting.Client.QBittorrent.Username = "admin"
	defaultConfig.Setting.Client.QBittorrent.Password = "adminadmin"
	defaultConfig.Setting.Client.QBittorrent.DownloadPath = "./download/incomplete"

	defaultConfig.Setting.DataPath = "./data"
	defaultConfig.Setting.SavePath = "./download/anime"
	defaultConfig.Setting.DownloadPath = defaultConfig.Setting.Client.QBittorrent.DownloadPath

	defaultConfig.Setting.Category = "AnimeGo"
	defaultConfig.Setting.Tag = "{year}年{quarter}月新番"

	defaultConfig.Setting.WebApi.Host = "0.0.0.0"
	defaultConfig.Setting.WebApi.Port = 7991
	defaultConfig.Setting.WebApi.AccessKey = "animego123"

	defaultConfig.Setting.Proxy.Enable = false
	defaultConfig.Setting.Proxy.Url = "http://127.0.0.1:7890"

	defaultConfig.Setting.Key.Themoviedb = "d3d8430aefee6c19520d0f7da145daf5"
}

func defaultPluginComment() {
	configComment["plugin_help"] = `按顺序依次执行启用的插件
列表类型，每一项需要有以下参数：
  enable: 启用
  type: 插件类型，目前仅支持 'python'(py) 和 'builtin' 插件类型。builtin为内置插件
  file: 插件文件，相对于 'data/plugin' 文件夹的路径，或内置插件名
  args: [可空]插件额外参数，字典类型，会覆盖同名参数
  vars: [可空]插件全局变量，字典类型，如果变量名前缀或后缀不是'__'将会自动补充，即在插件中变量名前后缀始终为'__'，
    会覆盖插件脚本中同名变量，具体变量和作用参考订阅插件文档`
}

func defaultPlugin() {
	defaultConfig.Plugin.Feed = []PluginInfo{
		{
			Enable: false,
			Type:   "builtin",
			File:   "builtin_mikan_rss.py",
			Vars: map[string]any{
				"url":  "",
				"cron": "0 0/20 * * * ?",
			},
		},
	}
	defaultConfig.Plugin.Filter = []PluginInfo{
		{
			Enable: true,
			Type:   "py",
			File:   "filter/default.py",
		},
	}
	defaultConfig.Plugin.Rename = []PluginInfo{
		{
			Enable: true,
			Type:   "builtin",
			File:   "builtin_rename.py",
		},
	}
	defaultConfig.Plugin.Parser = []PluginInfo{
		{
			Enable: true,
			Type:   "builtin",
			File:   "builtin_parser.py",
		},
	}
}

func defaultAdvancedComment() {
	configComment["update_delay_second_help"] = `每隔这一段时间，都会将下载队列添加到下载项、更新下载进度和重命名已完成下载项
等待过程是异步的，等待期间不影响操作
在下载项较多、等待时间过少时会出现请求超时，所以有个最小等待时间为2秒的限制
默认为10，最小值为2`

	configComment["rename_help"] = `下载状态顺序为: 创建下载项->下载->下载完成->做种->做种完成
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
	defaultConfig.Advanced.Download.Rename = "wait_move"

	defaultConfig.Advanced.Feed.DelaySecond = 5

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
		defaultConfig.Version = ConfigVersion
		defaultSettingComment()
		defaultSetting()
		defaultPluginComment()
		defaultPlugin()
		defaultAdvancedComment()
		defaultAdvanced()
		isInit = true
	}
}

func DefaultConfig() *Config {
	defaultAll()
	return defaultConfig
}

func Config2Bytes(config *Config) ([]byte, error) {
	defaultAll()
	yaml := encoder.NewEncoder(config,
		encoder.WithComments(encoder.CommentsOnHead),
		encoder.WithCommentsMap(configComment),
	)
	content, err := yaml.Encode()
	if err != nil {
		return nil, err
	}
	return content, nil
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
	data, err := Config2Bytes(defaultConfig)
	if err != nil {
		return err
	}
	// 所有者可读可写，其他用户只读
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
