package cache

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
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

func TestBolt_Put(t *testing.T) {
	db := NewBolt()
	db.Open("1.db")
	e := AnimeEntity{
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
	db.Put("test", "key", e, 0)
}

func TestBolt_Get(t *testing.T) {
	db := NewBolt()
	db.Open("1.db")
	var data AnimeEntity
	err := db.Get("test", "key3", &data)
	if err != nil {
		panic(err)
	}
	fmt.Println(data, data.AnimeExtra)
}

func TestBolt_BatchPut(t *testing.T) {
	db := NewBolt()
	db.Open("1.db")
	es := []interface{}{
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
	db.BatchPut("test", ks, es, 0)
}

func TestBolt_List(t *testing.T) {
	db := NewBolt()
	db.Open("/Users/wetor/GoProjects/AnimeGo/data/cache/bolt.db")
	db.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("hash2name"))

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k[1:], v[8:])
			return nil
		})
		return nil
	})
}

func TestBolt_GetAll(t *testing.T) {
	db := NewBolt()
	db.Open("/Users/wetor/GoProjects/AnimeGo/data/cache/bolt.db")
	v := ""
	k := ""
	db.GetAll("hash2name", &k, &v, func(k1, v1 interface{}) {
		fmt.Println("-----")
		fmt.Println(*k1.(*string))
		fmt.Println(*v1.(*string))
	})
}
