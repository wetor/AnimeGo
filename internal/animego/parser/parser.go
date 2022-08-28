package parser

import (
	"GoBangumi/internal/errors"
	"GoBangumi/internal/models"
	"regexp"
	"strconv"
	"strings"
)

const (
	MatchTitleEp = 0
	MatchEpNum   = 1
)

var epRegx = []*regexp.Regexp{
	// 匹配ep，https://github.com/EstrellaXD/Auto_Bangumi/blob/97f078818a4f5b8513116a6032224d4e2f1dd7d9/src/parser/analyser/raw_parser.py
	regexp.MustCompile(`(.*|\[.*])( -? \d+ |\[\d+]|\[\d+.?[vV]\d{1}]|[第]\d+[话話集]|\[\d+.?END])(.*)`),
	// 取出数字
	regexp.MustCompile(`\d+`),
}

const (
	MatchTitleResolution = 0
	MatchTitleSubtitle   = 1
	MatchTitleSource     = 2
)

var tagRegx = []*regexp.Regexp{
	// 匹配分辨率 Resolution
	regexp.MustCompile(`1080|720|2160|4K`),
	// 匹配字幕 Subtitle
	regexp.MustCompile(`[简繁日字幕]|CH|BIG5|GB`),
	// 匹配源 Source
	regexp.MustCompile(`B-Global|[Bb]aha|[Bb]ilibili|AT-X|Web`),
}

var tagSplitRegx = regexp.MustCompile(`[\[\]()（）]`)

func ParseTitle(title string) (*models.ParseResult, error) {

	str := title
	str = strings.ReplaceAll(str, "【", "[")
	str = strings.ReplaceAll(str, "】", "]")
	res := epRegx[MatchTitleEp].FindStringSubmatch(str)
	if res == nil {
		return nil, errors.ParseAnimeTitleErr
	}
	if len(res) < 4 {
		return nil, errors.ParseAnimeTitleErr
	}
	titleBody := res[1]
	_ = titleBody
	titleEp := res[2]
	titleTags := res[3]
	result := &models.ParseResult{}
	// ep
	epStr := epRegx[MatchEpNum].FindString(titleEp)
	ep, err := strconv.Atoi(epStr)
	if err != nil || ep == 0 {
		return nil, errors.ParseAnimeTitleErr
	}
	result.Ep = ep

	// tags
	tags := strings.Split(tagSplitRegx.ReplaceAllString(titleTags, " "), " ")
	for _, tag := range tags {
		if len(result.Resolution) == 0 && tagRegx[MatchTitleResolution].MatchString(tag) {
			result.Resolution = tag
		} else if len(result.Subtitle) == 0 && tagRegx[MatchTitleSubtitle].MatchString(tag) {
			result.Subtitle = tag
		} else if len(result.Source) == 0 && tagRegx[MatchTitleSource].MatchString(tag) {
			result.Source = tag
		}
	}
	return result, nil
}
