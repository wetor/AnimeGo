package cache

import (
	"GoBangumi/config"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "10")
	flag.Parse()
	defer glog.Flush()
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
	v := db.Get(DefaultBucket, "key")
	fmt.Println(v)
}
