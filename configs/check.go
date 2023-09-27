package configs

import "github.com/wetor/AnimeGo/pkg/log"

const (
	UpdateDelaySecondMin = 2
)

func (c *Config) Check() {
	if c.Advanced.RefreshSecond < UpdateDelaySecondMin {
		log.Warnf("配置项advanced.update_delay_second值范围错误: %v, 已修改为: %v", UpdateDelaySecondMin)
		c.Advanced.RefreshSecond = UpdateDelaySecondMin
	}
}
