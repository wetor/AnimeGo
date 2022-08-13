package utils

import (
	"fmt"
	"testing"
)

func TestConvertModel(t *testing.T) {
	//src := &qbapi.TorrentListItem{
	//	Name:        "测试标题",
	//	ContentPath: "1111111",
	//}
	//
	//dst := &models.TorrentItem{}
	//ConvertModel(src, dst)
	//fmt.Println(dst)
}

func TestApiGet(t *testing.T) {
	n, e := ApiGet("https://google.com", nil, "socks5://127.0.0.1:7891")
	fmt.Println(n, e)
}
