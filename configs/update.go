package configs

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/configs/version/v_110"
	"github.com/wetor/AnimeGo/configs/version/v_120"
	"github.com/wetor/AnimeGo/configs/version/v_130"
	"github.com/wetor/AnimeGo/configs/version/v_140"
	"github.com/wetor/AnimeGo/configs/version/v_141"
	"github.com/wetor/AnimeGo/configs/version/v_150"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
	encoder "github.com/wetor/AnimeGo/third_party/yaml-encoder"
)

type ConfigOnlyVersion struct {
	Version string `yaml:"version" json:"version"`
}

type Version struct {
	Name       string
	Desc       string
	UpdateFunc func(string) // 从上个版本升级到当前版本的升级函数
}

var (
	versions = []string{
		"1.1.0",
		"1.2.0",
		"1.3.0",
		"1.4.0",
		"1.4.1",
		"1.5.0",
		"1.5.1",
	}
	ConfigVersion = versions[len(versions)-1] // 当前配置文件版本

	versionList = []Version{
		{
			Name:       versions[0],
			UpdateFunc: func(s string) {},
		},
		{
			Name:       versions[1],
			Desc:       "插件配置结构变更；移除了自定义缓存文件、日志文件和临时文件功能",
			UpdateFunc: update_110_120,
		},
		{
			Name:       versions[2],
			Desc:       "插件配置结构变更；移除了js插件支持，增加了定时任务插件支持",
			UpdateFunc: update_120_130,
		},
		{
			Name:       versions[3],
			Desc:       "插件配置结构变更，支持设置参数",
			UpdateFunc: update_130_140,
		},
		{
			Name:       versions[4],
			Desc:       "新增重命名插件",
			UpdateFunc: update_140_141,
		},
		{
			Name:       versions[5],
			Desc:       "新增标题解析插件",
			UpdateFunc: update_141_150,
		},
		{
			Name:       versions[6],
			Desc:       "新增域名重定向设置",
			UpdateFunc: update_150_151,
		},
	}
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
	for i, v := range versions {
		if oldVer == v {
			oldIndex = i
			break
		}
	}
	if oldIndex < 0 {
		log.Fatal("配置文件升级失败：当前配置文件版本号错误 " + oldVer)
	}

	newVer := ConfigVersion
	newIndex := -1
	for i, v := range versions {
		if newVer == v {
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
		log.Printf("======= %s => %s =======\n", versions[i-1], versions[i])
		versionList[i].UpdateFunc(oldFile)
		if len(versionList[i].Desc) > 0 {
			log.Println("------------升级说明------------")
			log.Println(versionList[i].Desc)
		}
	}
	log.Println("===========子流程结束===========")
	log.Printf("配置文件升级完成：%s => %s\n", oldVer, newVer)
	log.Println("请确认配置后重新启动")
	return true
}

func update_110_120(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}
	oldConfig := &v_110.Config{}
	err = yaml.Unmarshal(data, oldConfig)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}

	newConfig := &v_120.Config{}
	err = copier.Copy(newConfig, oldConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	newConfig.Version = "1.2.0"
	log.Printf("[变动] 配置项(setting.filter.javascript) 变更为 setting.filter.plugin\n")
	if len(oldConfig.Filter.JavaScript) > 0 {
		newConfig.Filter.Plugin = make([]v_120.PluginInfo, len(oldConfig.Filter.JavaScript))
		for i, filename := range oldConfig.Filter.JavaScript {
			newConfig.Filter.Plugin[i] = v_120.PluginInfo{
				Enable: true,
				Type:   xpath.Ext(filename)[1:],
				File:   strings.TrimPrefix(filename, "plugin/"),
			}
		}
	}
	constant.Init(&constant.Options{
		DataPath: newConfig.DataPath,
	})
	log.Printf("[移除] 配置项(setting.advanced.xpath.db_file)\n")
	_ = utils.CreateMutiDir(constant.CachePath)
	_ = os.Rename(xpath.Join(oldConfig.DataPath, oldConfig.Advanced.Path.DbFile), constant.CacheFile)
	log.Printf("[移除] 配置项(setting.advanced.xpath.log_file)\n")
	_ = utils.CreateMutiDir(constant.LogPath)
	_ = os.Rename(xpath.Join(oldConfig.DataPath, oldConfig.Advanced.Path.LogFile), constant.LogFile)
	log.Printf("[移除] 配置项(setting.advanced.xpath.temp_path)\n")

	content, err := encodeConfig(newConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	err = os.WriteFile(file, content, 0644)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
}

func update_120_130(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}
	oldConfig := &v_120.Config{}
	err = yaml.Unmarshal(data, oldConfig)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}

	newConfig := &v_130.Config{}
	err = copier.Copy(newConfig, oldConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	newConfig.Version = "1.3.0"
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
	content, err := encodeConfig(newConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	err = os.WriteFile(file, content, 0644)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	// 强制写入
	assets.WritePlugins(assets.Dir, xpath.Join(newConfig.DataPath, assets.Dir), false)
}

func update_130_140(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}
	oldConfig := &v_130.Config{}
	err = yaml.Unmarshal(data, oldConfig)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}

	newConfig := &v_140.Config{}
	err = copier.Copy(newConfig, oldConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	newConfig.Version = "1.4.0"

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

	content, err := encodeConfig(newConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	err = os.WriteFile(file, content, 0644)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	// 强制写入
	assets.WritePlugins(assets.Dir, xpath.Join(newConfig.DataPath, assets.Dir), false)
}

func update_140_141(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}
	oldConfig := &v_140.Config{}
	err = yaml.Unmarshal(data, oldConfig)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}

	newConfig := &v_141.Config{}
	err = copier.Copy(newConfig, oldConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	newConfig.Version = "1.4.1"

	log.Println("[移除] 配置项(advanced.feed.multi_goroutine)")
	log.Println("[新增] 配置项(plugin.rename)")
	newConfig.Plugin.Rename = []v_141.PluginInfo{
		{
			Enable: true,
			Type:   "builtin",
			File:   "builtin_rename.py",
		},
	}
	content, err := encodeConfig(newConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	err = os.WriteFile(file, content, 0644)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	// 强制写入
	assets.WritePlugins(assets.Dir, xpath.Join(newConfig.DataPath, assets.Dir), false)
}

func update_141_150(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}
	oldConfig := &v_141.Config{}
	err = yaml.Unmarshal(data, oldConfig)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}

	newConfig := &v_150.Config{}
	err = copier.Copy(newConfig, oldConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	newConfig.Version = "1.5.0"

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
		fmt.Println(p)
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
	_ = os.Remove(xpath.Join(oldConfig.DataPath, "cache", "bolt.db"))
	content, err := encodeConfig(newConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	err = os.WriteFile(file, content, 0644)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	// 强制写入
	assets.WritePlugins(assets.Dir, xpath.Join(newConfig.DataPath, assets.Dir), false)
}

func update_150_151(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}
	oldConfig := &v_150.Config{}
	err = yaml.Unmarshal(data, oldConfig)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}

	newConfig := DefaultConfig()
	err = copier.Copy(newConfig, oldConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	newConfig.Version = "1.5.1"

	log.Println("[新增] 配置项(advanced.redirect.mikan)")
	log.Println("[新增] 配置项(advanced.redirect.bangumi)")
	log.Println("[新增] 配置项(advanced.redirect.themoviedb)")

	content, err := encodeConfig(newConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	err = os.WriteFile(file, content, 0644)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	// 强制写入
	assets.WritePlugins(assets.Dir, xpath.Join(newConfig.DataPath, assets.Dir), false)

}

func encodeConfig(conf any) ([]byte, error) {
	defaultSettingComment()
	defaultAdvancedComment()
	yml := encoder.NewEncoder(conf,
		encoder.WithComments(encoder.CommentsOnHead),
		encoder.WithCommentsMap(configComment),
	)
	content, err := yml.Encode()
	if err != nil {
		return nil, err
	}
	return content, nil
}

func BackupConfig(file string, version string) error {
	dir, name := xpath.Split(file)
	ext := xpath.Ext(name)
	name = strings.TrimSuffix(name, ext)
	timeStr := time.Now().Format("20060102150405")
	name = fmt.Sprintf("%s-%s-%s%s", name, version, timeStr, ext)
	oldFile, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	out := xpath.Join(dir, name)
	err = os.WriteFile(out, oldFile, 0644)
	if err != nil {
		return err
	}
	log.Printf("备份原配置文件到：'%s'\n", out)
	return nil
}
