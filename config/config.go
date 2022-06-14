package config

import (
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Feed struct {
		Rss map[string]*Rss
	}

	Client map[string]*Client

	*Directory
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
type Directory struct {
	Download string
	Data     string
	Cache    string
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

func QBt() *Client {
	nameList := []string{"qbittorrent", "qBittorrent", "qbt", "qBt"}
	for _, name := range nameList {
		if client, has := conf.Client[name]; has {
			return client
		}
	}
	return nil
}
func Mikan() *Rss {
	nameList := []string{"mikan", "Mikan", "mikanani"}
	for _, name := range nameList {
		if rss, has := conf.Feed.Rss[name]; has {
			return rss
		}
	}
	return nil
}
func Dir() *Directory {
	return conf.Directory
}
