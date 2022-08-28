package themoviedb

import (
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

// RemoveNameSuffix
//  @Description: 删除番剧名的一些后缀，使得themoviedb能够正常搜索到
//  @param name string 番剧名
//  @param step int 后缀移除阶段，从0开始。共分为四个阶段:
//                    step1: 删除末尾 "第X季"或"X季"
//                    step2: 删除末尾 "1st Season"、"Xnd Season"或"Season X"或"X"
//                    step3: 删除 "X篇"、"X季"之后内容
//                    step4: 删除"IV"或"2"、"3"等或" "之后内容
//  @return nextName string 删除后缀后的番剧名
//  @return nextStep int 下一个step，-1表示四个阶段全部完成
//  @return err error
//
func RemoveNameSuffix(name string, step int) (nextName string, nextStep int, err error) {
	if step < 0 {
		return "", 0, ParseAnimeNameErr
	}
	if step >= len(nameRegxStep) {
		return "", 0, ParseAnimeNameErr
	}
	nextName = name
	nextStep = step
	for ; nextStep < len(nameRegxStep); nextStep++ {
		has := nameRegxStep[nextStep].MatchString(nextName)
		if has {
			nextName = nameRegxStep[nextStep].ReplaceAllString(nextName, "")
			break
		}
	}
	nextStep++ // 下一步
	if nextStep >= len(nameRegxStep) {
		nextStep = -1
	}
	return nextName, nextStep, nil
}
