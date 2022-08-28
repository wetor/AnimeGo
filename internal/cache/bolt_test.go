package cache

import (
	"AnimeGo/internal/models"
	"fmt"
	"testing"
)

func TestBolt_Put(t *testing.T) {

	db := NewBolt()
	db.Open(".")
	e := models.AnimeEntity{
		ID:      666,
		Name:    "测试番剧名称",
		NameCN:  "测试番剧中文",
		AirDate: "2022-09-01",
		Eps:     10,
		AnimeExtra: &models.AnimeExtra{
			MikanID:      777,
			ThemoviedbID: 888,
		},
	}
	db.Put("test", "key", e, 0)
}

func TestBolt_Get(t *testing.T) {
	db := NewBolt()
	db.Open(".")
	var data models.AnimeEntity
	err := db.Get("test", "key", &data)
	if err != nil {
		panic(err)
	}
	fmt.Println(data, data.AnimeExtra)
}
