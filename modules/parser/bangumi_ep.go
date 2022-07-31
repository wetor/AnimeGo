package parser

import (
	"GoBangumi/models"
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

func NewBangumiEp() Parser {
	return &BangumiEp{}
}

func (p *BangumiEp) Parse(opt *models.ParseOptions) *models.ParseResult {

	str := opt.Name
	str = strings.ReplaceAll(str, "【", "[")
	str = strings.ReplaceAll(str, "】", "]")
	res := epRegxStep[0].FindStringSubmatch(str)
	if res == nil {
		return nil
	}
	epStr := epRegxStep[1].FindString(res[2])
	ep, err := strconv.Atoi(epStr)
	if err != nil || ep == 0 {
		return nil
	}
	return &models.ParseResult{
		ParseEpResult: &models.ParseEpResult{
			Ep: ep,
		},
	}
}
