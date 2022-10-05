package configs

import (
	"AnimeGo/internal/models"
	"AnimeGo/internal/utils"
	"log"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v3"
)

func Init(path string) *Config {
	if len(path) == 0 {
		path = "../data/config/animego.yaml"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}
	conf := &Config{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		log.Fatal("配置文件加载错误：", err)
	}
	return conf
}

func (c *Config) InitDir() {
	c.TempPath = path.Join(c.DataPath, c.TempPath)

	err := utils.CreateMutiDir(c.DataPath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", c.DataPath)
	}
	err = utils.CreateMutiDir(c.SavePath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", c.SavePath)
	}
	err = utils.CreateMutiDir(c.TempPath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", c.TempPath)
	}
	dbDir := path.Join(c.DataPath, path.Dir(c.DbFile))
	err = utils.CreateMutiDir(dbDir)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", dbDir)
	}
	logDir := path.Join(c.DataPath, path.Dir(c.LogFile))
	err = utils.CreateMutiDir(logDir)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", logDir)
	}

	c.DbFile = path.Join(c.DataPath, c.DbFile)
	c.LogFile = path.Join(c.DataPath, c.LogFile)
	c.JavaScript = path.Join(c.DataPath, c.JavaScript)
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
