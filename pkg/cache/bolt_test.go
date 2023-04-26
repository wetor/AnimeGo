package cache_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
)

type AnimeEntity struct {
	ID      int    // bangumi id
	Name    string // 名称，从bgm获取
	NameCN  string // 中文名称，从bgm获取
	AirDate string // 最初播放日期，从bgm获取
	Eps     int    // 总集数，从bgm获取
	*AnimeExtra
}

type AnimeExtra struct {
	ThemoviedbID int    // themoviedb ID
	MikanID      int    // mikan id
	MikanUrl     string // mikan当前集的url
}

type Entity struct {
	ID      int    `json:"id"`      // Bangumi ID
	NameCN  string `json:"name_cn"` // 中文名
	Name    string `json:"name"`    // 原名
	Eps     int    `json:"eps"`     // 集数
	AirDate string `json:"airdate"` // 可空

	Type     int `json:"type"`
	Platform int `json:"platform"`
}

var (
	db     = cache.NewBolt()
	db_sub = cache.NewBolt(true)
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	db.Open("data/test.db")
	db_sub.Open("testdata/bolt_sub.bolt")
	m.Run()
	db.Close()
	db_sub.Close()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestBolt_Put(t *testing.T) {
	want := &AnimeEntity{
		ID:      666,
		Name:    "测试番剧名称",
		NameCN:  "测试番剧中文",
		AirDate: "2022-09-01",
		Eps:     10,
		AnimeExtra: &AnimeExtra{
			MikanID:      777,
			ThemoviedbID: 888,
		},
	}
	db.Add("test")
	db.Put("test", "key", want, 0)
	got := &AnimeEntity{}
	err := db.Get("test", "key", got)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Put() = %v, want %v", got, want)
	}
}

func TestBolt_BatchPut(t *testing.T) {
	want := []interface{}{
		&AnimeEntity{
			ID:      666,
			Name:    "测试番剧名称",
			NameCN:  "测试番剧中文",
			AirDate: "2022-09-01",
			Eps:     10,
			AnimeExtra: &AnimeExtra{
				MikanID:      666,
				ThemoviedbID: 888,
			},
		},
		&AnimeEntity{
			ID:      777,
			Name:    "测试番剧6666名称",
			NameCN:  "测试番剧中文",
			AirDate: "2022-01-01",
			Eps:     4,
			AnimeExtra: &AnimeExtra{
				MikanID:      777,
				ThemoviedbID: 888,
			},
		},
		&AnimeEntity{
			ID:      888,
			Name:    "测试番剧名asd称",
			NameCN:  "测试番剧中文",
			AirDate: "2022-06-01",
			Eps:     7,
			AnimeExtra: &AnimeExtra{
				MikanID:      888,
				ThemoviedbID: 88888,
			},
		},
	}
	ks := []interface{}{
		"key1", "key2", "key3",
	}
	db.Add("test")
	db.BatchPut("test", ks, want, 0)
	for i, k := range ks {
		got := &AnimeEntity{}
		err := db.Get("test", k, got)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(got, want[i]) {
			t.Errorf("BatchPut() = %v, want %v", got, want[i])
		}
	}
}

func TestBolt_GetAll(t *testing.T) {
	k := 10
	v := Entity{}
	db_sub.GetAll("bangumi_sub", &k, &v, func(k1, v1 interface{}) {
		fmt.Println(*k1.(*int))
		fmt.Println(*v1.(*Entity))
	})
}

func TestBolt_GetAll_Sub(t *testing.T) {

	key := 51
	got := &Entity{}
	db_sub.Get("bangumi_sub", key, got)
	want := &Entity{
		ID:       51,
		NameCN:   "CLANNAD",
		Name:     "CLANNAD -クラナド-",
		Eps:      22,
		AirDate:  "2007-10-04",
		Type:     2,
		Platform: 1,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get() = %v, want %v", got, want)
	}
}
