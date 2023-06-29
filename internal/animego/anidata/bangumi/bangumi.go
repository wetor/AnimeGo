package bangumi

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/pkg/log"
	mem "github.com/wetor/AnimeGo/pkg/memorizer"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/third_party/bangumi/res"
)

const (
	SubjectBucket = "bangumi_sub"
)

var (
	Host = func() string {
		if len(anidata.RedirectBangumi) > 0 {
			return anidata.RedirectBangumi
		}
		return "https://api.bgm.tv"
	}
	Bucket  = "bangumi"
	infoApi = func(id int) string {
		return fmt.Sprintf("%s/v0/subjects/%d", Host(), id)
	}
)

type Bangumi struct {
	cacheInit           bool
	cacheParseAnimeInfo mem.Func
}

func (a *Bangumi) Name() string {
	return "Bangumi"
}

func (a *Bangumi) RegisterCache() {
	a.cacheInit = true
	a.cacheParseAnimeInfo = mem.Memorized(Bucket, anidata.Cache, func(params *mem.Params, results *mem.Results) error {
		entity, err := a.parseAnimeInfo(params.Get("bangumiID").(int))
		if err != nil {
			return err
		}
		results.Set("entity", entity)
		return nil
	})
}

func (a *Bangumi) GetCache(bangumiID int, filters any) (entity any, err error) {
	if !a.cacheInit {
		a.RegisterCache()
	}
	e, err := a.loadAnimeInfo(bangumiID)
	if err == nil && e != nil {
		log.Infof("使用Bangumi本地缓， %d", bangumiID)
		return e, nil
	}

	results := mem.NewResults("entity", &Entity{})
	err = a.cacheParseAnimeInfo(mem.NewParams("bangumiID", bangumiID).
		TTL(anidata.CacheTime[Bucket]), results)
	if err != nil {
		return nil, errors.Wrap(err, "获取Bangumi信息失败")
	}
	entity = results.Get("entity").(*Entity)
	return entity, nil
}

// Get
//
//	@Description: 通过bangumiID和指定ep集数，获取番剧信息和剧集信息
//	@receiver Bangumi
//	@param bangumiID int
//	@param ep int
//	@return entity *Entity
func (a *Bangumi) Get(bangumiID int, filters any) (entity any, err error) {
	entity, err = a.parseAnimeInfo(bangumiID)
	if err != nil {
		return nil, errors.Wrap(err, "获取Bangumi信息失败")
	}
	return entity, err
}

// parseAnimeInfo
//
//	@Description: 解析番剧信息
//	@receiver Bangumi
//	@param bangumiID int
//	@return entity *Entity
func (a *Bangumi) parseAnimeInfo(bangumiID int) (entity *Entity, err error) {
	uri := infoApi(bangumiID)
	resp := res.SubjectV0{}
	err = request.Get(uri, &resp)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrRequest{Name: a.Name()})
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

func (a *Bangumi) loadAnimeInfo(bangumiID int) (entity *Entity, err error) {
	entity = &Entity{}
	anidata.BangumiCacheLock.Lock()
	defer anidata.BangumiCacheLock.Unlock()
	err = anidata.BangumiCache.Get(SubjectBucket, bangumiID, entity)
	if err != nil {
		// log.DebugErr(err)
		return nil, err
	}
	return entity, nil
}

// Check interface is satisfied
var _ api.AniDataGet = &Bangumi{}
