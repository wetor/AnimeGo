package themoviedb

import (
	"github.com/wetor/AnimeGo/pkg/anisource"
	"github.com/wetor/AnimeGo/pkg/errors"
	mem "github.com/wetor/AnimeGo/pkg/memorizer"
	"github.com/wetor/AnimeGo/pkg/request"
)

var (
	Host                    = "https://api.themoviedb.org"
	Bucket                  = "themoviedb"
	MatchSeasonDays         = 90
	CacheSecond     int64   = 7 * 24 * 60 * 60
	MinSimilar      float64 = 0.75
)

type Themoviedb struct {
	Key                    string
	cacheInit              bool
	cacheParseThemoviedbID mem.Func
	cacheParseAnimeSeason  mem.Func
}

func (t *Themoviedb) RegisterCache() {
	if anisource.Cache == nil {
		panic(errors.NewAniError("需要先调用anisource.Init初始化缓存"))
	}
	t.cacheInit = true
	t.cacheParseThemoviedbID = mem.Memorized(Bucket, anisource.Cache, func(params *mem.Params, results *mem.Results) error {
		tmdbID, err := t.parseThemoviedbID(params.Get("name").(string))
		if err != nil {
			return err
		}
		results.Set("tmdbID", tmdbID)
		return nil
	})

	t.cacheParseAnimeSeason = mem.Memorized(Bucket, anisource.Cache, func(params *mem.Params, results *mem.Results) error {
		season, err := t.parseAnimeSeason(params.Get("tmdbID").(int), params.Get("airDate").(string))
		if err != nil {
			return err
		}
		results.Set("season", season)
		return nil
	})
}

func (t Themoviedb) ParseCache(name, airDate string) (tmdbID int, season int, err error) {
	if !t.cacheInit {
		t.RegisterCache()
	}
	results := mem.NewResults("tmdbID", 0, "season", 0)

	err = t.cacheParseThemoviedbID(mem.NewParams("name", name).TTL(CacheSecond), results)
	if err != nil {
		return
	}
	tmdbID = results.Get("tmdbID").(int)
	err = t.cacheParseAnimeSeason(mem.NewParams("tmdbID", tmdbID, "airDate", airDate).TTL(CacheSecond), results)
	if err != nil {
		return
	}
	season = results.Get("season").(int)
	return
}

func (t Themoviedb) Parse(name, airDate string) (tmdbID int, season int, err error) {
	tmdbID, err = t.parseThemoviedbID(name)
	if err != nil {
		return
	}
	season, err = t.parseAnimeSeason(tmdbID, airDate)
	if err != nil {
		return
	}
	return
}

func (t Themoviedb) parseThemoviedbID(name string) (tmdbID int, err error) {
	resp := FindResponse{}
	result, err := RemoveNameSuffix(name, func(innerName string) (interface{}, error) {
		err := request.Get(idApi(t.Key, innerName), &resp)
		if err != nil {
			return 0, err
		}
		if resp.TotalResults == 1 {
			return resp.Result[0].ID, nil
		} else if resp.TotalResults > 1 {
			// 筛选与original name完全相同的番剧
			for _, result := range resp.Result {
				if result.OriginalName == name {
					return result.ID, nil
				}
			}
			tmdbID = resp.Result[0].ID
			// 按照相似度排序筛选
			maxSimilar := float64(0)
			for _, result := range resp.Result {
				similar := SimilarText(result.OriginalName, name)
				if similar > maxSimilar {
					maxSimilar = similar
					tmdbID = result.ID
				}
			}
			if maxSimilar >= MinSimilar {
				return tmdbID, nil
			}
			return 0, errors.NewAniError("匹配Seasons失败，番剧名未找到")
		} else {
			// 未找到结果
			return nil, nil
		}
	})
	if err != nil {
		return 0, err
	}
	return result.(int), nil
}

func (t Themoviedb) parseAnimeSeason(tmdbID int, airDate string) (season int, err error) {
	resp := InfoResponse{}
	err = request.Get(infoApi(t.Key, tmdbID), &resp)
	if err != nil {
		return 0, err
	}
	if resp.Seasons == nil || len(resp.Seasons) == 0 {
		return 0, errors.NewAniError("匹配Seasons失败，可能此番剧未开播")
	}
	season = resp.Seasons[0].SeasonNumber
	min := 36500
	for _, r := range resp.Seasons {
		if r.SeasonNumber == 0 || r.Name == "Specials" {
			continue
		}
		// TODO: 待优化，通过比较此季度番剧的初放送日期，筛选差值最小的季
		if s := StrTimeSubAbs(r.AirDate, airDate); s < min {
			min = s
			season = r.SeasonNumber
		}
	}
	if min > MatchSeasonDays {
		return 0, errors.NewAniError("匹配Seasons失败，可能此番剧未开播")
	}
	return season, nil
}
