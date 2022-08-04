package themoviedb

import (
	"GoBangumi/internal/models"
	"GoBangumi/store"
	"GoBangumi/utils/logger"
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger.Init()
	defer logger.Flush()
	store.Init(nil)
	m.Run()
	fmt.Println("end")
}
func TestNewThemoviedb(t *testing.T) {
	tmdb := NewThemoviedb()
	b := tmdb.Parse(&models.AnimeParseOptions{
		Name: "カードファイト!! ヴァンガード will+Dress",
		Date: "2022-07-04",
	})
	fmt.Println(b)
}

func TestTime(t *testing.T) {

	//fmt.Println(utils.GetTimeRangeDay("2022-04-10", 1))
}
