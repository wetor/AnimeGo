package themoviedb

import (
	"encoding/gob"
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/pkg/errors"
	mem "github.com/wetor/AnimeGo/pkg/memorizer"
	"github.com/wetor/AnimeGo/pkg/request"
)

var (
	Host                    = "https://api.themoviedb.org"
	Bucket                  = "themoviedb"
	MatchSeasonDays         = 90
	MinSimilar      float64 = 0.75
)

type Themoviedb struct {
	Key                    string
	cacheInit              bool
	cacheParseThemoviedbID mem.Func
	cacheParseAnimeSeason  mem.Func
}

func (t *Themoviedb) RegisterCache() {
	if anidata.Cache == nil {
		panic(errors.NewAniError("需要先调用anidata.Init初始化缓存"))
	}
	t.cacheInit = true
	t.cacheParseThemoviedbID = mem.Memorized(Bucket, anidata.Cache, func(params *mem.Params, results *mem.Results) error {
		entity, err := t.parseThemoviedbID(params.Get("name").(string))
		if err != nil {
			return err
		}
		results.Set("entity", entity)
		return nil
	})

	t.cacheParseAnimeSeason = mem.Memorized(Bucket, anidata.Cache, func(params *mem.Params, results *mem.Results) error {
		seasonInfo, err := t.parseAnimeSeason(params.Get("tmdbID").(int), params.Get("airDate").(string))
		if err != nil {
			return err
		}
		results.Set("seasonInfo", seasonInfo)
		return nil
	})
}

func (t Themoviedb) ParseCache(name, airDate string) (entity *Entity, seasonInfo *SeasonInfo, err error) {
	if !t.cacheInit {
		t.RegisterCache()
	}
	results := mem.NewResults("entity", &Entity{}, "seasonInfo", &SeasonInfo{})

	err = t.cacheParseThemoviedbID(mem.NewParams("name", name).
		TTL(anidata.CacheTime[Bucket]), results)
	if err != nil {
		return
	}
	entity = results.Get("entity").(*Entity)
	err = t.cacheParseAnimeSeason(mem.NewParams("tmdbID", entity.ID, "airDate", airDate).
		TTL(anidata.CacheTime[Bucket]), results)
	if err != nil {
		return
	}
	seasonInfo = results.Get("seasonInfo").(*SeasonInfo)
	return
}

func (t Themoviedb) Parse(name, airDate string) (entity *Entity, seasonInfo *SeasonInfo, err error) {
	entity, err = t.parseThemoviedbID(name)
	if err != nil {
		return
	}
	seasonInfo, err = t.parseAnimeSeason(entity.ID, airDate)
	if err != nil {
		return
	}
	return
}

func (t Themoviedb) parseThemoviedbID(name string) (entity *Entity, err error) {
	resp := FindResponse{}
	result, err := RemoveNameSuffix(name, func(innerName string) (interface{}, error) {
		fmt.Println(idApi(t.Key, innerName))
		err := request.Get(idApi(t.Key, innerName), &resp)
		if err != nil {
			return nil, err
		}
		if resp.TotalResults == 1 {
			return resp.Result[0], nil
		} else if resp.TotalResults > 1 {
			// 筛选与original name完全相同的番剧
			for _, result := range resp.Result {
				if result.Name == name {
					return result, nil
				}
			}

			// 按照相似度排序筛选
			temp := &Entity{}
			maxSimilar := float64(0)
			for _, result := range resp.Result {
				similar := SimilarText(result.Name, name)
				if similar > maxSimilar {
					maxSimilar = similar
					temp = result
				}
			}
			if maxSimilar >= MinSimilar {
				return temp, nil
			}
			return 0, errors.NewAniError("匹配Seasons失败，番剧名未找到")
		} else {
			// 未找到结果
			return nil, nil
		}
	})
	if err != nil {
		return nil, err
	}

	return result.(*Entity), nil
}

func (t Themoviedb) parseAnimeSeason(tmdbID int, airDate string) (seasonInfo *SeasonInfo, err error) {
	resp := InfoResponse{}
	err = request.Get(infoApi(t.Key, tmdbID), &resp)
	if err != nil {
		return nil, err
	}
	if resp.Seasons == nil || len(resp.Seasons) == 0 {
		return nil, errors.NewAniError("匹配Seasons失败，可能此番剧未开播")
	}
	seasonInfo = resp.Seasons[0]
	min := 36500
	for _, r := range resp.Seasons {
		if r.Season == 0 || r.EpName == "Specials" {
			continue
		}
		// TODO: 待优化，通过比较此季度番剧的初放送日期，筛选差值最小的季
		if s := StrTimeSubAbs(r.AirDate, airDate); s < min {
			min = s
			seasonInfo = r
		}
	}
	if min > MatchSeasonDays {
		return nil, errors.NewAniError("匹配Seasons失败，可能此番剧未开播")
	}
	seasonInfo.EpName = ""
	return seasonInfo, nil
}

func init() {
	gob.Register(&Entity{})
	gob.Register(&SeasonInfo{})
}
