package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/wetor/AnimeGo/pkg/xpath"
	"sort"
	"strconv"
	"strings"
)

func format(format string, p map[string]any) string {
	args, i := make([]string, len(p)*2), 0
	for k, v := range p {
		args[i] = "{" + k + "}"
		args[i+1] = fmt.Sprint(v)
		i += 2
	}
	return strings.NewReplacer(args...).Replace(format)
}

var filenameMap = map[string]any{
	`/`: "",
	`\`: "",
	`[`: "(",
	`]`: ")",
	`:`: "-",
	`;`: "-",
	`=`: "-",
	`,`: "-",
}

func FileName(filename string) string {
	return format(filename, filenameMap)
}

func Md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func AnimeToFileName(a *AnimeEntity, index int) string {
	if a.Ep[index].Type == AnimeEpUnknown {
		file := strings.TrimSuffix(xpath.Base(a.Ep[index].Src), xpath.Ext(a.Ep[index].Src))
		return xpath.Join(fmt.Sprintf("S%02d", a.Season), file)
	}
	return xpath.Join(fmt.Sprintf("S%02d", a.Season), fmt.Sprintf("E%03d", a.Ep[index].Ep))
}

func AnimeToFilePath(a *AnimeEntity) []string {
	filePath := make([]string, 0, len(a.Ep))
	dir := FileName(a.NameCN)
	for i, ep := range a.Ep {
		filePath = append(filePath, xpath.Join(dir, AnimeToFileName(a, i)+xpath.Ext(ep.Src)))
	}
	return filePath
}

func AnimeToFilePathSrc(a *AnimeEntity) []string {
	filePath := make([]string, 0, len(a.Ep))
	for _, ep := range a.Ep {
		filePath = append(filePath, ep.Src)
	}
	return filePath
}

func ToMetaData(tmdbID, bangumiID int) string {
	nfoTemplate := "<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\"?>\n<tvshow>\n"
	if tmdbID != 0 {
		nfoTemplate += "  <tmdbid>{tmdbid}</tmdbid>\n"
	}
	nfoTemplate += "  <bangumiid>{bangumiid}</bangumiid>\n</tvshow>"
	return format(nfoTemplate, map[string]any{
		"tmdbid":    tmdbID,
		"bangumiid": bangumiID,
	})
}

func ToIntervals(eps []*AnimeEpEntity) string {
	// 初始化区间的起始值和结束值
	start, end := eps[0].Ep, eps[0].Ep
	// 初始化区间的字符串表示
	var intervalSlice []string
	for i := 1; i < len(eps); i++ {
		num := eps[i].Ep
		// 如果当前数与上一个数相差1，表示当前数与上一个数在同一个区间内
		if num-eps[i-1].Ep == 1 {
			end = num // 更新区间的结束值
		} else {
			// 否则，表示当前数与上一个数不在同一个区间内，需要将上一个区间输出，并重新开始一个新的区间
			if start == end {
				intervalSlice = append(intervalSlice, strconv.Itoa(start))
			} else {
				intervalSlice = append(intervalSlice, fmt.Sprintf("%d-%d", start, end))
			}
			start, end = num, num // 更新区间的起始值和结束值
		}
	}
	// 将最后一个区间输出
	if start == end {
		intervalSlice = append(intervalSlice, strconv.Itoa(start))
	} else {
		intervalSlice = append(intervalSlice, fmt.Sprintf("%d-%d", start, end))
	}
	// 对区间切片进行排序
	sort.Strings(intervalSlice)
	// 将排序后的区间切片合并成字符串
	intervals := strings.Join(intervalSlice, ",")
	return intervals
}
