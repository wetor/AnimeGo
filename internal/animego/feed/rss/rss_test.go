package rss_test

import (
	"fmt"
	"testing"

	"github.com/wetor/AnimeGo/internal/animego/feed/rss"
)

func TestRss_Parse(t *testing.T) {
	r := rss.NewRss(&rss.Options{File: "testdata/Mikan.xml"})
	items := r.Parse()
	for _, item := range items {
		fmt.Println(item.Url, item.Name, item.Length, item.Date, item.Hash)
	}
}
