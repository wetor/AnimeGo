package utils

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	epTitleRegx  = regexp.MustCompile(`(.*|\[.*])( -? \d+|\[\d+]|\[\d+.?[vV]\d{1}]|[第]\d+[话話集]|\[\d+.?END])(.*)`)
	epNumberRegx = regexp.MustCompile(`\d+`)
)

func ParseEp(name string) (ep int) {
	str := strings.NewReplacer("【", "[", "】", "]").Replace(name)
	res := epTitleRegx.FindStringSubmatch(str)
	if len(res) < 3 {
		return 0
	}
	titleBody := res[1]
	_ = titleBody
	titleEp := res[2]
	// titleTags := res[3]
	epStr := epNumberRegx.FindString(titleEp)
	ep, err := strconv.Atoi(epStr)
	if err != nil {
		return 0
	}
	return ep
}
