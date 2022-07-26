package bangumi

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/modules/parser"
	"GoBangumi/store"
	"GoBangumi/utils"
	"fmt"
	"github.com/golang/glog"
	"net/url"
)

var ThemoviedbIdApi = func(query string) string {
	url_, _ := url.Parse(config.Advanced().Themoviedb().Host + "/3/discover/tv")
	q := url_.Query()
	q.Set("api_key", config.KeyTmdb())
	q.Set("language", "zh-CN")
	q.Set("timezone", "Asia/Shanghai")
	q.Set("with_genres", "16")
	q.Set("with_text_query", query)
	return url_.String() + "?" + q.Encode()
}
var ThemoviedbInfoApi = func(id int) string {
	return fmt.Sprintf("%s/3/tv/%d?api_key=%s", config.Advanced().Themoviedb().Host, id, config.KeyTmdb())
}

type Themoviedb struct {
}

func NewThemoviedb() Bangumi {
	return &Themoviedb{}
}
func (b *Themoviedb) Parse(opt *models.BangumiParseOptions) *models.Bangumi {

	id := b.parseThemoviedb1(opt.Name)
	season := b.parseThemoviedb2(id, opt.Date)
	return &models.Bangumi{
		BangumiSeason: season,
	}
}

// parseThemoviedb1
//  @Cache name_tmdb name
//
func (b *Themoviedb) parseThemoviedb1(name string) (tmdbID int) {
	// 通过name查询缓存中的tmdbID
	tmp := store.Cache.Get("name_tmdb", name)
	if tmp != nil {
		if val, ok := tmp.(int); ok {
			glog.V(5).Infof("解析Themoviedb，步骤1，缓存\n")
			return val
		}
	}
	glog.V(5).Infof("解析Themoviedb，步骤1，获取tmdb ID\n")
	resp := &models.ThemoviedbIdResponse{}
	nameParser := parser.NewBangumiName()
	step := 0
	for {
		status, err := utils.ApiGet(ThemoviedbIdApi(name), resp, config.Proxy())
		if err != nil {
			glog.Errorln(err)
			return 0
		}
		if status != 200 && resp == nil {
			glog.Errorf("Themoviedb查找错误，状态码：%d，%s\n", status, name)
			return 0
		}
		if resp.TotalResults == 0 {
			glog.Errorln("Themoviedb中未找到番剧：" + name)
			result := nameParser.Parse(&models.ParseNameOptions{
				Name:      name,
				StartStep: step,
			})
			if result == nil {
				return 0
			}
			step = result.NextStep
			name = result.Name
			glog.Errorln("Themoviedb重新查找番剧名：" + name)
			continue

		} else {
			tmdbID = resp.Result[0].ID
			break
		}
	}
	store.Cache.Put("name_tmdb", name, tmdbID, config.Advanced().Themoviedb().CacheIdExpire)
	return tmdbID
}
func (b *Themoviedb) parseThemoviedb2(id int, date string) (season *models.BangumiSeason) {
	cacheKey := fmt.Sprintf("%d_%s", id, date)
	tmp := store.Cache.Get("tmdb_season", cacheKey)
	if tmp != nil {
		if val, ok := tmp.(*models.BangumiSeason); ok {
			glog.V(5).Infof("解析Themoviedb，步骤2，缓存\n")
			return val
		}
	}
	glog.V(5).Infof("解析Themoviedb，步骤2，获取信息\n")
	resp := &models.ThemoviedbResponse{}
	status, err := utils.ApiGet(ThemoviedbInfoApi(id), resp, config.Proxy())
	if err != nil {
		glog.Errorln(err)
		return nil
	}
	if status != 200 {
		glog.Errorln("Status:", status)
		return nil
	}
	season = &models.BangumiSeason{
		Season: 1,
	}
	if resp.Seasons == nil || len(resp.Seasons) == 0 {
		return season
	}
	season.Season = resp.Seasons[0].SeasonNumber
	min := 36500
	for _, r := range resp.Seasons {
		if r.SeasonNumber == 0 || r.Name == "Specials" {
			continue
		}
		if s := utils.StrTimeSubAbs(r.AirDate, date); s < min {
			min = s
			season.Season = r.SeasonNumber
		}
	}
	conf := config.Advanced().Themoviedb()
	if min > conf.MatchSeasonDays {
		glog.Errorln("Themoviedb匹配Seasons失败，可能此番剧未开播")
		return nil
	}
	if season.Season == 0 {
		glog.Errorln("Themoviedb匹配Seasons失败")
		return nil
	}
	store.Cache.Put("tmdb_season", cacheKey, season, conf.CacheSeasonExpire)
	return season
}
