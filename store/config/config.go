package config

import (
	"GoBangumi/models"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	*models.Config
}

func NewConfig(path string) *Config {
	if len(path) == 0 {
		path = "../data/config/conf.yaml"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		zap.S().Fatal("配置文件加载错误：", err)
	}
	conf := &Config{}
	conf.Config = &models.Config{}
	err = yaml.Unmarshal(data, conf.Config)
	if err != nil {
		zap.S().Fatal("配置文件加载错误：", err)
	}
	return conf
}

func (c *Config) ClientQBt() *models.Client {
	if client, has := c.Client["qbittorrent"]; has {
		return client
	}
	return nil
}
func (c *Config) RssMikan() *models.Rss {
	if rss, has := c.Feed.Rss["mikan"]; has {
		return rss
	}
	return nil
}
func (c *Config) KeyTmdb() string {
	if key, has := c.Key["themoviedb"]; has {
		return key
	}
	return ""
}

func (c *Config) Proxy() string {
	if c.ProxyConf.Enable {
		return c.ProxyConf.Url
	} else {
		return ""
	}
}
