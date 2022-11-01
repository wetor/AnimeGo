package cache

import (
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

	db := NewBolt()
	db.Open("1.a")
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
	db.Put("test", "key", e, 0)
}

func TestBolt_Get(t *testing.T) {
	db := NewBolt()
	db.Open("bolt.db")
	var data AnimeEntity
	err := db.Get("test", "key", &data)
	if err != nil {
		panic(err)
	}
	fmt.Println(data, data.AnimeExtra)
}
