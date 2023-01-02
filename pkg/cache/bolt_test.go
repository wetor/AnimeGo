package cache

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
	"reflect"
	"testing"
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

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	m.Run()
	fmt.Println("end")
}

func TestBolt_Put(t *testing.T) {
	db := NewBolt()
	db.Open("data/1.db")
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
	db := NewBolt()
	db.Open("data/1.db")
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

func TestBolt_List(t *testing.T) {
	db := NewBolt()
	db.Open("data/1.db")
	db.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("test"))

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v[8:])
			return nil
		})
		return nil
	})
}

func TestBolt_GetAll(t *testing.T) {
	db := NewBolt()
	db.Open("data/1.db")
	k := ""
	v := AnimeEntity{}
	db.GetAll("test", &k, &v, func(k1, v1 interface{}) {
		fmt.Println("-----")
		fmt.Println(*k1.(*string))
		fmt.Println(*v1.(*AnimeEntity))
	})
}

func TestBolt_GetAll_Sub(t *testing.T) {
	type Entity struct {
		ID      int    `json:"id"`      // Bangumi ID
		NameCN  string `json:"name_cn"` // 中文名
		Name    string `json:"name"`    // 原名
		Eps     int    `json:"eps"`     // 集数
		AirDate string `json:"airdate"` // 可空

		Type     int `json:"type"`
		Platform int `json:"platform"`
	}
	db := NewBolt()
	db.Open("data/bolt_sub.db")
	key := 302286
	got := &Entity{}
	db.Get("bangumi_sub", key, got)
	want := &Entity{
		ID:       302286,
		NameCN:   "死神 千年血战篇",
		Name:     "BLEACH 千年血戦篇",
		Eps:      13,
		AirDate:  "2022-10-10",
		Type:     2,
		Platform: 1,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get() = %v, want %v", got, want)
	}
}
