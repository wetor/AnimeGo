package bangumi

import (
	"fmt"

	"github.com/google/wire"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/pkg/request"
	"github.com/wetor/AnimeGo/pkg/log"
	mem "github.com/wetor/AnimeGo/pkg/memorizer"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/third_party/bangumi/model"
	"github.com/wetor/AnimeGo/third_party/bangumi/res"
)

var (
	Host = func(host string) string {
		if len(host) > 0 {
			return host
		}
		return constant.BangumiDefaultHost
	}
	infoApi = func(host string, id int) string {
		return fmt.Sprintf("%s/v0/subjects/%d", Host(host), id)
	}
	searchApi = func(host string, limit, offset int) string {
		return fmt.Sprintf("%s/v0/search/subjects?limit=%d&offset=%d", Host(host), limit, offset)
	}
)

var Set = wire.NewSet(
	NewBangumi,
)

type Bangumi struct {
	cacheInit           bool
	cacheParseAnimeInfo mem.Func

	*Options
}

func NewBangumi(opts *Options) *Bangumi {
	return &Bangumi{
		Options: opts,
	}
}

func (a *Bangumi) Name() string {
	return "Bangumi"
}

func (a *Bangumi) RegisterCache() {
	a.cacheInit = true
	a.cacheParseAnimeInfo = mem.Memorized(constant.BangumiBucket, a.Cache, func(params *mem.Params, results *mem.Results) error {
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
		TTL(a.CacheTime), results)
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

func (a *Bangumi) Search(name string, filters any) (int, error) {
	entity, err := a.searchAnimeInfo(name)
	if err != nil {
		return 0, errors.Wrap(err, "查询BangumiID失败")
	}
	return entity.ID, nil
}

func (a *Bangumi) SearchCache(name string, filters any) (int, error) {

	return 0, nil
}

func (a *Bangumi) searchAnimeInfo(name string) (entity *Entity, err error) {
	uri := searchApi(a.Host, 10, 0)
	resp := res.SearchPaged{}
	result, err := utils.RemoveNameSuffix(name, func(innerName string) (any, error) {
		req := res.Req{
			Keyword: innerName,
			Sort:    "match",
			Filter: res.ReqFilter{
				Type: []model.SubjectType{model.SubjectAnime},
			},
		}

		err := request.Post(uri, req, &resp)
		if err != nil {
			log.DebugErr(err)
			//return 0, errors.WithStack(&exceptions.ErrRequest{Name: a.Name()})
		}

		if resp.Total == 1 {
			return resp.Data[0], nil
		} else if resp.Total > 1 {
			// 筛选与original name完全相同的番剧
			for _, result := range resp.Data {
				if result.NameCN == name || result.Name == name {
					return result, nil
				}
			}

			// 按照相似度排序筛选
			var temp *res.ReponseSubject
			maxSimilar := float64(0)
			for _, result := range resp.Data {
				// 分别对比中文和原名，选取最相似的
				similarCN := utils.SimilarText(result.NameCN, name)
				similar := utils.SimilarText(result.Name, name)
				if similarCN > similar {
					similar = similarCN
				}
				if similar > maxSimilar {
					maxSimilar = similar
					temp = result
				}
			}
			if maxSimilar >= constant.BangumiMinSimilar {
				return temp, nil
			}
			err = errors.WithStack(&exceptions.ErrAniDataSearch{AniData: a})
			log.DebugErr(err)
			return nil, err
		} else {
			// 未找到结果
			return nil, nil
		}
	})
	if err != nil {
		return nil, err
	}
	sub := result.(*res.ReponseSubject)
	return &Entity{
		ID:      int(sub.ID),
		NameCN:  sub.NameCN,
		Name:    sub.Name,
		AirDate: sub.Date,
	}, nil

}

// parseAnimeInfo
//
//	@Description: 解析番剧信息
//	@receiver Bangumi
//	@param bangumiID int
//	@return entity *Entity
func (a *Bangumi) parseAnimeInfo(bangumiID int) (entity *Entity, err error) {
	uri := infoApi(a.Host, bangumiID)
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
	a.BangumiCacheLock.Lock()
	defer a.BangumiCacheLock.Unlock()
	err = a.BangumiCache.Get(constant.BangumiSubjectBucket, bangumiID, entity)
	if err != nil {
		// log.DebugErr(err)
		return nil, err
	}
	if entity.Eps == 0 || len(entity.AirDate) == 0 {
		return nil, errors.WithStack(&exceptions.ErrBangumiCacheNotFound{BangumiID: bangumiID})
	}
	return entity, nil
}

// Check interface is satisfied
var _ api.AniDataGet = &Bangumi{}
