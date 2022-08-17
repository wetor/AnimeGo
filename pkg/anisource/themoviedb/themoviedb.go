package themoviedb

import (
	"GoBangumi/pkg/anisource"
	"GoBangumi/pkg/request"
	"fmt"
	"net/url"
)

var (
	Host            = "https://api.themoviedb.org"
	DefaultSeason   = 1
	MatchSeasonDays = 90
)

var idApi = func(key string, query string) string {
	url_, _ := url.Parse(Host + "/3/discover/tv")
	q := url_.Query()
	q.Set("api_key", key)
	q.Set("language", "zh-CN")
	q.Set("timezone", "Asia/Shanghai")
	q.Set("with_genres", "16")
	q.Set("with_text_query", query)
	return url_.String() + "?" + q.Encode()
}

var infoApi = func(key string, id int) string {
	return fmt.Sprintf("%s/3/tv/%d?api_key=%s", Host, id, key)
}

type Themoviedb struct {
	Key string
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
