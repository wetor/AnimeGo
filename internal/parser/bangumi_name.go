package parser

import (
	"GoBangumi/internal/errors"
	"go.uber.org/zap"
	"regexp"
)

var nameRegxStep = []*regexp.Regexp{
	// 删除末尾 "第X季"或"X季"
	regexp.MustCompile(`\s?第?(\d{1,2}|(一|二|三|四|五|伍|六|七|八|九|十))(期|部|季|篇|章|編)$`),
	// 删除末尾 "1st Season"、"Xnd Season"或"Season X"或"X"
	regexp.MustCompile(`\s?(\d{1,2}(st|nd|rd|th)\s?Season|Season\s?\d{1,2}|\d{1,2})$`),
	// 删除 "X篇"、"X季"之后内容
	regexp.MustCompile(`\s(.*?)(期|部|季|篇|章|編).*$`),
	// 删除"IV"或"2"、"3"等或" "之后内容
	regexp.MustCompile(`\s?((V|X|IX|IV|V?I{1,3})|[2-9]|[1-9]\d).*$|\s\S+$`),
}

// Parse
//  @Description: 从 nameRegxStep[opt.StartStep] 开始执行，并返回下一步的index
//  @receiver *BangumiName
//  @param opt *models.ParseOptions
//  @return *models.ParseResult
//
func ParseName(name string, step int) (string, int, error) {
	if step < 0 {
		return "", 0, errors.ParseBangumiNameErr
	}
	if step >= len(nameRegxStep) {
		zap.S().Warn("BangumiName Step错误")
		return "", 0, errors.ParseBangumiNameErr
	}
	str := name
	i := step
	for ; i < len(nameRegxStep); i++ {
		has := nameRegxStep[i].MatchString(str)
		if has {
			str = nameRegxStep[i].ReplaceAllString(str, "")
			break
		}
	}
	i++ // 下一步
	if i >= len(nameRegxStep) {
		i = -1
	}

	return str, i, nil
}
