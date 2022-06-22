package bangumi

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/models/bgm/res"
	"GoBangumi/utils"
	"fmt"
	"github.com/golang/glog"
)

const (
	BangumiBaseApi = "https://api.bgm.tv" // Bangumi 域名
)

var BangumiInfoApi = func(id int) string {
	return fmt.Sprintf("%s/v0/subjects/%d", BangumiBaseApi, id)
}
var BangumiEpApi = func(id, ep int) string {
	// TODO: 支持根据上传日期，判断当前ep数
	offset := ep - 2 // 缓存当前ep的前一集
	if offset < 0 {
		offset = 0
	}
	limit := 3  // 共缓存三集
	epType := 0 // 仅番剧本体
	return fmt.Sprintf("%s/v0/episodes?subject_id=%d&type=%d&limit=%d&offset=%d", BangumiBaseApi, id, epType, limit, offset)
}

type Bgm struct {
}

func NewBgm() Bangumi {
	return &Bgm{}
}
func (b *Bgm) Parse(opt *models.BangumiParseOptions) *models.Bangumi {
	resp := b.parseBgm1(opt.ID)

	info := &models.Bangumi{
		ID:      int(resp.ID),
		NameCN:  resp.NameCN,
		Name:    resp.Name,
		AirDate: *resp.Date,
	}
	ep := b.parseBgm2(opt.ID, opt.Ep, opt.Date)
	if ep != nil {
		info.BangumiEp = &models.BangumiEp{
			Ep:       opt.Ep,
			Date:     ep.Airdate,
			Duration: ep.Duration,
			EpDesc:   ep.Description,
			EpName:   ep.Name,
			EpNameCN: ep.NameCN,
			EpID:     int(ep.ID),
			Eps:      int(resp.Eps),
		}
	} else {
		info.BangumiEp = &models.BangumiEp{
			Ep:  opt.Ep,
			Eps: int(resp.Eps),
		}
	}

	return info
}

func (b *Bgm) parseBgm1(bangumiID int) *res.SubjectV0 {
	url_ := BangumiInfoApi(bangumiID)
	resp := &res.SubjectV0{}
	status, err := utils.ApiGet(url_, resp, config.Proxy())
	if err != nil {
		glog.Errorln(err)
		return nil
	}
	if status != 200 {
		glog.Errorln("解析bangumi失败，Status:", status)
		return nil
	}
	return resp
}
func (b *Bgm) parseBgm2(bangumiID, ep int, date string) *res.Episode {
	url_ := BangumiEpApi(bangumiID, ep)
	resp := &res.Paged{
		Data: make([]*res.Episode, 0, 3),
	}
	status, err := utils.ApiGet(url_, resp, config.Proxy())
	if err != nil {
		glog.Errorln(err)
		return nil
	}
	if status != 200 {
		glog.Errorln("解析bangumi ep失败，Status:", status)
		return nil
	}
	// TODO: 根据ep、date是否为空进行不同规则的判断
	for _, e := range resp.Data {
		s := utils.StrTimeSubAbs(date, e.Airdate)
		if ep == int(e.Ep) && s <= 30 {
			return e
		}
	}
	glog.Errorln("解析bangumi ep失败，没有匹配到剧集信息")
	return nil
}
