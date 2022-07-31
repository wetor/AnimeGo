package bangumi

import (
	"GoBangumi/models"
	"GoBangumi/models/bgm/res"
	"GoBangumi/store"
	"GoBangumi/utils"
	"fmt"
	"go.uber.org/zap"
)

var BangumiInfoApi = func(id int) string {
	return fmt.Sprintf("%s/v0/subjects/%d", store.Config.Advanced.BangumiConf.Host, id)
}

var BangumiEpApi = func(id, ep int) string {
	// TODO: 支持根据上传日期，判断当前ep数
	conf := store.Config.Advanced.BangumiConf
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
	tmp := store.Cache.Get(models.BgmInfoBucket, bangumiID)
	if tmp != nil {
		if val, ok := tmp.(*models.Bangumi); ok {
			zap.S().Debugf("解析Bangumi，步骤1，缓存")
			return val
		}
	}
	zap.S().Debugf("解析Bangumi，步骤1，获取信息")
	url_ := BangumiInfoApi(bangumiID)
	resp := &res.SubjectV0{}
	status, err := utils.ApiGet(url_, resp, store.Config.Proxy())
	if err != nil {
		zap.S().Warn(err)
		return nil
	}
	if status != 200 {
		zap.S().Warn("解析bangumi失败，Status:", status)
		return nil
	}
	info = &models.Bangumi{
		ID:      int(resp.ID),
		NameCN:  resp.NameCN,
		Name:    resp.Name,
		AirDate: *resp.Date,
		Eps:     int(resp.Eps),
	}
	store.Cache.Put(models.BgmInfoBucket, bangumiID, info, store.Config.Advanced.CacheInfoExpire)
	return info
}

func (b *Bgm) parseBgm2(bangumiID, ep int, date string) (epInfo *models.BangumiEp) {
	cacheKey := fmt.Sprintf("%d_%d", bangumiID, ep)
	tmp := store.Cache.Get(models.BgmEpBucket, cacheKey)
	if tmp != nil {
		if val, ok := tmp.(*models.BangumiEp); ok {
			zap.S().Debugf("解析Bangumi，步骤2，缓存")
			return val
		}
	}
	zap.S().Debugf("解析Bangumi，步骤2，获取Ep信息")
	conf := store.Config.Advanced.BangumiConf
	url_ := BangumiEpApi(bangumiID, ep)
	resp := &res.Paged{
		Data: make([]*res.Episode, 0, conf.MatchEpRange*2+1),
	}
	status, err := utils.ApiGet(url_, resp, store.Config.Proxy())
	if err != nil {
		zap.S().Warn(err)
		return nil
	}
	if status != 200 {
		zap.S().Warn("解析bangumi ep失败，Status:", status)
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
		zap.S().Warn("解析bangumi ep失败，没有匹配到剧集信息")
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
	store.Cache.Put(models.BgmEpBucket, cacheKey, epInfo, conf.CacheEpExpire)
	return epInfo
}
