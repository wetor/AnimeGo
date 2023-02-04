package rss

import (
	"fmt"
	"testing"

	"github.com/wetor/AnimeGo/internal/animego/feed"
	"github.com/wetor/AnimeGo/internal/utils"
)

func TestRss_Parse(t *testing.T) {
	_ = utils.CreateMutiDir(feed.TempPath)
	feed.Init(&feed.Options{
		TempPath: "data",
	})
	r := NewRss("https://mikanani.me/RSS/MyBangumi?token=ky5DTt%2fMyAjCH2oKEN81FQ%3d%3d", "mikan")
	items := r.Parse()
	for _, item := range items {
		fmt.Println(item.Url, item.Name, item.Length, item.Date, item.Hash())
	}
}
