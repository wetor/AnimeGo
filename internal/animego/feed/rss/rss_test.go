package rss_test

import (
	"fmt"
	"testing"

	"github.com/wetor/AnimeGo/internal/animego/feed/rss"
	"github.com/wetor/AnimeGo/test"
)

func TestRss_Parse(t *testing.T) {
	r := rss.NewRss(&rss.Options{File: test.GetDataPath("feed", "Mikan.xml")})
	items, _ := r.Parse()
	for _, item := range items {
		fmt.Println(item.Url, item.Name)
	}
}
