package themoviedb

import (
	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/pkg/errors"
	mem "github.com/wetor/AnimeGo/pkg/memorizer"
	"github.com/wetor/AnimeGo/pkg/request"
)

var (
	Host = func() string {
		if len(anidata.RedirectThemoviedb) > 0 {
			return anidata.RedirectThemoviedb
		}
		return "https://api.themoviedb.org"
	}
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
		errors.NewAniError("需要先调用anidata.Init初始化缓存").TryPanic()
	}
	t.cacheInit = true
	t.cacheParseThemoviedbID = mem.Memorized(Bucket, anidata.Cache, func(params *mem.Params, results *mem.Results) error {
		entity := t.parseThemoviedbID(params.Get("name").(string))
		results.Set("entity", entity)
		return nil
	})

	t.cacheParseAnimeSeason = mem.Memorized(Bucket, anidata.Cache, func(params *mem.Params, results *mem.Results) error {
		seasonInfo := t.parseAnimeSeason(params.Get("tmdbID").(int), params.Get("airDate").(string))
		results.Set("seasonInfo", seasonInfo)
		return nil
	})
}

func (t *Themoviedb) Search(name string) int {
	entity := t.parseThemoviedbID(name)
	return entity.ID
}

func (t *Themoviedb) SearchCache(name string) int {
	if !t.cacheInit {
		t.RegisterCache()
	}
	results := mem.NewResults("entity", &Entity{})
	err := t.cacheParseThemoviedbID(mem.NewParams("name", name).
		TTL(anidata.CacheTime[Bucket]), results)
	errors.NewAniErrorD(err).TryPanic()
	entity := results.Get("entity").(*Entity)
	return entity.ID
}

func (t *Themoviedb) Get(id int, filters any) any {
	airDate := filters.(string)
	seasonInfo := t.parseAnimeSeason(id, airDate)
	return seasonInfo
}

func (t *Themoviedb) GetCache(id int, filters any) any {
	if !t.cacheInit {
		t.RegisterCache()
	}
	airDate := filters.(string)
	results := mem.NewResults("seasonInfo", &SeasonInfo{})
	err := t.cacheParseAnimeSeason(mem.NewParams("tmdbID", id, "airDate", airDate).
		TTL(anidata.CacheTime[Bucket]), results)
	errors.NewAniErrorD(err).TryPanic()
	seasonInfo := results.Get("seasonInfo").(*SeasonInfo)
	return seasonInfo
}

func (t *Themoviedb) parseThemoviedbID(name string) (entity *Entity) {
	resp := FindResponse{}
	result := RemoveNameSuffix(name, func(innerName string) any {
		err := request.Get(idApi(t.Key, innerName), &resp)
		errors.NewAniErrorD(err).TryPanic()

		if resp.TotalResults == 1 {
			return resp.Result[0]
		} else if resp.TotalResults > 1 {
			// 筛选与original name完全相同的番剧
			for _, result := range resp.Result {
				if result.Name == name {
					return result
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
				return temp
			}
			errors.NewAniError("匹配Seasons失败，番剧名未找到").TryPanic()
		} else {
			// 未找到结果
			return nil
		}
		return nil
	})

	return result.(*Entity)
}

func (t *Themoviedb) parseAnimeSeason(tmdbID int, airDate string) (seasonInfo *SeasonInfo) {
	resp := InfoResponse{}
	err := request.Get(infoApi(t.Key, tmdbID), &resp)
	errors.NewAniErrorD(err).TryPanic()
	if resp.Seasons == nil || len(resp.Seasons) == 0 {
		errors.NewAniError("匹配Seasons失败，可能此番剧未开播").TryPanic()
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
		errors.NewAniError("匹配Seasons失败，可能此番剧未开播").TryPanic()
	}
	seasonInfo.EpName = ""
	return seasonInfo
}

// Check interface is satisfied
var _ api.AniDataSearchGet = &Themoviedb{}
