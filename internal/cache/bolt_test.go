package cache

import (
	"GoBangumi/configs"
	"GoBangumi/internal/models"
	"GoBangumi/utils/logger"
	"fmt"
	"testing"
	"time"
)

var conf *configs.Config

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger.Init()
	defer logger.Flush()
	conf = configs.NewConfig("../../data/config/conf.yaml")
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
	b := c.toBytes(&models.AnimeEntity{
		ID:     1000,
		Name:   "测试日文",
		NameCN: "测试中文",
		AnimeExtra: &models.AnimeExtra{
			MikanID:  22,
			MikanUrl: "hasdtasdasdas",
		},
	}, time.Now().Unix()+30)
	fmt.Println(b)
	v, e := c.toValue(b)
	anime := v.(*models.AnimeEntity)
	fmt.Println(anime, anime.AnimeExtra, e)
}
