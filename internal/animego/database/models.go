package database

type AnimeDir struct {
	Dir       string
	SeasonDir map[int]string
}

type Options struct {
	SavePath string
}
