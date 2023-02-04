package bangumi

import (
	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	mem "github.com/wetor/AnimeGo/pkg/memorizer"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/third_party/bangumi/res"
)

const (
	SubjectBucket = "bangumi_sub"
)

var (
	Host   = "https://api.bgm.tv"
	Bucket = "bangumi"
)

type Bangumi struct {
	cacheInit           bool
	cacheParseAnimeInfo mem.Func
}

func (b *Bangumi) RegisterCache() {
	if anidata.Cache == nil {
		errors.NewAniError("需要先调用anidata.Init初始化缓存").TryPanic()
	}
	b.cacheInit = true
	b.cacheParseAnimeInfo = mem.Memorized(Bucket, anidata.Cache, func(params *mem.Params, results *mem.Results) error {
		entity := b.parseAnimeInfo(params.Get("bangumiID").(int))
		results.Set("entity", entity)
		return nil
	})
}

func (b Bangumi) ParseCache(bangumiID int) (entity *Entity) {
	if !b.cacheInit {
		b.RegisterCache()
	}

	if e, err := b.loadAnimeInfo(bangumiID); err == nil {
		if e != nil {
			log.Debugf("使用Bangumi Archive，%d", bangumiID)
			return e
		}
	}

	results := mem.NewResults("entity", &Entity{})
	err := b.cacheParseAnimeInfo(mem.NewParams("bangumiID", bangumiID).
		TTL(anidata.CacheTime[Bucket]), results)
	errors.NewAniErrorD(err).TryPanic()
	entity = results.Get("entity").(*Entity)

	return entity
}

// Parse
//  @Description: 通过bangumiID和指定ep集数，获取番剧信息和剧集信息
//  @receiver Bangumi
//  @param bangumiID int
//  @param ep int
//  @return entity *Entity
//
func (b Bangumi) Parse(bangumiID int) (entity *Entity) {
	entity = b.parseAnimeInfo(bangumiID)
	return entity
}

// parseAnimeInfo
//  @Description: 解析番剧信息
//  @receiver Bangumi
//  @param bangumiID int
//  @return entity *Entity
//
func (b Bangumi) parseAnimeInfo(bangumiID int) (entity *Entity) {
	uri := infoApi(bangumiID)
	resp := res.SubjectV0{}

	err := request.Get(uri, &resp)
	errors.NewAniErrorD(err).TryPanic()

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
	return entity
}

func (b Bangumi) loadAnimeInfo(bangumiID int) (entity *Entity, err error) {
	entity = &Entity{}
	anidata.BangumiCacheLock.Lock()
	err = anidata.BangumiCache.Get(SubjectBucket, bangumiID, entity)
	anidata.BangumiCacheLock.Unlock()
	return entity, err
}
