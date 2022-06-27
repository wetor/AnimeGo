package config

import (
	"github.com/golang/glog"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Feed struct {
		Rss map[string]*Rss
	}

	Client map[string]*Client
	Key    map[string]string
	Proxy  struct {
		Enable bool
		Url    string
	}
	*Settings     `yaml:"setting"`
	*AdvancedConf `yaml:"advanced"`
}

type Client struct {
	Url      string
	Username string
	Password string
}
type Rss struct {
	Name string
	Url  string
}

var conf = &Config{}

func Init(path string) {
	if len(path) == 0 {
		path = "../data/config/conf.yaml"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		glog.Errorln("配置文件加载错误：", err)
	}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		glog.Errorln("配置文件加载错误：", err)
	}
}

func ClientQBt() *Client {
	if client, has := conf.Client["qbittorrent"]; has {
		return client
	}
	return nil
}
func RssMikan() *Rss {
	if rss, has := conf.Feed.Rss["mikan"]; has {
		return rss
	}
	return nil
}
func KeyTmdb() string {
	if key, has := conf.Key["themoviedb"]; has {
		return key
	}
	return ""
}
func Setting() *Settings {
	return conf.Settings
}

func Advanced() *AdvancedConf {
	return conf.AdvancedConf
}

func Proxy() string {
	if conf.Proxy.Enable {
		return conf.Proxy.Url
	} else {
		return ""
	}
}
