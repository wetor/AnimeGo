package bangumi

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/models/bgm/res"
	"GoBangumi/modules/cache"
	"GoBangumi/store"
	"GoBangumi/utils"
	"fmt"
	"github.com/golang/glog"
)

var BangumiInfoApi = func(id int) string {
	return fmt.Sprintf("%s/v0/subjects/%d", config.Advanced().Bangumi().Host, id)
}
var BangumiEpApi = func(id, ep int) string {
	// TODO: 支持根据上传日期，判断当前ep数
	conf := config.Advanced().Bangumi()
	offset := ep - 1 - conf.MatchEpRange // 缓存当前ep的前一集
	if offset < 0 {
		offset = 0
	}
	limit := conf.MatchEpRange*2 + 1 // 共缓存三集
	epType := 0                      // 仅番剧本体
	return fmt.Sprintf("%s/v0/episodes?subject_id=%d&type=%d&limit=%d&offset=%d", conf.Host, id, epType, limit, offset)
}

type Bgm struct {
}

func NewBgm() Bangumi {
	return &Bgm{}
}
func (b *Bgm) Parse(opt *models.BangumiParseOptions) *models.Bangumi {
	info := b.parseBgm1(opt.ID)
	if info == nil {
		return nil
	}
	info.BangumiEp = b.parseBgm2(opt.ID, opt.Ep, opt.Date)
	if info.BangumiEp == nil {
		info.BangumiEp = &models.BangumiEp{
			Ep: opt.Ep,
		}
	}
	return info
}

func (b *Bgm) parseBgm1(bangumiID int) (info *models.Bangumi) {
	tmp := store.Cache.Get(cache.BgmInfoBucket, bangumiID)
	if tmp != nil {
		if val, ok := tmp.(*models.Bangumi); ok {
			glog.V(5).Infof("解析Bangumi，步骤1，缓存\n")
			return val
		}
	}
	glog.V(5).Infof("解析Bangumi，步骤1，获取信息\n")
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
	info = &models.Bangumi{
		ID:      int(resp.ID),
		NameCN:  resp.NameCN,
		Name:    resp.Name,
		AirDate: *resp.Date,
		Eps:     int(resp.Eps),
	}
	store.Cache.Put(cache.BgmInfoBucket, bangumiID, info, config.Advanced().Bangumi().CacheInfoExpire)
	return info
}
func (b *Bgm) parseBgm2(bangumiID, ep int, date string) (epInfo *models.BangumiEp) {
	cacheKey := fmt.Sprintf("%d_%d", bangumiID, ep)
	tmp := store.Cache.Get(cache.BgmEpBucket, cacheKey)
	if tmp != nil {
		if val, ok := tmp.(*models.BangumiEp); ok {
			glog.V(5).Infof("解析Bangumi，步骤2，缓存\n")
			return val
		}
	}
	glog.V(5).Infof("解析Bangumi，步骤2，获取Ep信息\n")
	conf := config.Advanced().Bangumi()
	url_ := BangumiEpApi(bangumiID, ep)
	resp := &res.Paged{
		Data: make([]*res.Episode, 0, conf.MatchEpRange*2+1),
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
	var respEp *res.Episode = nil
	for _, e := range resp.Data {
		s := utils.StrTimeSubAbs(date, e.Airdate)
		if ep == int(e.Ep) && s <= conf.MatchEpDays {
			respEp = e
			break
		}
	}
	if respEp == nil {
		glog.Errorln("解析bangumi ep失败，没有匹配到剧集信息")
		return nil
	}
	epInfo = &models.BangumiEp{
		Ep:       int(respEp.Ep),
		Date:     respEp.Airdate,
		Duration: respEp.Duration,
		EpDesc:   respEp.Description,
		EpName:   respEp.Name,
		EpNameCN: respEp.NameCN,
		EpID:     int(respEp.ID),
	}
	store.Cache.Put(cache.BgmEpBucket, cacheKey, epInfo, conf.CacheEpExpire)
	return epInfo
}
