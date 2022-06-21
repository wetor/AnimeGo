package config

import (
	"GoBangumi/models"
	"fmt"
	"strings"
	"time"
)

type Settings struct {
	DataPath  string `yaml:"data_path"`
	CachePath string `yaml:"cache_path"`

	SavePath    string `yaml:"save_path"`
	Category    string `yaml:"category"`     // 分类
	TagSrc      string `yaml:"tag"`          // 标签
	SeedingTime int    `yaml:"seeding_time"` // 做种时间，单位：分钟
}

type tagFormat map[string]interface{}

func format(format string, p tagFormat) string {
	args, i := make([]string, len(p)*2), 0
	for k, v := range p {
		args[i] = "{" + k + "}"
		args[i+1] = fmt.Sprint(v)
		i += 2
	}
	return strings.NewReplacer(args...).Replace(format)
}

func (s *Settings) Tag(info *models.Bangumi) string {
	date, _ := time.Parse("2006-01-02", info.AirDate)
	mouth := (int(date.Month()) + 2) / 3
	str := format(s.TagSrc, tagFormat{
		"year":          date.Year(),
		"quarter":       (mouth-1)*3 + 1,
		"quarter_index": mouth,
		"quarter_name":  []string{"冬", "春", "夏", "秋"}[mouth-1],
		"ep":            info.Ep,
		"week":          (int(date.Weekday())+6)%7 + 1,
		"week_name":     []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}[date.Weekday()],
	})
	return str
}
