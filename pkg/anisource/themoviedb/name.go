package themoviedb

import (
	"AnimeGo/pkg/errors"
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

// RemoveNameSuffix
//  @Description: 删除番剧名的一些后缀，使得themoviedb能够正常搜索到
//  @param name string 番剧名
//  @param fun Action 执行操作。
// 					返回obj, nil表示成功；
//					返回nil, nil进行下一步；步骤执行完则返回nil, err
//					返回nil, err则直接结束；
//  @return interface{} fun返回的obj
//  @return error
//
func RemoveNameSuffix(name string, fun func(string) (interface{}, error)) (interface{}, error) {
	currStep := 0
	for ; currStep < len(nameRegxStep)+1; currStep++ {
		result, err := fun(name)
		if err != nil {
			return nil, err
		}
		if result != nil {
			return result, nil
		}
		if currStep < len(nameRegxStep) {
			has := nameRegxStep[currStep].MatchString(name)
			if has {
				newName := nameRegxStep[currStep].ReplaceAllString(name, "")
				if len(newName) > 0 &&  len(newName) > len(name)/10  {
					name = newName
				}
				zap.S().Debugf("重新搜索：「%s」", name)
			}
		}
	}
	return nil, errors.NewAniError("解析番剧名失败")
}

// SimilarText
//  @Description: 字符串相似度计算
//  @param first string
//  @param second string
//  @param percent *float64
//  @return int
//
func SimilarText(first, second string) (percent float64) {
	var similarText func(string, string, int, int) int
	similarText = func(str1, str2 string, len1, len2 int) int {
		var sum, max int
		pos1, pos2 := 0, 0

		// Find the longest segment of the same section in two strings
		for i := 0; i < len1; i++ {
			for j := 0; j < len2; j++ {
				for l := 0; (i+l < len1) && (j+l < len2) && (str1[i+l] == str2[j+l]); l++ {
					if l+1 > max {
						max = l + 1
						pos1 = i
						pos2 = j
					}
				}
			}
		}

		if sum = max; sum > 0 {
			if pos1 > 0 && pos2 > 0 {
				sum += similarText(str1, str2, pos1, pos2)
			}
			if (pos1+max < len1) && (pos2+max < len2) {
				s1 := []byte(str1)
				s2 := []byte(str2)
				sum += similarText(string(s1[pos1+max:]), string(s2[pos2+max:]), len1-pos1-max, len2-pos2-max)
			}
		}

		return sum
	}

	l1, l2 := len(first), len(second)
	if l1+l2 == 0 {
		return 0
	}
	sim := similarText(first, second, l1, l2)
	return float64(sim*200) / float64(l1+l2)
}
