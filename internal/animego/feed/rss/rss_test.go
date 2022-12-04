package rss

import (
	"fmt"
	"github.com/wetor/AnimeGo/test"
	"testing"
)

func TestRss_Parse(t *testing.T) {
	test.TestInit()
	r := NewRss("https://share.dmhy.org/topics/rss/rss.xml", "dmhy")
	items := r.Parse()
	for _, item := range items {
		fmt.Println(item.Url, item.Name, item.Length, item.Date, item.Hash())
	}
}
