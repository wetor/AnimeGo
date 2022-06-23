package parser

import (
	"GoBangumi/models"
	"regexp"
	"strconv"
	"strings"
)

var epRegxStep = []*regexp.Regexp{
	// 匹配ep，https://github.com/EstrellaXD/Auto_Bangumi/blob/33454cf23578017cc92e31fd98f4c4d7351cdf7f/auto_bangumi/parser/analyser/raw_parser.py
	regexp.MustCompile(`(.*|\[.*])( -? \d{1,3} |\[\d{1,3}]|\[\d{1,3}.?[vV]\d{1}]|[第]\d{1,3}[话話集]|\[\d{1,3}.?END])(.*)`),
	// 取出数字
	regexp.MustCompile(`\d{1,3}`),
}

type BangumiEp struct {
}

func NewBangumiEp() Parser {
	return &BangumiEp{}
}

func (p *BangumiEp) Parse(opt *models.ParseNameOptions) *models.ParseResult {

	str := opt.Name
	str = strings.ReplaceAll(str, "【", "[")
	str = strings.ReplaceAll(str, "】", "]")
	res := epRegxStep[0].FindStringSubmatch(str)
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
