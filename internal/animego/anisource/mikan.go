package anisource

import (
	"github.com/google/wire"

	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
)

var MikanSet = wire.NewSet(
	NewMikanSource,
)

type Mikan struct {
	bangumiSource api.AniSource
	aniData       api.AniDataParse
}

func NewMikanSourceInterface(aniData api.AniDataParse, bangumiSource api.AniSource) *Mikan {
	return &Mikan{
		bangumiSource: bangumiSource,
		aniData:       aniData,
	}
}

func NewMikanSource(aniData *mikan.Mikan, bangumiSource *Bangumi) *Mikan {
	return &Mikan{
		bangumiSource: bangumiSource,
		aniData:       aniData,
	}
}

func (m Mikan) Parse(opts *models.AnimeParseOptions) (anime *models.AnimeEntity, err error) {
	var mikanEntity = &mikan.Entity{}
	var mikanUrl string
	switch input := opts.Input.(type) {
	case string:
		mikanUrl = input
	}

	// ------------------- 获取bangumiID -------------------
	if opts.AnimeParseOverride == nil || !opts.OverrideMikan() {
		log.Debugf("[AniSource] 解析Mikan，%s", mikanUrl)
		entity, err := m.aniData.ParseCache(mikanUrl)
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
		Input: &models.MikanEntity{
			MikanID:   mikanEntity.MikanID,
			BangumiID: mikanEntity.BangumiID,
		},
		AnimeParseOverride: opts.AnimeParseOverride,
	})
}
