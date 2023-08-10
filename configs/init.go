package configs

import (
	"log"
	"os"
	"path"

	"github.com/caarlos0/env/v9"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/pkg/utils"
)

func InitUpdateConfig(configFile string, backup bool) {
	// 尝试升级配置文件
	if UpdateConfig(configFile, backup) {
		// 升级完成之后需要重启
		os.Exit(0)
	}
}

func InitDefaultConfig(configFile string) {
	log.Printf("未找到配置文件（%s），开始初始化默认配置\n", configFile)
	dir := path.Dir(configFile)
	if !utils.IsExist(dir) {
		err := utils.CreateMutiDir(dir)
		if err != nil {
			panic(err)
		}
	}
	// 写入默认配置文件
	err := DefaultFile(configFile)
	if err != nil {
		panic(err)
	}
	log.Printf("初始化默认配置完成（%s）\n", configFile)
}

func InitEnvConfig(configFile, saveFile string) {
	prefix := "ANIMEGO_"
	envs := &Environment{}
	opts := env.Options{
		Prefix: prefix,
	}
	err := env.ParseWithOptions(envs, opts)
	if err != nil {
		panic(err)
	}
	config := Load(configFile)
	err = Env2Config(envs, config, prefix)
	if err != nil {
		panic(err)
	}

	if envs.ProxyUrl != nil {
		if len(*envs.ProxyUrl) == 0 {
			config.Setting.Proxy.Enable = false
		} else {
			config.Setting.Proxy.Enable = true
		}
	}

	data, err := Config2Bytes(config)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(saveFile, data, constant.FilePerm)
	if err != nil {
		panic(err)
	}
}
