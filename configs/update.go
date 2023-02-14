package configs

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"

	"github.com/wetor/AnimeGo/configs/version/v_110"
	"github.com/wetor/AnimeGo/configs/version/v_120"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
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
	ConfigVersion = "1.3.0" // 当前配置文件版本

	versions = []string{
		"1.1.0",
		"1.2.0",
		"1.3.0",
	}
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
		newConfig.Filter.Plugin = make([]models.Plugin, len(oldConfig.Filter.JavaScript))
		for i, filename := range oldConfig.Filter.JavaScript {
			newConfig.Filter.Plugin[i] = models.Plugin{
				Enable: true,
				Type:   "js",
				File:   strings.TrimPrefix(filename, "plugin/"),
			}
		}
	}
	constant.Init(&constant.Options{
		DataPath: newConfig.DataPath,
	})
	log.Printf("[移除] 配置项(setting.advanced.path.db_file)\n")
	_ = utils.CreateMutiDir(constant.CachePath)
	_ = os.Rename(path.Join(oldConfig.DataPath, oldConfig.Advanced.Path.DbFile), constant.CacheFile)
	log.Printf("[移除] 配置项(setting.advanced.path.log_file)\n")
	_ = utils.CreateMutiDir(constant.LogPath)
	_ = os.Rename(path.Join(oldConfig.DataPath, oldConfig.Advanced.Path.LogFile), constant.LogFile)
	log.Printf("[移除] 配置项(setting.advanced.path.temp_path)\n")

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

	newConfig := DefaultConfig()
	err = copier.Copy(newConfig, oldConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	newConfig.Version = "1.3.0"
	log.Printf("[变动] 配置项(setting.filter.javascript) 变更为 plugin.filter\n")
	if len(oldConfig.Filter.Plugin) > 0 {
		newConfig.Plugin.Filter = make([]models.Plugin, 0, len(oldConfig.Filter.Plugin))
		for _, p := range oldConfig.Filter.Plugin {
			if strings.ToLower(p.Type) == "js" || strings.ToLower(p.Type) == "javascript" {
				log.Printf("[移除] 不支持javascript插件 %s\n", p.File)
			} else {
				newConfig.Plugin.Filter = append(newConfig.Plugin.Filter, p)
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
	dir, name := filepath.Split(file)
	ext := path.Ext(name)
	name = strings.TrimSuffix(name, ext)
	timeStr := time.Now().Format("20060102150405")
	name = fmt.Sprintf("%s-%s-%s%s", name, version, timeStr, ext)
	oldFile, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	out := path.Join(dir, name)
	err = os.WriteFile(out, oldFile, 0644)
	if err != nil {
		return err
	}
	log.Printf("备份原配置文件到：'%s'\n", out)
	return nil
}
