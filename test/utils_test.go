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
	test.LogBatchCompare(r,
		"一",
		"二",
		3,
		"三",
	)
}
