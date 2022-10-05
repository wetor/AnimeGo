package bangumi

import (
	"AnimeGo/pkg/anisource"
	mem "AnimeGo/pkg/memorizer"
	"AnimeGo/pkg/request"
	"AnimeGo/third_party/bangumi/res"
	"encoding/gob"
)

var (
	Host               = "https://api.bgm.tv"
	Bucket             = "bangumi"
	MatchEpRange       = 10
	CacheSecond  int64 = 3 * 24 * 60 * 60
)

type Bangumi struct {
	cacheInit             bool
	cacheParseAnimeInfo   mem.Func
	cacheParseAnimeEpInfo mem.Func
}

func (b *Bangumi) RegisterCache() {
	if anisource.Cache == nil {
		panic("需要先调用anisource.Init初始化缓存")
	}
	b.cacheInit = true
	b.cacheParseAnimeInfo = mem.Memorized(Bucket, anisource.Cache, func(params *mem.Params, results *mem.Results) error {
		entity, err := b.parseAnimeInfo(params.Get("bangumiID").(int))
		if err != nil {
			return err
		}
		results.Set("entity", entity)
		return nil
	})

	b.cacheParseAnimeEpInfo = mem.Memorized(Bucket, anisource.Cache, func(params *mem.Params, results *mem.Results) error {
		epInfo, err := b.parseAnimeEpInfo(
			params.Get("bangumiID").(int),
			params.Get("ep").(int),
			params.Get("eps").(int),
		)
		if err != nil {
			return err
		}
		results.Set("epInfo", epInfo)
		return nil
	})
}

func (b Bangumi) ParseCache(bangumiID, ep int) (entity *Entity, epInfo *Ep, err error) {
	if !b.cacheInit {
		b.RegisterCache()
	}
	results := mem.NewResults("entity", &Entity{}, "epInfo", &Ep{})

	err = b.cacheParseAnimeInfo(mem.NewParams("bangumiID", bangumiID).TTL(CacheSecond), results)
	if err != nil {
		return nil, nil, err
	}
	entity = results.Get("entity").(*Entity)
	err = b.cacheParseAnimeEpInfo(
		mem.NewParams("bangumiID", bangumiID, "ep", ep, "eps", entity.Eps).TTL(CacheSecond), results)
	if err != nil {
		return nil, nil, err
	}
	epInfo = results.Get("epInfo").(*Ep)
	return entity, epInfo, nil
}

// Parse
//  @Description: 通过bangumiID和指定ep集数，获取番剧信息和剧集信息
//  @receiver Bangumi
//  @param bangumiID int
//  @param ep int
//  @return entity *Entity
//  @return epInfo *Ep
//  @return err error
//
func (b Bangumi) Parse(bangumiID, ep int) (entity *Entity, epInfo *Ep, err error) {
	entity, err = b.parseAnimeInfo(bangumiID)
	if err != nil {
		return nil, nil, err
	}
	epInfo, err = b.parseAnimeEpInfo(bangumiID, ep, entity.Eps)
	if err != nil {
		return nil, nil, err
	}
	return entity, epInfo, nil
}

// parseAnimeInfo
//  @Description: 解析番剧信息
//  @receiver Bangumi
//  @param bangumiID int
//  @return entity *Entity
//  @return err error
//
func (b Bangumi) parseAnimeInfo(bangumiID int) (entity *Entity, err error) {
	uri := infoApi(bangumiID)
	resp := res.SubjectV0{}

	err = request.Get(&request.Param{
		Uri:      uri,
		Proxy:    anisource.Proxy,
		BindJson: &resp,
		Retry:    anisource.Retry,
		Timeout:  anisource.Timeout,
	})
	if err != nil {
		return nil, err
	}
	entity = &Entity{
		ID:     int(resp.ID),
		NameCN: resp.NameCN,
		Name:   resp.Name,
	}
	if resp.Eps != 0 {
		entity.Eps = int(resp.Eps)
	} else {
		entity.Eps = int(resp.TotalEpisodes)
	}
	if resp.Date != nil {
		entity.AirDate = *resp.Date
	}
	return entity, nil
}

// parseBnagumiEpInfo
//  @Description: 解析番剧ep信息，暂时无用
//  @receiver Bangumi
//  @param bangumiID int
//  @param ep int
//  @param eps int 总集数，用于计算筛选范围，减少遍历范围
//  @return epInfo *Ep
//  @return err error
//
func (b Bangumi) parseAnimeEpInfo(bangumiID, ep, eps int) (epInfo *Ep, err error) {
	uri := epInfoApi(bangumiID, ep, eps)
	resp := &res.Paged{
		Data: make([]*res.Episode, 0, MatchEpRange*2+1),
	}
	err = request.Get(&request.Param{
		Uri:      uri,
		Proxy:    anisource.Proxy,
		BindJson: &resp,
		Retry:    anisource.Retry,
		Timeout:  anisource.Timeout,
	})
	if err != nil {
		return nil, err
	}
	var respEp *res.Episode = nil
	for _, e := range resp.Data {
		if ep == int(e.Ep) {
			respEp = e
			break
		}
	}
	if respEp == nil {
		return nil, NotMatchEpErr
	}
	epInfo = &Ep{
		Ep:       int(respEp.Ep),
		Date:     respEp.Airdate,
		Duration: respEp.Duration,
		EpDesc:   respEp.Description,
		EpName:   respEp.Name,
		EpNameCN: respEp.NameCN,
		EpID:     int(respEp.ID),
	}
	return epInfo, nil
}

func init() {
	gob.Register(&Entity{})
	gob.Register(&Ep{})
}
