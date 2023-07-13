package configs

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

var ConfigFile = "./data/animego.yaml"

func Load(file string) *Config {
	if len(file) == 0 {
		file = xpath.Abs(ConfigFile)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}
	ConfigFile = file

	conf := &Config{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}

	return conf
}

func (c *Config) InitDir() {

	err := utils.CreateMutiDir(c.DataPath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", c.DataPath)
	}

	err = utils.CreateMutiDir(c.DownloadPath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", c.DownloadPath)
	}

	err = utils.CreateMutiDir(c.SavePath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", c.SavePath)
	}

	err = utils.CreateMutiDir(constant.CachePath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", constant.CachePath)
	}

	err = utils.CreateMutiDir(constant.LogPath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", constant.LogPath)
	}

	err = utils.CreateMutiDir(constant.TempPath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", constant.TempPath)
	}

	// 路径设置转绝对路径
	c.Setting.DataPath = xpath.Abs(c.Setting.DataPath)
	c.Setting.DownloadPath = xpath.Abs(c.Setting.DownloadPath)
	c.Setting.SavePath = xpath.Abs(c.Setting.SavePath)
}

func (c *Config) Proxy() string {
	if c.Setting.Proxy.Enable {
		return c.Setting.Proxy.Url
	} else {
		return ""
	}
}
