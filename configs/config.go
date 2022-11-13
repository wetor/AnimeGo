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

	conf.Path.TempPath = path.Join(conf.DataPath, conf.Path.TempPath)
	absPath, err := filepath.Abs(conf.SavePath)
	if err != nil {
		log.Fatalf("save_path不是正确的路径，%s", conf.SavePath)
	}
	conf.SavePath = absPath

	conf.Path.DbFile = path.Join(conf.DataPath, conf.Path.DbFile)
	conf.Path.LogFile = path.Join(conf.DataPath, conf.Path.LogFile)
	for i := range conf.Filter.JavaScript {
		conf.Filter.JavaScript[i] = path.Join(conf.DataPath, conf.Filter.JavaScript[i])
	}
	return conf
}

func (c *Config) InitDir() {

	err := utils.CreateMutiDir(c.DataPath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", c.DataPath)
	}
	err = utils.CreateMutiDir(c.SavePath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", c.SavePath)
	}

	err = utils.CreateMutiDir(c.Path.TempPath)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", c.Path.TempPath)
	}
	dbDir := path.Dir(c.Path.DbFile)
	err = utils.CreateMutiDir(dbDir)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", dbDir)
	}
	logDir := path.Dir(c.Path.LogFile)
	err = utils.CreateMutiDir(logDir)
	if err != nil {
		log.Fatalf("创建文件夹失败，%s", logDir)
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
