package mikan

import (
	"GoBangumi/pkg/anisource"
	mem "GoBangumi/pkg/memorizer"
	"fmt"
	"github.com/antchfx/htmlquery"
	"net/url"
	"strconv"
	"strings"
)

const (
	IdXPath         = "//a[@class='mikan-rss']"                                 // Mikan番剧id获取XPath
	BangumiUrlXPath = "//p[@class='bangumi-info']/a[contains(@href, 'bgm.tv')]" // Mikan番剧信息中bangumi id获取XPath
)

var (
	Host              = "https://mikanani.me"
	Bucket            = "mikan"
	CacheSecond int64 = 30 * 24 * 60 * 60
)

type Mikan struct {
	cacheInit                bool
	cacheParseMikanID        mem.Func
	cacheparseMikanBangumiID mem.Func
}

func (m *Mikan) RegisterCache() {
	if anisource.Cache == nil {
		panic("需要先调用anisource.Init初始化缓存")
	}
	m.cacheInit = true
	m.cacheParseMikanID = mem.Memorized(Bucket, anisource.Cache, func(params *mem.Params, results *mem.Results) error {
		mikanID, err := m.parseMikanID(params.Get("mikanUrl").(string))
		if err != nil {
			return err
		}
		results.Set("mikanID", mikanID)
		return nil
	})

	m.cacheparseMikanBangumiID = mem.Memorized(Bucket, anisource.Cache, func(params *mem.Params, results *mem.Results) error {
		bangumiID, err := m.parseMikanBangumiID(params.Get("mikanID").(int))
		if err != nil {
			return err
		}
		results.Set("bangumiID", bangumiID)
		return nil
	})
}

func (m Mikan) ParseCache(url string) (mikanID int, bangumiID int, err error) {
	if !m.cacheInit {
		m.RegisterCache()
	}
	results := mem.NewResults("mikanID", 0, "bangumiID", 0)

	err = m.cacheParseMikanID(mem.NewParams("mikanUrl", url).TTL(CacheSecond), results)
	if err != nil {
		return 0, 0, err
	}
	mikanID = results.Get("mikanID").(int)

	err = m.cacheparseMikanBangumiID(mem.NewParams("mikanID", mikanID).TTL(CacheSecond), results)
	if err != nil {
		return mikanID, 0, err
	}
	bangumiID = results.Get("bangumiID").(int)
	return mikanID, bangumiID, nil
}

// Parse
//  @Description: 通过mikan剧集的url，解析两次网页，分别获取到mikanID和bangumiID
//  @receiver Mikan
//  @param url string mikan剧集的url
//  @return mikanID int
//  @return bangumiID int
//  @return err error
//
func (m Mikan) Parse(url string) (mikanID int, bangumiID int, err error) {
	mikanID, err = m.parseMikanID(url)
	if err != nil {
		return 0, 0, err
	}
	bangumiID, err = m.parseMikanBangumiID(mikanID)
	if err != nil {
		return mikanID, 0, err
	}
	return mikanID, bangumiID, nil
}

// parseMikanID
//  @Description: 解析网页取出mikanID
//  @receiver Mikan
//  @param mikanUrl string
//  @return mikanID int
//  @return err error
//
func (m Mikan) parseMikanID(mikanUrl string) (mikanID int, err error) {
	doc, err := htmlquery.LoadURL(mikanUrl)
	if err != nil {
		return 0, err
	}
	miaknLink := htmlquery.FindOne(doc, IdXPath)
	href := htmlquery.SelectAttr(miaknLink, "href")
	u, err := url.Parse(href)
	if err != nil {
		return 0, err
	}
	query := u.Query()
	if query.Has("bangumiId") {
		mikanID, err = strconv.Atoi(query.Get("bangumiId"))
		if err != nil {
			return 0, err
		}
	} else {
		return 0, ParseMikanIDErr
	}
	return mikanID, nil
}

// parseMikanBangumiID
//  @Description: 解析网页取出bangumiID
//  @receiver Mikan
//  @param mikanID int
//  @return bangumiID int
//  @return err error
//
func (m Mikan) parseMikanBangumiID(mikanID int) (bangumiID int, err error) {
	url_ := fmt.Sprintf("%s/Home/bangumi/%d", Host, mikanID)
	doc, err := htmlquery.LoadURL(url_)
	if err != nil {
		return 0, err
	}
	bangumiUrl := htmlquery.FindOne(doc, BangumiUrlXPath)
	href := htmlquery.SelectAttr(bangumiUrl, "href")

	hrefSplit := strings.Split(href, "/")
	bangumiID, err = strconv.Atoi(hrefSplit[len(hrefSplit)-1])
	if err != nil {
		return 0, err
	}
	return bangumiID, nil
}
