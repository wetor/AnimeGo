package bangumi

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/modules/parser"
	"GoBangumi/utils"
	"fmt"
	"github.com/golang/glog"
	"net/url"
)

const (
	ThemoviedbBaseApi = "https://api.themoviedb.org"
)

var ThemoviedbIdApi = func(query string) string {
	url_, _ := url.Parse(ThemoviedbBaseApi + "/3/discover/tv")
	q := url_.Query()
	q.Set("api_key", config.TMDB())
	q.Set("language", "zh-CN")
	q.Set("timezone", "Asia/Shanghai")
	q.Set("with_genres", "16")
	q.Set("with_text_query", query)
	return url_.String() + "?" + q.Encode()
}
var ThemoviedbInfoApi = func(id int) string {
	return fmt.Sprintf("%s/3/tv/%d?api_key=%s", ThemoviedbBaseApi, id, config.TMDB())
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

func (b *Themoviedb) parseThemoviedb1(name string) int {
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
			result := nameParser.ParseBangumiName(&models.ParseBangumiNameOptions{
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
			return resp.Result[0].ID
		}
	}

}
func (b *Themoviedb) parseThemoviedb2(id int, date string) *models.BangumiSeason {

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
	season := &models.BangumiSeason{
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
	if min > 90 {
		glog.Errorln("Themoviedb匹配Seasons失败，可能此番剧未开播")
		return nil
	}
	if season.Season == 0 {
		glog.Errorln("Themoviedb匹配Seasons失败")
		return nil
	}
	return season
}
