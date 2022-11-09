package cache

import (
	"encoding/gob"
	"fmt"
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
	gob.Register(&AnimeEntity{})
	gob.Register(&AnimeExtra{})
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
	gob.Register(&AnimeEntity{})
	gob.Register(&AnimeExtra{})
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
