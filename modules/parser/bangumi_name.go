package parser

import (
	"GoBangumi/models"
	"github.com/golang/glog"
	"regexp"
)

var regStep = []*regexp.Regexp{
	// 删除末尾 "第X季"或"X季"
	regexp.MustCompile(`\s?第?(\d{1,2}|(一|二|三|四|五|伍|六|七|八|九|十))(期|部|季|篇|章|編)$`),
	// 删除末尾 "1st Season"、"Xnd Season"或"Season X"或"X"
	regexp.MustCompile(`\s?(\d{1,2}(st|nd|rd|th)\s?Season|Season\s?\d{1,2}|\d{1,2})$`),
	// 删除 "X篇"、"X季"之后内容
	regexp.MustCompile(`\s(.*?)(期|部|季|篇|章|編).*$`),
	// 删除"IV"或"2"、"3"等或" "之后内容
	regexp.MustCompile(`\s?((V|X|IX|IV|V?I{1,3})|[2-9]|[1-9]\d).*$|\s\S+$`),
}

type BangumiName struct {
}

func NewBangumiName() Parser {
	return &BangumiName{}
}

func (p *BangumiName) ParseBangumiName(opt *models.ParseBangumiNameOptions) *models.BangumiName {
	if opt.StartStep < 0 {
		return nil
	}
	if opt.StartStep >= len(regStep) {
		glog.Errorln("BangumiName Step错误")
		return nil
	}
	str := opt.Name
	i := opt.StartStep
	for ; i < len(regStep); i++ {
		has := regStep[i].MatchString(str)
		if has {
			str = regStep[i].ReplaceAllString(str, "")
			break
		}
	}
	i++ // 下一步
	if i >= len(regStep) {
		i = -1
	}

	return &models.BangumiName{
		Name:     str,
		NextStep: i,
	}
}
