package themoviedb

import (
	"fmt"
	"github.com/wetor/AnimeGo/internal/constant"
	"net/url"
)

var idApi = func(query string, isMovie bool) string {
	uri := "/3/discover/tv"
	if isMovie {
		uri = "/3/discover/movie"
	}
	url_, _ := url.Parse(constant.ThemoviedbHost + uri)
	q := url_.Query()
	q.Set("sort_by", "first_air_date.desc")
	q.Set("language", "zh-CN")
	q.Set("timezone", "Asia/Shanghai")
	q.Set("with_genres", "16")
	q.Set("with_text_query", query)
	return url_.String() + "?" + q.Encode()
}

var infoApi = func(id int, isMovie bool) string {
	uri := "/3/tv"
	if isMovie {
		uri = "/3/movie"
	}
	return fmt.Sprintf("%s/%d", constant.ThemoviedbHost+uri, id)
}
