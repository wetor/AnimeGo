package configs

import (
	"AnimeGo/internal/models"
	"AnimeGo/internal/utils"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func NewConfig(_path string) *Config {
	if len(_path) == 0 {
		_path = "../data/config/conf.yaml"
	}
	data, err := os.ReadFile(_path)
	if err != nil {
		zap.S().Fatal("配置文件加载错误：", err)
	}
	conf := &Config{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		zap.S().Fatal("配置文件加载错误：", err)
	}
	conf.CachePath = path.Join(conf.DataPath, conf.CachePath)
	conf.SavePath = path.Join(conf.DataPath, conf.SavePath)
	conf.Filter.JavaScript = path.Join(conf.DataPath, conf.Filter.JavaScript)
	return conf
}

func (c *Config) ClientQBt() *Client {
	if client, has := c.Client["qbittorrent"]; has {
		return client
	}
	return nil
}

func (c *Config) RssMikan() *Rss {
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

func (s *Setting) Tag(info *models.AnimeEntity) string {
	date, _ := time.Parse("2006-01-02", info.AirDate)
	mouth := (int(date.Month()) + 2) / 3
	tag := utils.Format(s.TagSrc, utils.FormatMap{
		"year":          date.Year(),
		"quarter":       (mouth-1)*3 + 1,
		"quarter_index": mouth,
		"quarter_name":  []string{"冬", "春", "夏", "秋"}[mouth-1],
		"ep":            info.Ep,
		"week":          (int(date.Weekday())+6)%7 + 1,
		"week_name":     []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}[date.Weekday()],
	})
	return tag
}
