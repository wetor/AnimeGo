package bangumi

import (
	"GoBangumi/pkg/anisource"
	"GoBangumi/pkg/request"
	"GoBangumi/third_party/bangumi/res"
	"fmt"
)

var (
	Host         = "https://api.bgm.tv"
	MatchEpRange = 10
)

var infoApi = func(id int) string {
	return fmt.Sprintf("%s/v0/subjects/%d", Host, id)
}

var epInfoApi = func(id, ep, eps int) string {
	rang := 2
	if eps > 15 {
		rang = 4
	} else if eps > 40 {
		rang = 6
	} else if eps > 80 {
		rang = 10
	} else if eps > 150 {
		rang = 20
	}
	offset := ep - 1 - rang
	if offset < 0 {
		offset = 0
	}
	limit := rang*2 + 1

	epType := 0 // 仅番剧本体
	return fmt.Sprintf("%s/v0/episodes?subject_id=%d&type=%d&limit=%d&offset=%d", Host, id, epType, limit, offset)
}

type Bangumi struct {
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
