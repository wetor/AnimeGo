package memorizer

import (
	"encoding/json"
	"fmt"
	"github.com/wetor/AnimeGo/pkg/cache"
	"go.uber.org/zap"
	"strconv"
	"testing"
	"time"
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
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	m.Run()
	fmt.Println("end")
}

func toString(in interface{}) string {
	data, _ := json.Marshal(in)
	return string(data)
}

func DoSomething(params *Params, results *Results) error {

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
	dosomething := Memorized("test", db, DoSomething)

	res := NewResults("ThemovieID", &Type{})
	dosomething(NewParams("mikanID", 1001, "bangumiID", 3333).TTL(1), res)
	fmt.Println(toString(res))

	dosomething(NewParams("mikanID", 1001, "bangumiID", 3333), res)
	fmt.Println(toString(res))

	time.Sleep(2 * time.Second)

	dosomething(NewParams("mikanID", 1001, "bangumiID", 3333).TTL(1), res)
	fmt.Println(toString(res))
}
