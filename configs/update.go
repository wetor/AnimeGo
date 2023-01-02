package configs

import (
	"fmt"
	encoder "github.com/wetor/yaml-encoder"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

var (
	versionMap = map[string]int{
		"1.0.0": 0,
		"1.1.0": 1,
	}
	updateFunc = []func(string){
		update_100_110,
	}
	updateConfig *Config
)

func UpdateConfig(oldFile string, backup bool) (restart bool) {
	// 载入配置文件
	data, err := os.ReadFile(oldFile)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}
	updateConfig = &Config{}
	err = yaml.Unmarshal(data, updateConfig)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}

	// 取出配置文件版本号和最新版本号
	oldVer := updateConfig.Version
	newVer := os.Getenv("ANIMEGO_CONFIG_VERSION")
	// 版本号转换升级函数index
	oldIndex, ok := versionMap[oldVer]
	if !ok {
		log.Fatal("配置文件升级失败：当前配置文件版本号错误 " + oldVer)
	}
	newIndex, ok := versionMap[newVer]
	if !ok {
		log.Fatal("配置文件升级失败：待升级版本号错误 " + newVer)
	}

	// 版本号相同，无需升级
	if oldIndex == newIndex {
		return false
	}
	log.Printf("配置文件升级：%s => %s\n", oldVer, newVer)
	if backup {
		err = backupConfig(oldFile)
		if err != nil {
			log.Fatal("配置文件备份失败：", err)
		}
	}

	log.Println("===========升级子流程===========")
	// 执行升级函数
	for i := oldIndex + 1; i <= newIndex; i++ {
		updateFunc[i-1](oldFile)
	}
	log.Println("===========子流程结束===========")
	log.Printf("配置文件升级完成：%s => %s\n", oldVer, newVer)
	log.Println("请确认配置后重新启动")
	return true
}

func update_100_110(file string) {
	log.Println("======= 1.0.0 => 1.1.0 =======")
	updateConfig.Version = "1.1.0"
	tmp := updateConfig.Setting.SavePath
	updateConfig.Setting.SavePath = path.Join(tmp, "anime")
	log.Printf("[变动] 配置项(setting.save_path)：'%s' => '%s'\n", tmp, updateConfig.Setting.SavePath)
	updateConfig.Setting.DownloadPath = path.Join(tmp, "incomplete")
	log.Printf("[新增] 配置项(setting.download_path)：'%s'\n", updateConfig.Setting.DownloadPath)
	updateConfig.Advanced.Download.Rename = "link_delete"
	log.Printf("[新增] 配置项(setting.advanced.download.rename)：%v\n", updateConfig.Advanced.Download.Rename)

	log.Printf("[移除] 配置项(setting.advanced.download.queue_max_num)\n")
	log.Printf("[移除] 配置项(setting.advanced.download.queue_delay_second)\n")

	content, err := encodeConfig(updateConfig)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}
	err = os.WriteFile(file, content, 0644)
	if err != nil {
		log.Fatal("配置文件升级失败：", err)
	}

	log.Println("配置文件升级完成：1.0.0 => 1.1.0")
	log.Println("------------升级说明------------")
	log.Printf("存储动画的路径变更为 '%s'，请注意修改设置\n", updateConfig.Setting.SavePath)
}

func encodeConfig(conf *Config) ([]byte, error) {
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

func backupConfig(file string) error {
	dir, name := path.Split(file)
	ext := path.Ext(name)
	name = strings.TrimSuffix(name, ext)
	timeStr := time.Now().Format("20060102150405")
	name = fmt.Sprintf("%s-%s%s", name, timeStr, ext)
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
