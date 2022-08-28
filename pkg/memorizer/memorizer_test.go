package memorizer

import (
	"GoBangumi/internal/cache"
	"GoBangumi/internal/models"
	"encoding/json"
	"fmt"
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

func toString(in interface{}) string {
	data, _ := json.Marshal(in)
	return string(data)
}

func DoSomething(params *Params, results *Results) error {

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
	db.Open(".")
	dosomething := Memorized(models.DefaultBucket, db, DoSomething)

	res := NewResults("ThemovieID", &Type{})
	dosomething(NewParams("mikanID", 1001, "bangumiID", 3333).TTL(1), res)
	fmt.Println(toString(res))

	dosomething(NewParams("mikanID", 1001, "bangumiID", 3333), res)
	fmt.Println(toString(res))

	time.Sleep(2 * time.Second)

	dosomething(NewParams("mikanID", 1001, "bangumiID", 3333), res)
	fmt.Println(toString(res))
}
