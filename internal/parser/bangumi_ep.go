package parser

import (
	"GoBangumi/internal/errors"
	"regexp"
	"strconv"
	"strings"
)

var epRegxStep = []*regexp.Regexp{
	// 匹配ep，https://github.com/EstrellaXD/Auto_Bangumi/blob/97f078818a4f5b8513116a6032224d4e2f1dd7d9/src/parser/analyser/raw_parser.py
	regexp.MustCompile(`(.*|\[.*])( -? \d+ |\[\d+]|\[\d+.?[vV]\d{1}]|[第]\d+[话話集]|\[\d+.?END])(.*)`),
	// 取出数字
	regexp.MustCompile(`\d+`),
}

type BangumiEp struct {
}

func NewBangumiEp() *BangumiEp {
	return &BangumiEp{}
}

func ParseEp(title string) (int, error) {

	str := title
	str = strings.ReplaceAll(str, "【", "[")
	str = strings.ReplaceAll(str, "】", "]")
	res := epRegxStep[0].FindStringSubmatch(str)
	if res == nil {
		return 0, errors.ParseBangumiEpErr
	}
	epStr := epRegxStep[1].FindString(res[2])
	ep, err := strconv.Atoi(epStr)
	if err != nil || ep == 0 {
		return 0, errors.ParseBangumiEpErr
	}
	return ep, nil
}
