package parser

import (
	"GoBangumi/internal/errors"
	"GoBangumi/internal/models"
	"regexp"
	"strconv"
	"strings"
)

var epRegx = []*regexp.Regexp{
	// 匹配ep，https://github.com/EstrellaXD/Auto_Bangumi/blob/97f078818a4f5b8513116a6032224d4e2f1dd7d9/src/parser/analyser/raw_parser.py
	regexp.MustCompile(`(.*|\[.*])( -? \d+ |\[\d+]|\[\d+.?[vV]\d{1}]|[第]\d+[话話集]|\[\d+.?END])(.*)`),
	// 取出数字
	regexp.MustCompile(`\d+`),
}

var tagRegx = []*regexp.Regexp{
	// 匹配分辨率 Resolution
	regexp.MustCompile(`1080|720|2160|4K`),
	// 匹配字幕 Subtitle
	regexp.MustCompile(`[简繁日字幕]|CH|BIG5|GB`),
	// 匹配源 Source
	regexp.MustCompile(`B-Global|[Bb]aha|[Bb]ilibili|AT-X|Web`),
}

var tagSplitRegx = regexp.MustCompile(`[\[\]()（）]`)

func ParseEp(title string) (int, error) {

	str := title
	str = strings.ReplaceAll(str, "【", "[")
	str = strings.ReplaceAll(str, "】", "]")
	res := epRegx[0].FindStringSubmatch(str)
	if res == nil {
		return 0, errors.ParseBangumiEpErr
	}
	epStr := epRegx[1].FindString(res[2])
	ep, err := strconv.Atoi(epStr)
	if err != nil || ep == 0 {
		return 0, errors.ParseBangumiEpErr
	}
	return ep, nil
}

func ParseTag(title string) (*models.ParseTagResult, error) {

	str := title
	str = strings.ReplaceAll(str, "【", "[")
	str = strings.ReplaceAll(str, "】", "]")

	tags := strings.Split(tagSplitRegx.ReplaceAllString(str, "  "), " ")

	result := &models.ParseTagResult{}

	for _, tag := range tags {
		if tagRegx[0].MatchString(tag) {
			result.Resolution = tag
		} else if tagRegx[1].MatchString(tag) {
			result.Subtitle = tag
		} else if tagRegx[2].MatchString(tag) {
			result.Source = tag
		}
	}
	return result, nil
}
