package mikan

import (
	"github.com/wetor/AnimeGo/internal/animego/anidata/mikan"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
)

type Mikan struct {
	bangumiSource api.AniSource
}

func NewMikanSource(bangumiSource api.AniSource) api.AniSource {
	return &Mikan{
		bangumiSource: bangumiSource,
	}
}

func (m Mikan) Parse(opts *models.AnimeParseOptions) (anime *models.AnimeEntity, err error) {
	var mikanEntity = &mikan.Entity{}

	// ------------------- 获取bangumiID -------------------
	if opts.AnimeParseOverride == nil || !opts.OverrideMikan() {
		log.Debugf("[AniSource] 解析Mikan，%s", opts.Input)
		entity, err := anisource.Mikan().ParseCache(opts.Input)
		if err != nil {
			return nil, err
		}
		mikanEntity = entity.(*mikan.Entity)
	} else {
		mikanEntity.MikanID = opts.AnimeParseOverride.MikanID
		mikanEntity.BangumiID = opts.AnimeParseOverride.BangumiID
	}
	// ------------------- 通过bangumiID获取信息 -------------------
	return m.bangumiSource.Parse(&models.AnimeParseOptions{
		Input: models.MikanEntity{
			MikanID:   mikanEntity.MikanID,
			BangumiID: mikanEntity.BangumiID,
		},
		AnimeParseOverride: opts.AnimeParseOverride,
	})
}
