package themoviedb

import (
	"fmt"
	"net/url"
)

var idApi = func(key string, query string) string {
	url_, _ := url.Parse(Host() + "/3/discover/tv")
	q := url_.Query()
	q.Set("api_key", key)
	q.Set("sort_by", "first_air_date.desc")
	q.Set("language", "zh-CN")
	q.Set("timezone", "Asia/Shanghai")
	q.Set("with_genres", "16")
	q.Set("with_text_query", query)
	return url_.String() + "?" + q.Encode()
}

var infoApi = func(key string, id int) string {
	return fmt.Sprintf("%s/3/tv/%d?api_key=%s", Host(), id, key)
}
