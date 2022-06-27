package utils

import (
	"GoBangumi/models"
	"fmt"
	"github.com/xxxsen/qbapi"
	"testing"
	"time"
)

func TestConvertModel(t *testing.T) {
	src := &qbapi.TorrentListItem{
		Name:        "测试标题",
		ContentPath: "1111111",
	}

	dst := &models.TorrentItem{}
	ConvertModel(src, dst)
	fmt.Println(dst)
}

func TestApiGet(t *testing.T) {
	n, e := ApiGet("https://google.com", nil, "socks5://127.0.0.1:7891")
	fmt.Println(n, e)
}

func TestToBytes(t *testing.T) {
	b := ToBytes(&models.Bangumi{
		ID:     1000,
		Name:   "测试日文",
		NameCN: "测试中文",
		BangumiExtra: &models.BangumiExtra{
			SubID:  22,
			SubUrl: "hasdtasdasdas",
		},
	}, time.Now().Unix()+30)
	fmt.Println(b)
	v, e := ToValue(b)
	bgm := v.(*models.Bangumi)
	fmt.Println(bgm, bgm.BangumiExtra, e)
}
