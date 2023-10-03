package models

import (
	"fmt"
	"strconv"
)

const (
	AnimeFlagNone          = 0
	AnimeFlagEpParseFailed = 1 << (iota - 1)
	AnimeFlagSeasonParseFailed
)

// AnimeEntity 动画信息结构体
//
//	必须要有的值
//	  NameCN: 用于保存文件名，可用 Name 和 ID 替代
//	  Season: 用于保存文件名
//	  Ep: 用于保存文件名
//	可选值
//	  ID: bangumi id，用于生成nfo文件
//	  ThemoviedbID: themoviedb id，用于生成nfo文件
type AnimeEntity struct {
	ID           int              `json:"id"`            // bangumi id
	ThemoviedbID int              `json:"themoviedb_id"` // themoviedb ID
	MikanID      int              `json:"mikan_id"`      // [暂时无用] rss id
	Name         string           `json:"name"`          // 名称，从bgm获取
	NameCN       string           `json:"name_cn"`       // 中文名称，从bgm获取
	Season       int              `json:"season"`        // 当前季，从themoviedb获取
	Eps          int              `json:"eps"`           // [暂时无用] 总集数，从bgm获取
	AirDate      string           `json:"air_date"`      // 最初播放日期，从bgm获取
	Ep           []*AnimeEpEntity `json:"ep"`
	Flag         int              `json:"flag"`
	Torrent      *AnimeTorrent    `json:"torrent"`
}

type AnimeEpType int8

const (
	AnimeEpUnknown AnimeEpType = iota
	AnimeEpNormal
	AnimeEpSpecial
)

type AnimeEpEntity struct {
	Type    AnimeEpType `json:"type"`     // ep类型。0:未知，1:正常剧集，2:SP
	Ep      int         `json:"ep"`       // 集数，Type=0时，不使用此参数
	Src     string      `json:"src"`      // 原文件名
	AirDate string      `json:"air_date"` // 首映日期
}

type AnimeTorrent struct {
	Hash string `json:"hash"`
	Url  string `json:"url"`
	File string `json:"file"`
}

func (a *AnimeEntity) Default() {
	if len(a.NameCN) == 0 {
		a.NameCN = a.Name
	}
	if len(a.NameCN) == 0 {
		a.NameCN = strconv.Itoa(a.ID)
	}
}

func (a *AnimeEntity) AnimeName() string {
	return a.NameCN
}

func (a *AnimeEntity) Hash() string {
	if a.Torrent != nil {
		return a.Torrent.Hash
	}
	return ""
}

func (a *AnimeEntity) EpKeys() []string {
	result := make([]string, len(a.Ep))
	for i := range a.Ep {
		result[i] = a.EpKey(i)
	}
	return result
}

func (a *AnimeEntity) EpKey(index int) string {
	switch a.Ep[index].Type {
	case AnimeEpUnknown:
		return fmt.Sprintf("key-%v-S%v-%v", a.NameCN, a.Season, a.Ep[index].Src)
	case AnimeEpSpecial:
		return fmt.Sprintf("key-%v-S%v-SP%v", a.NameCN, a.Season, a.Ep[index].Ep)
	}
	return a.EpKeyByEp(a.Ep[index].Ep)
}

func (a *AnimeEntity) EpKeyByEp(ep int) string {
	return fmt.Sprintf("key-%v-S%v-E%v", a.NameCN, a.Season, ep)
}

func (a *AnimeEntity) FullName() string {
	if a.Flag&AnimeFlagEpParseFailed > 0 {
		return fmt.Sprintf("%s[第%d季][第-集][%s]", a.NameCN, a.Season, a.Hash())
	}
	if len(a.Ep) == 1 {
		return fmt.Sprintf("%s[第%d季][第%d集]", a.NameCN, a.Season, a.Ep[0].Ep)
	}
	return fmt.Sprintf("%s[第%d季][%s集]", a.NameCN, a.Season, ToIntervals(a.Ep))

}

func (a *AnimeEntity) FileName(index int) string {
	return AnimeToFileName(a, index)
}

func (a *AnimeEntity) DirName() string {
	return FileName(a.NameCN)
}

func (a *AnimeEntity) FilePath() []string {
	return AnimeToFilePath(a)
}

func (a *AnimeEntity) FilePathSrc() []string {
	return AnimeToFilePathSrc(a)
}

func (a *AnimeEntity) Meta() string {
	return ToMetaData(a.ThemoviedbID, a.ID)
}
