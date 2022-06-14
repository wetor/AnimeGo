package utils

import (
	"GoBangumi/model"
	"fmt"
	"github.com/xxxsen/qbapi"
	"testing"
)

func TestConvertModel(t *testing.T) {
	src := &qbapi.TorrentListItem{
		Name:        "测试标题",
		ContentPath: "1111111",
	}

	dst := &model.TorrentItem{}
	ConvertModel(src, dst)
	fmt.Println(dst)
}

func TestHttpGet(t *testing.T) {
	err := HttpGet("https://mikanani.me/RSS/MyBangumi?token=CE6CRA3j0Sf4hsGI6eH3Fg%3d%3d", "../data/cache/rss.xml")
	if err != nil {
		return
	}
}
