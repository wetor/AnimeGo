package themoviedb

type FindResponse struct {
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

type InfoResponse struct {
	ID               int             `json:"id"`
	LastAirDate      string          `json:"last_air_date"`
	LastEpisodeToAir *ItemResponse   `json:"last_episode_to_air"`
	NextEpisodeToAir *ItemResponse   `json:"next_episode_to_air"`
	NumberOfEpisodes int             `json:"number_of_episodes"`
	NumberOfSeasons  int             `json:"number_of_seasons"`
	OriginalName     string          `json:"original_name"`
	Seasons          []*ItemResponse `json:"seasons"`
}
type ItemResponse struct {
	AirDate       string `json:"air_date"`
	EpisodeNumber int    `json:"episode_number"`
	ID            int    `json:"id"`
	EpisodeCount  int    `json:"episode_count"`
	Name          string `json:"name"`
	SeasonNumber  int    `json:"season_number"`
}
