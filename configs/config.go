package configs

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
	"log"
	"os"
	"path"
	"path/filepath"
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
	c.Path.TempPath = path.Join(c.DataPath, c.Path.TempPath)

	err := utils.CreateMutiDir(c.DataPath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", c.DataPath)
	}
	err = utils.CreateMutiDir(c.SavePath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", c.SavePath)
	}
	absPath, err := filepath.Abs(c.SavePath)
	if err != nil {
		log.Fatalf("save_path不是正确的路径，%s", c.SavePath)
	}
	c.SavePath = absPath
	err = utils.CreateMutiDir(c.Path.TempPath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", c.Path.TempPath)
	}
	dbDir := path.Join(c.DataPath, path.Dir(c.Path.DbFile))
	err = utils.CreateMutiDir(dbDir)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", dbDir)
	}
	logDir := path.Join(c.DataPath, path.Dir(c.Path.LogFile))
	err = utils.CreateMutiDir(logDir)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", logDir)
	}

	c.Path.DbFile = path.Join(c.DataPath, c.Path.DbFile)
	c.Path.LogFile = path.Join(c.DataPath, c.Path.LogFile)
	for i := range c.Filter.JavaScript {
		c.Filter.JavaScript[i] = path.Join(c.DataPath, c.Filter.JavaScript[i])
	}
}

func (c *Config) Proxy() string {
	if c.Setting.Proxy.Enable {
		return c.Setting.Proxy.Url
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
