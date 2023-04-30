package models

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/wetor/AnimeGo/pkg/xpath"
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
	Torrent      *AnimeTorrent    `json:"torrent"`
}

type AnimeEpEntity struct {
	Type    int    `json:"type"` // ep类型。0:正常剧集，1:SP
	Ep      int    `json:"ep"`
	Src     string `json:"src"`
	AirDate string `json:"air_date"`
}

type AnimeTorrent struct {
	Hash string `json:"hash"`
	Url  string `json:"url"`
}

func (b *AnimeEntity) Default() {
	if len(b.NameCN) == 0 {
		b.NameCN = b.Name
	}
	if len(b.NameCN) == 0 {
		b.NameCN = strconv.Itoa(b.ID)
	}
}

// getIntervals 将整数序列转换为区间形式
func (b *AnimeEntity) getIntervals() string {
	// 初始化区间的起始值和结束值
	start, end := b.Ep[0].Ep, b.Ep[0].Ep
	// 初始化区间的字符串表示
	var intervalSlice []string
	for i := 1; i < len(b.Ep); i++ {
		num := b.Ep[i].Ep
		// 如果当前数与上一个数相差1，表示当前数与上一个数在同一个区间内
		if num-b.Ep[i-1].Ep == 1 {
			end = num // 更新区间的结束值
		} else {
			// 否则，表示当前数与上一个数不在同一个区间内，需要将上一个区间输出，并重新开始一个新的区间
			if start == end {
				intervalSlice = append(intervalSlice, strconv.Itoa(start))
			} else {
				intervalSlice = append(intervalSlice, fmt.Sprintf("%d-%d", start, end))
			}
			start, end = num, num // 更新区间的起始值和结束值
		}
	}
	// 将最后一个区间输出
	if start == end {
		intervalSlice = append(intervalSlice, strconv.Itoa(start))
	} else {
		intervalSlice = append(intervalSlice, fmt.Sprintf("%d-%d", start, end))
	}
	// 对区间切片进行排序
	sort.Strings(intervalSlice)
	// 将排序后的区间切片合并成字符串
	intervals := strings.Join(intervalSlice, ",")
	return intervals
}

func (b *AnimeEntity) FullName() string {
	if len(b.Ep) == 1 {
		return fmt.Sprintf("%s[第%d季][第%d集]", b.NameCN, b.Season, b.Ep[0].Ep)
	}
	return fmt.Sprintf("%s[第%d季][%s集]", b.NameCN, b.Season, b.getIntervals())

}

func (b *AnimeEntity) FileName(index int) string {
	return xpath.Join(fmt.Sprintf("S%02d", b.Season), fmt.Sprintf("E%03d", b.Ep[index].Ep))
}

func (b *AnimeEntity) DirName() string {
	return Filename(b.NameCN)
}

func (b *AnimeEntity) FilePath() []string {
	filePath := make([]string, 0, len(b.Ep))
	for i, ep := range b.Ep {
		filePath = append(filePath, xpath.Join(b.DirName(), b.FileName(i)+xpath.Ext(ep.Src)))
	}
	return filePath
}

func (b *AnimeEntity) FilePathSrc() []string {
	filePath := make([]string, 0, len(b.Ep))
	for _, ep := range b.Ep {
		filePath = append(filePath, ep.Src)
	}
	return filePath
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
