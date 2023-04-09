package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/wetor/AnimeGo/internal/models"
)

var (
	epRegxStr  = `(\d+|\d+.?\-.?\d+)`
	titleRegx  = regexp.MustCompile(strings.ReplaceAll(`(.*|\[.*])( -? {ep}|\[{ep}]|\[{ep}.?[vV]\d{1}]|[第]{ep}[话話集]|\[{ep}.?END])(.*)`, "{ep}", epRegxStr))
	numberRegx = regexp.MustCompile(epRegxStr)
)

func Parse(name string) *models.AnimeEpEntity {
	str := strings.NewReplacer("【", "[", "】", "]").Replace(name)
	fmt.Println(str)
	res := titleRegx.FindStringSubmatch(str)
	fmt.Println(res[2])
	titleBody := res[1]
	_ = titleBody
	titleEp := res[2]
	titleTags := res[3]
	_ = titleTags
	// ep
	epStr := numberRegx.FindString(titleEp)
	ep, err := strconv.Atoi(epStr)
	if err != nil {
		return nil
	}
	return &models.AnimeEpEntity{
		Ep: ep,
	}
}
