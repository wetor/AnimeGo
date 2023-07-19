package test_test

import (
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
