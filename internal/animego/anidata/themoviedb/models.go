package themoviedb

type Options struct {
	name    string
	airDate string
}

type Entity struct {
	ID      int    `json:"id"`             // tmdb ID
	NameCN  string `json:"name"`           // 中文名
	Name    string `json:"original_name"`  // 原名
	AirDate string `json:"first_air_date"` // 番剧第一季开播时间
	*SeasonInfo
}

type FindResponse struct {
	Page         int       `json:"page"`
	TotalPages   int       `json:"total_pages"`
	TotalResults int       `json:"total_results"`
	Result       []*Entity `json:"results"`
}

type SeasonInfo struct {
	Season  int    `json:"season_number"`
	AirDate string `json:"air_date"` // 番剧当前季度开播时间
	EpID    int    `json:"id"`
	EpName  string `json:"name"`
	Ep      int    `json:"episode_number"`
	Eps     int    `json:"episode_count"`
}

type InfoResponse struct {
	ID               int           `json:"id"`
	LastAirDate      string        `json:"last_air_date"`
	LastEpisodeToAir *SeasonInfo   `json:"last_episode_to_air"`
	NextEpisodeToAir *SeasonInfo   `json:"next_episode_to_air"`
	NumberOfEpisodes int           `json:"number_of_episodes"`
	NumberOfSeasons  int           `json:"number_of_seasons"`
	OriginalName     string        `json:"original_name"`
	Seasons          []*SeasonInfo `json:"seasons"`
}
