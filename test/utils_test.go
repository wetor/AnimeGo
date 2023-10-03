package test_test

import (
	"fmt"
	"github.com/wetor/AnimeGo/test"
	"io"
	"strings"
	"testing"
)

func TestBatchCompare(t *testing.T) {
	r := strings.NewReader(`第一行
		第二行
		skip1
		skip2
		skip3
		第三行`)
	test.LogBatchCompare(r, nil,
		"一",
		"二",
		3,
		"三",
	)

	r = strings.NewReader(`第一行
		第三行`)
	test.LogBatchCompare(r, nil,
		"一",
		"#!二",
		"三",
	)

	r = strings.NewReader(`第一行
		第二行
		skip1
		skip2
		skip3
		第三行`)
	test.LogBatchCompare(r, nil,
		[]string{"一", "二"},
		3,
		"三",
	)

	r = strings.NewReader(`第一行
		第二行
		第三行
		第四行`)
	test.LogBatchCompare(r, nil,
		[]string{"一", "二", "三"},
		"四",
	)
}

func TestBatchCompare2(t *testing.T) {
	var r io.Reader
	r = strings.NewReader(`第一行
		第一行
		skip1
		skip2
		skip3
		第三行`)
	test.LogBatchCompare(r, nil,
		map[string]any{"一": 2, "skip": 3},
		"三",
	)

	r = strings.NewReader(`第一行
		第一行
		skip1
		skip2
		skip3
		第三行`)
	test.LogBatchCompare(r, nil,
		map[string]any{
			"一":    test.NewRange(0, 2),
			"skip": 3,
		},
		"三",
	)

	r = strings.NewReader(`第一行
第一行
skip1
skip2
skip3
第三行`)
	test.LogBatchCompare(r, func(line, match string) bool {
		return line == match
	},
		map[string]any{"第一行": 2, "skip1": 1, "skip2": 1, "skip3": 1},
		"第三行",
	)

}

func TestCompareJSON(t *testing.T) {
	type BaseDBEntity struct {
		Hash     string `json:"hash"`
		Name     string `json:"name"`
		CreateAt int64  `json:"create_at"`
		UpdateAt int64  `json:"update_at"`
	}
	type StateDB struct {
		Seeded     bool `json:"seeded"`     // 是否做种
		Downloaded bool `json:"downloaded"` // 是否已下载完成
		Renamed    bool `json:"renamed"`    // 是否已重命名/移动
		Scraped    bool `json:"scraped"`    // 是否已经完成搜刮
	}
	type EpisodeDBEntity struct {
		BaseDBEntity `json:"info"`
		StateDB      `json:"state"`
		Season       int  `json:"season"`
		Type         int8 `json:"type"`
		Ep           int  `json:"ep"`
	}

	jsonFile := `
{
  "info": {
    "hash": "95d89a8afb97b3a7ed75fa3f7559adac",
    "name": "动画1",
    "create_at": 1692585997,
    "update_at": 1692585999
  },
  "state": {
    "downloaded": true,
    "seeded": true,
    "renamed": true,
    "scraped": true
  },
  "season": 2,
  "type": 1,
  "ep": 1
}
`
	data := &EpisodeDBEntity{
		BaseDBEntity: BaseDBEntity{
			Hash: "95d89a8afb97b3a7ed75fa3f7559adac",
			Name: "动画1",
		},
		StateDB: StateDB{
			Downloaded: true,
			Seeded:     true,
			Renamed:    true,
			Scraped:    true,
		},
		Season: 3,
		Type:   1,
		Ep:     1,
	}
	skipFiled := []string{"create_at", "update_at"}
	isEqual, err := test.CompareJSON([]byte(jsonFile), data, skipFiled)
	if err != nil {
		fmt.Println(err)
		return
	}
	if isEqual {
		fmt.Println("两个结构体相等")
	} else {
		fmt.Println("两个结构体不相等")
	}
}
