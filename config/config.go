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

	Client    map[string]*Client
	Key       map[string]string
	*Settings `yaml:"setting"`
	Proxy     struct {
		Enable bool
		Url    string
	}
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
func TMDB() string {
	nameList := []string{"themoviedb", "TheMovieDB", "tmdb", "TMDB"}
	for _, name := range nameList {
		if key, has := conf.Key[name]; has {
			return key
		}
	}
	return ""
}
func Setting() *Settings {
	return conf.Settings
}

func Proxy() string {
	if conf.Proxy.Enable {
		return conf.Proxy.Url
	} else {
		return ""
	}
}
