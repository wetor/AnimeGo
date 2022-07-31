package cache

import (
	"GoBangumi/models"
	"GoBangumi/store/config"
	"GoBangumi/utils/logger"
	"fmt"
	"testing"
	"time"
)

var conf *config.Config

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger.Init()
	defer logger.Flush()
	conf = config.NewConfig("../../data/config/conf.yaml")
	m.Run()
	fmt.Println("end")
}
func TestBolt_Put(t *testing.T) {

	db := NewBolt()
	db.Open(conf.Setting.CachePath)
	db.Put(models.DefaultBucket, "key", "这是测试文222", 0)
}

func TestBolt_Get(t *testing.T) {
	db := NewBolt()
	db.Open(conf.Setting.CachePath)
	v := db.Get(models.DefaultBucket, "key11")
	fmt.Println(v)
}

func TestToBytes(t *testing.T) {
	c := &Bolt{}
	b := c.toBytes(&models.Bangumi{
		ID:     1000,
		Name:   "测试日文",
		NameCN: "测试中文",
		BangumiExtra: &models.BangumiExtra{
			SubID:  22,
			SubUrl: "hasdtasdasdas",
		},
	}, time.Now().Unix()+30)
	fmt.Println(b)
	v, e := c.toValue(b)
	bgm := v.(*models.Bangumi)
	fmt.Println(bgm, bgm.BangumiExtra, e)
}
