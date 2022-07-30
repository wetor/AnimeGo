package cache

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/utils/logger"
	"fmt"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger.Init()
	defer logger.Flush()
	config.Init("../../data/config/conf.yaml")
	m.Run()
	fmt.Println("end")
}
func TestBolt_Put(t *testing.T) {

	db := NewBolt()
	db.Open(config.Setting().CachePath)
	db.Put(DefaultBucket, "key", "这是测试文222", 0)
}

func TestBolt_Get(t *testing.T) {
	db := NewBolt()
	db.Open(config.Setting().CachePath)
	v := db.Get(DefaultBucket, "key11")
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
