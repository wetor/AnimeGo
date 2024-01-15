package themoviedb

import (
	"fmt"
	"net/url"
)

var idApi = func(host string, key string, query string, isMovie bool) string {
	uri := "/3/discover/tv"
	if isMovie {
		uri = "/3/discover/movie"
	}
	url_, _ := url.Parse(Host(host) + uri)
	q := url_.Query()
	q.Set("api_key", key)
	q.Set("sort_by", "first_air_date.desc")
	q.Set("language", "zh-CN")
	q.Set("timezone", "Asia/Shanghai")
	q.Set("with_genres", "16")
	q.Set("with_text_query", query)
	return url_.String() + "?" + q.Encode()
}

var infoApi = func(host string, key string, id int, isMovie bool) string {
	uri := "/3/tv"
	if isMovie {
		uri = "/3/movie"
	}
	return fmt.Sprintf("%s/%d?api_key=%s", Host(host)+uri, id, key)
}
