package bangumi

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/modules/cache"
	"GoBangumi/modules/parser"
	"GoBangumi/store"
	"fmt"
	"github.com/antchfx/htmlquery"
	"go.uber.org/zap"
	"net/url"
	"strconv"
	"strings"
)

const (
	MikanIdXPath         = "//a[@class='mikan-rss']"                                 // Mikan番剧id获取XPath
	MikanBangumiUrlXPath = "//p[@class='bangumi-info']/a[contains(@href, 'bgm.tv')]" // Mikan番剧信息中bangumi id获取XPath
)

var MikanInfoUrl = func(id int) string {
	return fmt.Sprintf("%s/Home/Bangumi/%d", config.Advanced().Mikan().Host, id)
}

type Mikan struct {
}

func NewMikan() Bangumi {
	return &Mikan{}
}

func (b *Mikan) Parse(opt *models.BangumiParseOptions) *models.Bangumi {
	zap.S().Infof("获取「%s」信息开始...", opt.Name)
	// ------------------- 解析文件名获取ep -------------------
	epParser := parser.NewBangumiEp()
	ep := epParser.Parse(&models.ParseNameOptions{
		Name: opt.Name,
	})
	if ep == nil {
		zap.S().Warn("解析ep信息失败，结束此流程")
		return nil
	}
	// ------------------- 获取mikanID -------------------
	mikanID := b.parseMikan1(opt.Url)
	if mikanID == 0 {
		zap.S().Warn("获取Mikan ID失败，结束此流程")
		return nil
	}
	// ------------------- 获取bangumiID -------------------
	bangumiID := b.parseMikan2(mikanID)
	if bangumiID == 0 {
		zap.S().Warn("获取bangumi ID失败，结束此流程")
		return nil
	}
	// ------------------- 获取bangumi信息 -------------------
	info := b.parseBangumi(bangumiID, ep.Ep, opt.Date)
	if info == nil {
		zap.S().Warn("获取Bangumi信息失败，结束此流程")
		return nil
	}
	// ------------------- 获取tmdb信息(季度信息) -------------------
	info.BangumiSeason = b.parseThemoviedb(info.Name, info.AirDate)
	if info.BangumiSeason == nil || info.Season == 0 {
		zap.S().Warn("获取Themoviedb季度信息失败，结束此流程")
		return nil
	}
	info.BangumiExtra = &models.BangumiExtra{
		SubID:  mikanID,
		SubUrl: opt.Url,
	}
	zap.S().Infof("获取「%s」信息成功！更名为「%s」", opt.Name, info.FullName())
	return info
}

// parseMikan1
//  @Description: 解析mikan rss中的link页面，获取当前资源的mikan id
//  @receiver b
//  @param url_
//  @return int
//
func (b *Mikan) parseMikan1(url_ string) (mikanID int) {
	tmp := store.Cache.Get(cache.RssMikanBucket, url_)
	if tmp != nil {
		if val, ok := tmp.(int); ok {
			zap.S().Debugf("步骤1，解析Mikan，缓存")
			return val
		}
	}

	zap.S().Debugf("步骤1，解析Mikan，%s", url_)
	doc, err := htmlquery.LoadURL(url_)
	if err != nil {
		zap.S().Warn(err)
		return 0
	}
	miaknLink := htmlquery.FindOne(doc, MikanIdXPath)
	href := htmlquery.SelectAttr(miaknLink, "href")
	u, err := url.Parse(href)
	if err != nil {
		zap.S().Warn(err)
		return 0
	}
	query := u.Query()
	if query.Has("bangumiId") {
		id, err := strconv.Atoi(query.Get("bangumiId"))
		if err != nil {
			zap.S().Warn(err)
			return 0
		}
		mikanID = id
	}
	if mikanID == 0 {
		zap.S().Warn("获取Mikan ID失败")
		return 0
	}
	store.Cache.Put(cache.RssMikanBucket, url_, mikanID, config.Advanced().Mikan().CacheIdExpire)
	return mikanID
}

// parseMikan2
//  @Cache mikan_bangumi mikanID
//  @Description: 通过mikan id解析mikan番剧信息页面，获取bgm.tv id
//  @receiver b
//  @param mikanID
//  @return int
//
func (b *Mikan) parseMikan2(mikanID int) (bangumiID int) {
	// 通过mikanID查询缓存中的bangumiID
	tmp := store.Cache.Get(cache.MikanBangumiBucket, mikanID)
	if tmp != nil {
		if val, ok := tmp.(int); ok {
			zap.S().Debugf("步骤2，解析Mikan，缓存")
			return val
		}
	}

	url_ := MikanInfoUrl(mikanID)
	zap.S().Debugf("步骤2，解析Mikan，%s", url_)
	doc, err := htmlquery.LoadURL(url_)
	if err != nil {
		zap.S().Warn(err)
		return 0
	}
	bangumiUrl := htmlquery.FindOne(doc, MikanBangumiUrlXPath)
	href := htmlquery.SelectAttr(bangumiUrl, "href")

	//fmt.Println(href)
	hrefSplit := strings.Split(href, "/")
	bangumiID, err = strconv.Atoi(hrefSplit[len(hrefSplit)-1])
	if err != nil {
		zap.S().Warn(err)
		return 0
	}
	// mikanID和bangumiID对应关系固定，缓存
	store.Cache.Put(cache.MikanBangumiBucket, mikanID, bangumiID, config.Advanced().Mikan().CacheBangumiExpire)
	return bangumiID
}

//  parseBangumi
//  @Description: 解析bgm.tv，获取番剧信息和ep信息
//  @receiver b
//  @param bangumiID
//  @param ep
//  @return *models.Bangumi
//
func (b *Mikan) parseBangumi(bangumiID, ep int, date string) *models.Bangumi {
	zap.S().Debugf("步骤3，解析Bangumi，%d", bangumiID)
	bangumi := NewBgm()
	newBgm := bangumi.Parse(&models.BangumiParseOptions{
		ID:   bangumiID,
		Ep:   ep,
		Date: date,
	})
	return newBgm
}

// parseThemoviedb
//  @Description: 从Themoviedb网站获取当前季度
//  @receiver b
//  @param name
//  @param airDate
//  @return int
//
func (b *Mikan) parseThemoviedb(name, airDate string) *models.BangumiSeason {
	zap.S().Debugf("步骤4，解析Themoviedb，%s", name)
	tmdb := NewThemoviedb()
	newBgm := tmdb.Parse(&models.BangumiParseOptions{
		Name: name,
		Date: airDate,
	})
	return newBgm.BangumiSeason
}
