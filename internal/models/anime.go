package models

import (
	"fmt"
)

type AnimeEntity struct {
	ID      int    // bangumi id
	Name    string // 名称，从bgm获取
	NameCN  string // 中文名称，从bgm获取
	AirDate string // 最初播放日期，从bgm获取
	Eps     int    // 总集数，从bgm获取
	*AnimeSeason
	*AnimeEp
	*AnimeExtra
	*TorrentInfo
}

type AnimeSeason struct {
	Season int // 当前季，从themoviedb获取
}
type AnimeEp struct {
	Ep       int    // 当前集，从下载文件名解析
	Date     string // 当前集播放日期，从bgm获取
	Duration string // 当前集时长
	EpDesc   string // 当前集简介
	EpName   string // 当前集标题
	EpNameCN string // 当前集中文标题
	EpID     int    // 当前集bgm id

}
type AnimeExtra struct {
	SubID  int    // 其他id，mikan id
	SubUrl string // 其他url，mikan当前集的url
}

type TorrentInfo struct {
	Url  string // 当前集种子链接
	Hash string // 当前集种子Hash，唯一ID
}

func (b *AnimeEntity) FullName() string {
	// TODO: 重命名
	str := fmt.Sprintf("%s[第%d季][第%d集]", b.NameCN, b.Season, b.Ep)
	return str
}

func (b *AnimeEntity) DirName() string {
	str := fmt.Sprintf("%s-第%d季-第%d集", b.NameCN, b.Season, b.Ep)
	return str
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
