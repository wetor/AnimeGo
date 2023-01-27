package configs

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/utils"
)

var ConfigFile = "./data/animego.yaml"

func Init(file string) *Config {
	if len(file) == 0 {
		file = ConfigFile
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

	absPath, err := filepath.Abs(conf.DownloadPath)
	if err != nil {
		log.Fatalf("download_path不是正确的路径，%s", conf.DownloadPath)
	}
	conf.DownloadPath = absPath

	absPath, err = filepath.Abs(conf.SavePath)
	if err != nil {
		log.Fatalf("save_path不是正确的路径，%s", conf.SavePath)
	}
	conf.SavePath = absPath

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

	err = utils.CreateMutiDir(constant.TempPath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", constant.TempPath)
	}

	err = utils.CreateMutiDir(constant.CachePath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", constant.CachePath)
	}

	err = utils.CreateMutiDir(constant.LogPath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", constant.LogPath)
	}
}

func (c *Config) Proxy() string {
	if c.Setting.Proxy.Enable {
		return c.Setting.Proxy.Url
	} else {
		return ""
	}
}
