package models

import (
	"fmt"
	"github.com/wetor/AnimeGo/internal/utils"
	"path"
	"strconv"
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
	EpID         int    `json:"ep_id"`         // [暂时无用] 当前集bangumi ep id
	Eps          int    `json:"eps"`           // [暂时无用] 总集数，从bgm获取
	AirDate      string `json:"air_date"`      // [暂时无用] 最初播放日期，从bgm获取
	Date         string `json:"date"`          // [暂时无用] 当前集播放日期，从bgm获取
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
	return path.Join(fmt.Sprintf("S%02d", b.Season), fmt.Sprintf("E%03d", b.Ep))
}

func (b *AnimeEntity) DirName() string {
	if len(b.NameCN) == 0 {
		b.NameCN = b.Name
	}
	if len(b.NameCN) == 0 {
		b.NameCN = strconv.Itoa(b.ID)
	}
	return utils.Filename(b.NameCN)
}

func (b *AnimeEntity) Meta() string {
	nfoTemplate := "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n<tvshow>\n"
	if b.ThemoviedbID != 0 {
		nfoTemplate += "  <tmdbid>{tmdbid}</tmdbid>\n"
	}
	nfoTemplate += "  <bangumiid>{bangumiid}</bangumiid>\n</tvshow>"
	return utils.Format(nfoTemplate, utils.FormatMap{
		"tmdbid":    b.ThemoviedbID,
		"bangumiid": b.ID,
	})
}

type ThemoviedbIdResponse struct {
	Page         int `json:"page"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
	Result       []*struct {
		BackdropPath string `json:"backdrop_path"`
		FirstAirDate string `json:"first_air_date"`
		ID           int    `json:"id"`
		Name         string `json:"name"`
		OriginalName string `json:"original_name"`
		PosterPath   string `json:"poster_path"`
	} `json:"results"`
}

type ThemoviedbResponse struct {
	ID               int                       `json:"id"`
	LastAirDate      string                    `json:"last_air_date"`
	LastEpisodeToAir *ThemoviedbItemResponse   `json:"last_episode_to_air"`
	NextEpisodeToAir *ThemoviedbItemResponse   `json:"next_episode_to_air"`
	NumberOfEpisodes int                       `json:"number_of_episodes"`
	NumberOfSeasons  int                       `json:"number_of_seasons"`
	OriginalName     string                    `json:"original_name"`
	Seasons          []*ThemoviedbItemResponse `json:"seasons"`
}
type ThemoviedbItemResponse struct {
	AirDate       string `json:"air_date"`
	EpisodeNumber int    `json:"episode_number"`
	ID            int    `json:"id"`
	EpisodeCount  int    `json:"episode_count"`
	Name          string `json:"name"`
	SeasonNumber  int    `json:"season_number"`
}
