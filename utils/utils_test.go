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
