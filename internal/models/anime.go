package models

import (
	"fmt"
	"strconv"

	"github.com/wetor/AnimeGo/pkg/xpath"
)

// AnimeEntity 动画信息结构体
//  必须要有的值
//    NameCN: 用于保存文件名，可用 Name 和 ID 替代
//    Season: 用于保存文件名
//    Ep: 用于保存文件名
//  可选值
//    ID: bangumi id，用于生成nfo文件
//    ThemoviedbID: themoviedb id，用于生成nfo文件
type AnimeEntity struct {
	ID           int    `json:"id"`            // bangumi id
	ThemoviedbID int    `json:"themoviedb_id"` // themoviedb ID
	MikanID      int    `json:"mikan_id"`      // [暂时无用] rss id
	Name         string `json:"name"`          // 名称，从bgm获取
	NameCN       string `json:"name_cn"`       // 中文名称，从bgm获取
	Season       int    `json:"season"`        // 当前季，从themoviedb获取
	Ep           int    `json:"ep"`            // 当前集，从下载文件名解析
	Eps          int    `json:"eps"`           // [暂时无用] 总集数，从bgm获取
	AirDate      string `json:"air_date"`      // 最初播放日期，从bgm获取
	*DownloadInfo
}

type DownloadInfo struct {
	Url  string `json:"url"`  // 当前集下载链接
	Hash string `json:"hash"` // 当前集Hash，唯一ID
}

func (b *AnimeEntity) FullName() string {
	if len(b.NameCN) == 0 {
		b.NameCN = b.Name
	}
	if len(b.NameCN) == 0 {
		b.NameCN = strconv.Itoa(b.ID)
	}
	return fmt.Sprintf("%s[第%d季][第%d集]", b.NameCN, b.Season, b.Ep)
}

func (b *AnimeEntity) FileName() string {
	return xpath.Join(fmt.Sprintf("S%02d", b.Season), fmt.Sprintf("E%03d", b.Ep))
}

func (b *AnimeEntity) DirName() string {
	if len(b.NameCN) == 0 {
		b.NameCN = b.Name
	}
	if len(b.NameCN) == 0 {
		b.NameCN = strconv.Itoa(b.ID)
	}
	return Filename(b.NameCN)
}

func (b *AnimeEntity) Meta() string {
	nfoTemplate := "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n<tvshow>\n"
	if b.ThemoviedbID != 0 {
		nfoTemplate += "  <tmdbid>{tmdbid}</tmdbid>\n"
	}
	nfoTemplate += "  <bangumiid>{bangumiid}</bangumiid>\n</tvshow>"
	return Format(nfoTemplate, map[string]any{
		"tmdbid":    b.ThemoviedbID,
		"bangumiid": b.ID,
	})
}
