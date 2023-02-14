package memorizer_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/memorizer"
)

type Type struct {
	ID  int
	Str string
	*Type2
}

type Type2 struct {
	Str2 string
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/test.log",
		Debug: true,
	})
	m.Run()
	fmt.Println("end")
}

func toString(in interface{}) string {
	data, _ := json.Marshal(in)
	return string(data)
}

func DoSomething(params *memorizer.Params, results *memorizer.Results) error {

	fmt.Print("not cache: ")
	mikanID := params.Get("mikanID").(int)
	bangumiID := params.Get("bangumiID").(int)

	results.Set("ThemovieID", &Type{
		ID:  mikanID + bangumiID,
		Str: strconv.Itoa(mikanID),
		Type2: &Type2{
			Str2: strconv.Itoa(bangumiID),
		},
	})
	return nil
}

func TestMemorized(t *testing.T) {
	db := cache.NewBolt()
	db.Open("data/bolt.db")
	dosomething := memorizer.Memorized("test", db, DoSomething)

	res := memorizer.NewResults("ThemovieID", &Type{})
	dosomething(memorizer.NewParams("mikanID", 1001, "bangumiID", 3333).TTL(1), res)
	fmt.Println(toString(res))

	dosomething(memorizer.NewParams("mikanID", 1001, "bangumiID", 3333), res)
	fmt.Println(toString(res))

	time.Sleep(2 * time.Second)

	dosomething(memorizer.NewParams("mikanID", 1001, "bangumiID", 3333).TTL(1), res)
	fmt.Println(toString(res))
}
