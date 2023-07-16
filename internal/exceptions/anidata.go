package exceptions

import (
	"fmt"

	"github.com/wetor/AnimeGo/internal/api"
)

type ErrAniData struct {
	AniData api.AniData
}

func (e ErrAniData) Error() string {
	return fmt.Sprintf("处理 %s 失败", e.AniData.Name())
}

type ErrAniDataGet struct {
	AniData api.AniData
}

func (e ErrAniDataGet) Error() string {
	return fmt.Sprintf("获取 %s 信息失败", e.AniData.Name())
}

type ErrAniDataParse struct {
	AniData api.AniData
}

func (e ErrAniDataParse) Error() string {
	return fmt.Sprintf("解析 %s 信息失败", e.AniData.Name())
}

type ErrAniDataSearch struct {
	AniData api.AniData
}

func (e ErrAniDataSearch) Error() string {
	return fmt.Sprintf("查询 %s 信息失败", e.AniData.Name())
}

type ErrMikanParseHTML struct {
	Message string
}

func (e ErrMikanParseHTML) Error() string {
	if len(e.Message) == 0 {
		return "解析网页失败"
	}
	return fmt.Sprintf("解析 %s 失败，解析网页错误", e.Message)
}

type ErrThemoviedbMatchSeason struct {
	Message string
}

func (e ErrThemoviedbMatchSeason) Error() string {
	return fmt.Sprintf("匹配季度信息失败，%s", e.Message)
}

type ErrThemoviedbSearchName struct {
}

func (e ErrThemoviedbSearchName) Error() string {
	return "搜索番剧名失败"
}

type ErrBangumiCacheNotFound struct {
	BangumiID int
}

func (e ErrBangumiCacheNotFound) Error() string {
	return fmt.Sprintf("获取Bangumi %v 缓存数据异常，使用在线数据", e.BangumiID)
}
