package bangumi

import (
	"GoBangumi/models"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/golang/glog"
	"net/url"
	"strconv"
	"strings"
)

const (
	MikanBaseUrl         = "https://mikanani.me"                                     // Mikan 域名
	MikanIdXPath         = "//a[@class='mikan-rss']"                                 // Mikan番剧id获取XPath
	MikanBangumiUrlXPath = "//p[@class='bangumi-info']/a[contains(@href, 'bgm.tv')]" // Mikan番剧信息中bangumi id获取XPath
)

var MikanInfoUrl = func(id int) string {
	return fmt.Sprintf("%s/Home/Bangumi/%d", MikanBaseUrl, id)
}

type Mikan struct {
}

func NewMikan() Bangumi {
	return &Mikan{}
}

func (b *Mikan) Parse(opt *models.BangumiParseOptions) *models.Bangumi {
	glog.V(3).Infof("获取「%s」信息开始...\n", opt.Name)
	mikanID := b.parseMikan1(opt.Url)
	bangumiID := b.parseMikan2(mikanID)
	ep := 1
	// TODO: opt.Name为文件名，解析出ep数
	info := b.parseBangumi(bangumiID, ep, opt.Date)
	if info == nil {
		glog.Errorln("获取Bangumi信息失败，结束此流程")
		return nil
	}
	info.BangumiSeason = b.parseThemoviedb(info.Name, info.AirDate)
	if info.Season == 0 {
		glog.Errorln("获取Themoviedb季度信息失败，默认SE01")
		info.Season = 1
	}
	info.BangumiExtra = &models.BangumiExtra{
		SubID:  mikanID,
		SubUrl: opt.Url,
	}
	glog.V(3).Infof("获取「%s」信息成功！更名为「%s」\n", opt.Name, info.FullName())
	return info
}

// parseMikan1
//  @Description: 解析mikan rss中的link页面，获取当前资源的mikan id
//  @receiver b
//  @param url_
//  @return int
//
func (b *Mikan) parseMikan1(url_ string) int {
	glog.V(5).Infof("步骤1，解析Mikan，%s\n", url_)
	doc, err := htmlquery.LoadURL(url_)
	if err != nil {
		glog.Errorln(err)
		return 0
	}
	miaknLink := htmlquery.FindOne(doc, MikanIdXPath)
	href := htmlquery.SelectAttr(miaknLink, "href")
	u, err := url.Parse(href)
	if err != nil {
		glog.Errorln(err)
		return 0
	}
	query := u.Query()
	if query.Has("bangumiId") {
		id, err := strconv.Atoi(query.Get("bangumiId"))
		if err != nil {
			glog.Errorln(err)
			return 0
		}
		return id
	}
	glog.Errorln("获取Mikan ID失败")
	return 0
}

// parseMikan2
//  @Description: 通过mikan id解析mikan番剧信息页面，获取bgm.tv id
//  @receiver b
//  @param mikanID
//  @return int
//
func (b *Mikan) parseMikan2(mikanID int) int {
	url_ := MikanInfoUrl(mikanID)
	glog.V(5).Infof("步骤2，解析Mikan，%s\n", url_)
	doc, err := htmlquery.LoadURL(url_)
	if err != nil {
		glog.Errorln(err)
		return 0
	}
	bangumiUrl := htmlquery.FindOne(doc, MikanBangumiUrlXPath)
	href := htmlquery.SelectAttr(bangumiUrl, "href")

	//fmt.Println(href)
	hrefSplit := strings.Split(href, "/")
	bgmId, err := strconv.Atoi(hrefSplit[len(hrefSplit)-1])
	if err != nil {
		glog.Errorln(err)
		return 0
	}
	return bgmId
}

//  parseBangumi
//  @Description: 解析bgm.tv，获取番剧信息和ep信息
//  @receiver b
//  @param bangumiID
//  @param ep
//  @return *models.Bangumi
//
func (b *Mikan) parseBangumi(bangumiID, ep int, date string) *models.Bangumi {
	glog.V(5).Infof("步骤3，解析Bangumi，%d\n", bangumiID)
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
	glog.V(5).Infof("步骤4，解析Themoviedb，%s\n", name)
	tmdb := NewThemoviedb()
	newBgm := tmdb.Parse(&models.BangumiParseOptions{
		Name: name,
		Date: airDate,
	})
	return newBgm.BangumiSeason
}
