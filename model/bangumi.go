package model

import (
	"fmt"
)

type Bangumi struct {
	ID     int    // bgm id
	SubID  int    // 其他id，mikan id
	Name   string // 名称，从bgm获取
	NameJp string // 日文名，从bgm获取
	Date   string // 播放日期，从bgm获取
	Season int    // 当前季，从themoviedb获取
	Ep     int    // 当前集，从下载文件名解析
	Eps    int    // 总集数，从bgm获取
}

func (b *Bangumi) FullName() string {
	str := fmt.Sprintf("%s[第%d季][第%d集]", b.Name, b.Season, b.Ep)
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
