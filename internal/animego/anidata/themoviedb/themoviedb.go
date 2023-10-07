package themoviedb

import (
	"github.com/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/utils"

	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/pkg/log"
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

func (a *Themoviedb) Name() string {
	return "Themoviedb"
}

func (a *Themoviedb) RegisterCache() {
	a.cacheInit = true
	a.cacheParseThemoviedbID = mem.Memorized(Bucket, anidata.Cache, func(params *mem.Params, results *mem.Results) error {
		entity, err := a.parseThemoviedbID(params.Get("name").(string))
		if err != nil {
			return err
		}
		results.Set("entity", entity)
		return nil
	})

	a.cacheParseAnimeSeason = mem.Memorized(Bucket, anidata.Cache, func(params *mem.Params, results *mem.Results) error {
		seasonInfo, err := a.parseAnimeSeason(params.Get("tmdbID").(int), params.Get("airDate").(string))
		if err != nil {
			return err
		}
		results.Set("seasonInfo", seasonInfo)
		return nil
	})
}

func (a *Themoviedb) Search(name string, filters any) (int, error) {
	entity, err := a.parseThemoviedbID(name)
	if err != nil {
		return 0, errors.Wrap(err, "查询ThemoviedbID失败")
	}
	return entity.ID, nil
}

func (a *Themoviedb) SearchCache(name string, filters any) (int, error) {
	if !a.cacheInit {
		a.RegisterCache()
	}
	results := mem.NewResults("entity", &Entity{})
	err := a.cacheParseThemoviedbID(mem.NewParams("name", name).
		TTL(anidata.CacheTime[Bucket]), results)
	if err != nil {
		return 0, errors.Wrap(err, "查询ThemoviedbID失败")
	}
	entity := results.Get("entity").(*Entity)
	return entity.ID, nil
}

func (a *Themoviedb) Get(id int, filters any) (any, error) {
	airDate := filters.(string)
	seasonInfo, err := a.parseAnimeSeason(id, airDate)
	if err != nil {
		return nil, errors.Wrap(err, "获取Themoviedb信息失败")
	}
	return seasonInfo, nil
}

func (a *Themoviedb) GetCache(id int, filters any) (any, error) {
	if !a.cacheInit {
		a.RegisterCache()
	}
	airDate := filters.(string)
	results := mem.NewResults("seasonInfo", &SeasonInfo{})
	err := a.cacheParseAnimeSeason(mem.NewParams("tmdbID", id, "airDate", airDate).
		TTL(anidata.CacheTime[Bucket]), results)
	if err != nil {
		return nil, errors.Wrap(err, "获取Themoviedb信息失败")
	}
	seasonInfo := results.Get("seasonInfo").(*SeasonInfo)
	return seasonInfo, nil
}

func (a *Themoviedb) parseThemoviedbID(name string) (entity *Entity, err error) {
	resp := FindResponse{}
	result, err := utils.RemoveNameSuffix(name, func(innerName string) (any, error) {
		err := request.Get(idApi(a.Key, innerName, false), &resp)
		if err != nil {
			log.DebugErr(err)
			//return 0, errors.WithStack(&exceptions.ErrRequest{Name: a.Name()})
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
				similar := utils.SimilarText(result.Name, name)
				if similar > maxSimilar {
					maxSimilar = similar
					temp = result
				}
			}
			if maxSimilar >= MinSimilar {
				return temp, nil
			}
			err = errors.WithStack(&exceptions.ErrThemoviedbMatchSeason{Message: "番剧名未找到"})
			log.DebugErr(err)
			return nil, err
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

func (a *Themoviedb) parseAnimeSeason(tmdbID int, airDate string) (seasonInfo *SeasonInfo, err error) {
	resp := InfoResponse{}
	err = request.Get(infoApi(a.Key, tmdbID, false), &resp)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrRequest{Name: a.Name()})
	}
	if resp.Seasons == nil || len(resp.Seasons) == 0 {
		err = errors.WithStack(&exceptions.ErrThemoviedbMatchSeason{Message: "此番剧可能未开播"})
		log.DebugErr(err)
		return nil, err
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
		err = errors.WithStack(&exceptions.ErrThemoviedbMatchSeason{Message: "此番剧可能未开播"})
		log.DebugErr(err)
		return nil, err
	}
	seasonInfo.EpName = ""
	return seasonInfo, nil
}

// Check interface is satisfied
var _ api.AniDataSearchGet = &Themoviedb{}
