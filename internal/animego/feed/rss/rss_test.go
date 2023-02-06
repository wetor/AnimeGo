package rss_test

import (
	"fmt"
	"testing"

	"github.com/wetor/AnimeGo/internal/animego/feed/rss"
)

func TestRss_Parse(t *testing.T) {
	r := rss.NewRss("https://mikanani.me/RSS/MyBangumi?token=ky5DTt%2fMyAjCH2oKEN81FQ%3d%3d", "mikan")
	items := r.Parse()
	for _, item := range items {
		fmt.Println(item.Url, item.Name, item.Length, item.Date, item.Hash())
	}
}
