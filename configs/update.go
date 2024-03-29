package configs

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"

	"github.com/wetor/AnimeGo/configs/version/v_110"
	"github.com/wetor/AnimeGo/configs/version/v_120"
	"github.com/wetor/AnimeGo/configs/version/v_130"
	"github.com/wetor/AnimeGo/configs/version/v_140"
	"github.com/wetor/AnimeGo/configs/version/v_141"
	"github.com/wetor/AnimeGo/configs/version/v_150"
	"github.com/wetor/AnimeGo/configs/version/v_151"
	"github.com/wetor/AnimeGo/configs/version/v_152"
	"github.com/wetor/AnimeGo/configs/version/v_160"
	"github.com/wetor/AnimeGo/configs/version/v_161"
	"github.com/wetor/AnimeGo/configs/version/v_162"
	"github.com/wetor/AnimeGo/configs/version/v_170"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/dirdb"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

type ConfigOnlyVersion struct {
	Version string `yaml:"version" json:"version"`
}

type Version struct {
	Name       string
	Desc       string
	UpdateFunc func(string, string) // 从上个版本升级到当前版本的升级函数
}

var (
	versionList = []Version{
		{
			Name:       "1.1.0",
			UpdateFunc: func(f, v string) {},
		},
		{
			Name:       "1.2.0",
			Desc:       "插件配置结构变更；移除了自定义缓存文件、日志文件和临时文件功能",
			UpdateFunc: update(&v_110.Config{}, &v_120.Config{}, update_110_120),
		},
		{
			Name:       "1.3.0",
			Desc:       "插件配置结构变更；移除了js插件支持，增加了定时任务插件支持",
			UpdateFunc: update(&v_120.Config{}, &v_130.Config{}, update_120_130),
		},
		{
			Name:       "1.4.0",
			Desc:       "插件配置结构变更，支持设置参数",
			UpdateFunc: update(&v_130.Config{}, &v_140.Config{}, update_130_140),
		},
		{
			Name:       "1.4.1",
			Desc:       "新增重命名插件",
			UpdateFunc: update(&v_140.Config{}, &v_141.Config{}, update_140_141),
		},
		{
			Name:       "1.5.0",
			Desc:       "新增标题解析插件",
			UpdateFunc: update(&v_141.Config{}, &v_150.Config{}, update_141_150),
		},
		{
			Name:       "1.5.1",
			Desc:       "新增域名重定向设置",
			UpdateFunc: update(&v_150.Config{}, &v_151.Config{}, update_150_151),
		},
		{
			Name:       "1.5.2",
			Desc:       "新增下载器独立下载路径设置",
			UpdateFunc: update(&v_151.Config{}, &v_152.Config{}, update_151_152),
		},
		{
			Name:       "1.6.0",
			Desc:       "更改字段名，数据库迁移",
			UpdateFunc: update(&v_152.Config{}, &v_160.Config{}, update_152_160),
		},
		{
			Name:       "1.6.1",
			Desc:       "更改域名重定向设置，支持设置Mikan的Cookie",
			UpdateFunc: update(&v_160.Config{}, &v_161.Config{}, update_160_161),
		},
		{
			Name:       "1.6.2",
			Desc:       "新增Database设置",
			UpdateFunc: update(&v_161.Config{}, &v_162.Config{}, update_161_162),
		},
		{
			Name:       "1.7.0",
			Desc:       "更改下载器配置，新增Transmission客户端支持",
			UpdateFunc: update(&v_162.Config{}, &v_170.Config{}, update_162_170),
		},
		{
			Name:       "1.7.1",
			Desc:       "更改Themoviedb的ApiKey字段",
			UpdateFunc: update(&v_170.Config{}, DefaultConfig(), update_170_171),
		},
	}
	ConfigVersion = "1.7.1" // 当前配置文件版本
)

func UpdateConfig(oldFile string, backup bool) (restart bool) {
	// 载入旧版配置文件
	data, err := os.ReadFile(oldFile)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}
	configOnlyVersion := &ConfigOnlyVersion{}
	err = yaml.Unmarshal(data, configOnlyVersion)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}

	// 取出配置文件版本号和最新版本号
	oldVer := configOnlyVersion.Version
	// 版本号转换升级函数index
	oldIndex := -1
	for i, v := range versionList {
		if oldVer == v.Name {
			oldIndex = i
			break
		}
	}
	if oldIndex < 0 {
		log.Fatal("配置文件升级失败：当前配置文件版本号错误 " + oldVer)
	}

	newVer := ConfigVersion
	newIndex := -1
	for i, v := range versionList {
		if newVer == v.Name {
			newIndex = i
			break
		}
	}
	if newIndex < 0 {
		log.Fatal("配置文件升级失败：待升级版本号错误 " + newVer)
	}

	// 版本号相同，无需升级
	if oldIndex == newIndex {
		return false
	}
	log.Printf("配置文件升级：%s => %s\n", oldVer, newVer)
	if backup {
		err = BackupConfig(oldFile, oldVer)
		if err != nil {
			log.Fatal("配置文件备份失败：", err)
		}
	}

	log.Println("===========升级子流程===========")
	// 执行升级函数
	for i := oldIndex + 1; i <= newIndex; i++ {
		ver := versionList[i]
		log.Printf("======= %s => %s =======\n", versionList[i-1].Name, ver.Name)
		if len(ver.Desc) > 0 {
			log.Println("------------升级说明------------")
			log.Println(ver.Desc)
			log.Println("--------------------------------")
		}
		ver.UpdateFunc(oldFile, ver.Name)

	}
	log.Println("===========子流程结束===========")
	log.Printf("配置文件升级完成：%s => %s\n", oldVer, newVer)
	log.Println("请确认配置后重新启动")
	return true
}

func update(oldConfig, newConfig any, f func(any, any, string)) func(string, string) {
	return func(file, version string) {
		updateBefore(file, oldConfig, newConfig)
		f(oldConfig, newConfig, version)
		updateAfter(file, newConfig)
	}
}

func updateBefore(file string, oldConfig, newConfig any) {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}
	err = yaml.Unmarshal(data, oldConfig)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}

	err = copier.Copy(newConfig, oldConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
}

func updateAfter(file string, newConfig any) {
	content, err := encodeConfig(newConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	err = os.WriteFile(file, content, constant.WriteFilePerm)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
}

func update_110_120(old, new any, version string) {
	oldConfig := old.(*v_110.Config)
	newConfig := new.(*v_120.Config)

	newConfig.Version = version
	log.Printf("[变动] 配置项(setting.filter.javascript) 变更为 setting.filter.plugin\n")
	if len(oldConfig.Filter.JavaScript) > 0 {
		newConfig.Filter.Plugin = make([]v_120.PluginInfo, len(oldConfig.Filter.JavaScript))
		for i, filename := range oldConfig.Filter.JavaScript {
			newConfig.Filter.Plugin[i] = v_120.PluginInfo{
				Enable: true,
				Type:   path.Ext(filename)[1:],
				File:   strings.TrimPrefix(filename, "plugin/"),
			}
		}
	}
	constant.Init(&constant.Options{
		DataPath: newConfig.DataPath,
	})
	log.Printf("[移除] 配置项(setting.advanced.xpath.db_file)\n")
	_ = utils.CreateMutiDir(constant.CachePath)
	_ = os.Rename(path.Join(oldConfig.DataPath, oldConfig.Advanced.Path.DbFile), constant.CacheFile)
	log.Printf("[移除] 配置项(setting.advanced.xpath.log_file)\n")
	_ = utils.CreateMutiDir(constant.LogPath)
	_ = os.Rename(path.Join(oldConfig.DataPath, oldConfig.Advanced.Path.LogFile), constant.LogFile)
	log.Printf("[移除] 配置项(setting.advanced.xpath.temp_path)\n")
}

func update_120_130(old, new any, version string) {
	oldConfig := old.(*v_120.Config)
	newConfig := new.(*v_130.Config)

	newConfig.Version = version
	log.Printf("[变动] 配置项(setting.filter.javascript) 变更为 plugin.filter\n")
	if len(oldConfig.Filter.Plugin) > 0 {
		newConfig.Plugin.Filter = make([]v_130.PluginInfo, 0, len(oldConfig.Filter.Plugin))
		for _, p := range oldConfig.Filter.Plugin {
			if strings.ToLower(p.Type) == "js" || strings.ToLower(p.Type) == "javascript" {
				log.Printf("[移除] 不支持javascript插件 %s\n", p.File)
			} else {
				newConfig.Plugin.Filter = append(newConfig.Plugin.Filter, v_130.PluginInfo{
					Enable: p.Enable,
					Type:   p.Type,
					File:   p.File,
				})
			}
		}
	}
	// 强制写入
	assets.WritePlugins(assets.Dir, path.Join(newConfig.DataPath, assets.Dir), false)
}

func update_130_140(old, new any, version string) {
	oldConfig := old.(*v_130.Config)
	newConfig := new.(*v_140.Config)
	newConfig.Version = version

	log.Printf("[变动] 配置项(setting.feed.mikan) 变更为 plugin.feed 中的一个插件:\n")
	log.Printf("\t__name__: %s\n", oldConfig.Setting.Feed.Mikan.Name)
	log.Printf("\t__url__: %s\n", oldConfig.Setting.Feed.Mikan.Url)
	log.Printf("\t__cron__: %s\n", "0 0/20 * * * ?")
	log.Printf("\t默认关闭，每整20分钟执行\n")
	newConfig.Plugin.Feed = []v_140.PluginInfo{
		{
			Enable: false,
			Type:   "builtin",
			File:   "builtin_mikan_rss.py",
			Vars: map[string]any{
				"__name__": oldConfig.Setting.Feed.Mikan.Name,
				"__url__":  oldConfig.Setting.Feed.Mikan.Url,
				"__cron__": "0 0/20 * * * ?",
			},
		},
	}

	log.Printf("[新增] 配置项(plugin.filter) 支持 args\n")
	newConfig.Plugin.Filter = make([]v_140.PluginInfo, len(oldConfig.Plugin.Filter))
	_ = copier.Copy(&newConfig.Plugin.Filter, &oldConfig.Plugin.Filter)
	log.Printf("[新增] 配置项(plugin.schedule) 支持 args\n")
	newConfig.Plugin.Schedule = make([]v_140.PluginInfo, len(oldConfig.Plugin.Schedule))
	_ = copier.Copy(&newConfig.Plugin.Schedule, &oldConfig.Plugin.Schedule)

	// 强制写入
	assets.WritePlugins(assets.Dir, path.Join(newConfig.DataPath, assets.Dir), false)
}

func update_140_141(old, new any, version string) {
	// oldConfig := old.(*v_140.Config)
	newConfig := new.(*v_141.Config)
	newConfig.Version = version

	log.Println("[移除] 配置项(advanced.feed.multi_goroutine)")
	log.Println("[新增] 配置项(plugin.rename)")
	newConfig.Plugin.Rename = []v_141.PluginInfo{
		{
			Enable: true,
			Type:   "builtin",
			File:   "builtin_rename.py",
		},
	}
	// 强制写入
	assets.WritePlugins(assets.Dir, path.Join(newConfig.DataPath, assets.Dir), false)
}

func update_141_150(old, new any, version string) {
	oldConfig := old.(*v_141.Config)
	newConfig := new.(*v_150.Config)
	newConfig.Version = version

	log.Println("[移除] 配置项(advanced.download.ignore_size_max_kb)")
	log.Println("[新增] 配置项(plugin.parser)")
	newConfig.Plugin.Parser = []v_150.PluginInfo{
		{
			Enable: true,
			Type:   "builtin",
			File:   "builtin_parser.py",
		},
	}
	for i, p := range newConfig.Plugin.Feed {
		for key, val := range p.Vars {
			oldKey := key
			key = strings.TrimPrefix(key, "__")
			key = strings.TrimSuffix(key, "__")
			if key != oldKey {
				p.Vars[key] = val
				delete(p.Vars, oldKey)
			}
		}
		newConfig.Plugin.Feed[i] = p
	}
	log.Println("[清理] 清理缓存(data/cache/bolt.db)")
	_ = os.Remove(path.Join(oldConfig.DataPath, "cache", "bolt.db"))

	// 强制写入
	assets.WritePlugins(assets.Dir, path.Join(newConfig.DataPath, assets.Dir), false)
}

func update_150_151(old, new any, version string) {
	// oldConfig := old.(*v_150.Config)
	newConfig := new.(*v_151.Config)
	newConfig.Version = version

	log.Println("[新增] 配置项(advanced.redirect.mikan)")
	log.Println("[新增] 配置项(advanced.redirect.bangumi)")
	log.Println("[新增] 配置项(advanced.redirect.themoviedb)")

	// 强制写入
	assets.WritePlugins(assets.Dir, path.Join(newConfig.DataPath, assets.Dir), false)
}

func update_151_152(old, new any, version string) {
	oldConfig := old.(*v_151.Config)
	newConfig := new.(*v_152.Config)
	newConfig.Version = version

	log.Println("[新增] 配置项(setting.client.qbittorrent.download_path)")
	newConfig.Setting.Client.QBittorrent.DownloadPath = oldConfig.Setting.DownloadPath

	// 强制写入
	assets.WritePlugins(assets.Dir, path.Join(newConfig.DataPath, assets.Dir), false)
}

func update_152_160(old, new any, version string) {
	oldConfig := old.(*v_152.Config)
	newConfig := new.(*v_160.Config)
	newConfig.Version = version
	constant.Init(&constant.Options{
		DataPath: newConfig.DataPath,
	})
	log.Println("[变动] 配置项(advanced.update_delay_second) 变更为 advanced.refresh_second")
	newConfig.Advanced.RefreshSecond = oldConfig.Advanced.UpdateDelaySecond

	// 强制写入
	assets.WritePlugins(assets.Dir, path.Join(newConfig.DataPath, assets.Dir), false)

	log.Println("--------------------------------")
	log.Println("[数据迁移] 数据库部分迁移到文件标记")

	bolt2dirdb(constant.CacheFile, xpath.P(newConfig.Setting.SavePath))
	log.Println("--------------------------------")
}

func bolt2dirdb(boltPath, savePath string) {
	if !utils.IsExist(boltPath) {
		return
	}
	bolt := cache.NewBolt()
	bolt.Open(boltPath)
	defer bolt.Close()

	keys := bolt.ListKey("name2status")
	for _, key := range keys {

		status := models.DownloadStatus{}
		_ = bolt.Get("name2status", key, &status)
		entity := models.AnimeEntity{}
		_ = bolt.Get("name2entity", key, &entity)

		base := models.BaseDBEntity{
			Hash:     status.Hash,
			Name:     key,
			CreateAt: utils.Unix(),
			UpdateAt: utils.Unix(),
		}
		if len(status.Path) > 0 {
			animePath := xpath.Root(status.Path[0])
			_ = write(path.Join(savePath, animePath, constant.DatabaseAnimeDBName), models.AnimeDBEntity{
				BaseDBEntity: base,
			})
		}
		for i, f := range status.Path {
			file := path.Join(savePath, f)
			if utils.IsExist(file) {
				filename := fmt.Sprintf(constant.DatabaseEpisodeDBFmt, strings.TrimSuffix(f, path.Ext(f)))
				_ = write(path.Join(savePath, filename), models.EpisodeDBEntity{
					BaseDBEntity: base,
					StateDB: models.StateDB{
						Seeded:     true,
						Downloaded: true,
						Renamed:    true,
						Scraped:    true,
					},
					Season: entity.Season,
					Ep:     entity.Ep[i].Ep,
					Type:   entity.Ep[i].Type,
				})
			}
		}
	}
}

func write(file string, data any) error {
	f := dirdb.NewFile(file)
	err := f.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	err = f.DB.Marshal(data)
	if err != nil {
		return err
	}
	log.Printf("write %s: %+v\n", file, data)
	return nil
}

func update_160_161(old, new any, version string) {
	oldConfig := old.(*v_160.Config)
	newConfig := new.(*v_161.Config)
	newConfig.Version = version
	log.Println("[新增] 配置项(advanced.anidata.mikan.cookie)")
	log.Println("[变动] 配置项(advanced.redirect.mikan) 变更为 advanced.anidata.mikan.redirect")
	newConfig.Advanced.AniData.Mikan.Redirect = oldConfig.Advanced.Redirect.Mikan
	log.Println("[变动] 配置项(advanced.redirect.bangumi) 变更为 advanced.anidata.bangumi.redirect")
	newConfig.Advanced.AniData.Bangumi.Redirect = oldConfig.Advanced.Redirect.Bangumi
	log.Println("[变动] 配置项(advanced.redirect.themoviedb) 变更为 advanced.anidata.themoviedb.redirect")
	newConfig.Advanced.AniData.Themoviedb.Redirect = oldConfig.Advanced.Redirect.Themoviedb

	// 强制写入
	assets.WritePlugins(assets.Dir, path.Join(newConfig.DataPath, assets.Dir), false)
}

func update_161_162(old, new any, version string) {
	// oldConfig := old.(*v_161.Config)
	newConfig := new.(*v_162.Config)
	newConfig.Version = version

	log.Println("[新增] 配置项(advanced.database.refresh_database_cron)")
	newConfig.Advanced.Database.RefreshDatabaseCron = "0 0 6 * * *"

	// 强制写入
	assets.WritePlugins(assets.Dir, path.Join(newConfig.DataPath, assets.Dir), false)
}

func update_162_170(old, new any, version string) {
	oldConfig := old.(*v_162.Config)
	newConfig := new.(*v_170.Config)
	newConfig.Version = version

	log.Println("[新增] 配置项(setting.client.client)")
	log.Println("[变动] 配置项(setting.client.qbittorrent) 变更为 setting.client")
	newConfig.Setting.Client.Client = "QBittorrent"
	newConfig.Setting.Client.Username = oldConfig.Setting.Client.QBittorrent.Username
	newConfig.Setting.Client.Password = oldConfig.Setting.Client.QBittorrent.Password
	newConfig.Setting.Client.Url = oldConfig.Setting.Client.QBittorrent.Url
	log.Println("[变动] 配置项(advanced.download.seeding_time_minute) 变更为 advanced.client.seeding_time_minute")
	newConfig.Advanced.Client.SeedingTimeMinute = oldConfig.Advanced.Download.SeedingTimeMinute

	// 强制写入
	assets.WritePlugins(assets.Dir, path.Join(newConfig.DataPath, assets.Dir), false)
}

func update_170_171(old, new any, version string) {
	oldConfig := old.(*v_170.Config)
	newConfig := new.(*Config)
	newConfig.Version = version

	log.Println("[变动] 配置项(advanced.anidata) 变更为 advanced.source")
	newConfig.Advanced.Source.Mikan.Redirect = oldConfig.Advanced.AniData.Mikan.Redirect
	newConfig.Advanced.Source.Mikan.Cookie = oldConfig.Advanced.AniData.Mikan.Cookie
	newConfig.Advanced.Source.Bangumi.Redirect = oldConfig.Advanced.AniData.Bangumi.Redirect
	newConfig.Advanced.Source.Themoviedb.Redirect = oldConfig.Advanced.AniData.Themoviedb.Redirect
	log.Println("[变动] 配置项(setting.key.themoviedb) 变更为 advanced.source.themoviedb.api_key")
	newConfig.Advanced.Source.Themoviedb.ApiKey = oldConfig.Setting.Key.Themoviedb

	// 强制写入
	assets.WritePlugins(assets.Dir, path.Join(newConfig.DataPath, assets.Dir), false)
}
