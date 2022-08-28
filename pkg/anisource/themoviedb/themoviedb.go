package themoviedb

import (
	"AnimeGo/pkg/anisource"
	mem "AnimeGo/pkg/memorizer"
	"AnimeGo/pkg/request"
)

var (
	Host                  = "https://api.themoviedb.org"
	Bucket                = "themoviedb"
	DefaultSeason         = 1
	MatchSeasonDays       = 90
	CacheSecond     int64 = 7 * 24 * 60 * 60
)

type Themoviedb struct {
	Key                    string
	cacheInit              bool
	cacheParseThemoviedbID mem.Func
	cacheParseAnimeSeason  mem.Func
}

func (t *Themoviedb) RegisterCache() {
	if anisource.Cache == nil {
		panic("需要先调用anisource.Init初始化缓存")
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
		return 0, DefaultSeason, err
	}
	tmdbID = results.Get("tmdbID").(int)
	err = t.cacheParseAnimeSeason(mem.NewParams("tmdbID", tmdbID, "airDate", airDate).TTL(CacheSecond), results)
	if err != nil {
		return tmdbID, DefaultSeason, err
	}
	season = results.Get("season").(int)
	return tmdbID, season, nil
}

func (t Themoviedb) Parse(name, airDate string) (tmdbID int, season int, err error) {
	tmdbID, err = t.parseThemoviedbID(name)
	if err != nil {
		return 0, DefaultSeason, err
	}
	season, err = t.parseAnimeSeason(tmdbID, airDate)
	if err != nil {
		return tmdbID, DefaultSeason, err
	}
	return tmdbID, season, nil
}

func (t Themoviedb) parseThemoviedbID(name string) (tmdbID int, err error) {
	resp := FindResponse{}
	step := 0
	for step >= 0 {
		err = request.Get(&request.Param{
			Uri:      idApi(t.Key, name),
			Proxy:    anisource.Proxy,
			BindJson: &resp,
		})
		if err != nil {
			return 0, err
		}
		if resp.TotalResults != 0 {
			tmdbID = resp.Result[0].ID
			for _, result := range resp.Result {
				if result.OriginalName == name {
					tmdbID = result.ID
					break
				}
			}
			return tmdbID, nil
		} else {
			result, nextStep, err := RemoveNameSuffix(name, step)
			if err != nil {
				return 0, err
			}
			step = nextStep
			name = result
			continue
		}
	}
	return 0, NotFoundAnimeNameErr
}

func (t Themoviedb) parseAnimeSeason(tmdbID int, airDate string) (season int, err error) {
	resp := InfoResponse{}
	err = request.Get(&request.Param{
		Uri:      infoApi(t.Key, tmdbID),
		Proxy:    anisource.Proxy,
		BindJson: &resp,
	})
	if err != nil {
		return DefaultSeason, err
	}
	if resp.Seasons == nil || len(resp.Seasons) == 0 {
		return DefaultSeason, NotMatchSeasonErr
	}
	season = resp.Seasons[0].SeasonNumber
	min := 36500
	for _, r := range resp.Seasons {
		if r.SeasonNumber == 0 || r.Name == "Specials" {
			continue
		}
		// TODO: 待优化，通过比较此季度番剧的初放松日期，筛选差值最小的季
		if s := StrTimeSubAbs(r.AirDate, airDate); s < min {
			min = s
			season = r.SeasonNumber
		}
	}
	if min > MatchSeasonDays {
		return DefaultSeason, NotMatchSeasonErr
	}
	return season, nil
}
