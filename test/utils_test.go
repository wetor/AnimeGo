package test_test

import (
	"github.com/wetor/AnimeGo/test"
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
第一行
skip1
skip2
skip3
第三行`)
	test.LogBatchCompare(r, nil,
		map[string]int{"一": 2, "skip": 3},
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
		map[string]int{"第一行": 2, "skip1": 1, "skip2": 1, "skip3": 1},
		"第三行",
	)
}
