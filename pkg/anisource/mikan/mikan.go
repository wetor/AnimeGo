package mikan

import (
	"AnimeGo/pkg/anisource"
	mem "AnimeGo/pkg/memorizer"
	"AnimeGo/pkg/request"
	"bytes"
	"encoding/gob"
	"fmt"
	"golang.org/x/net/html"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
)

const (
	IdXPath         = "//a[@class='mikan-rss']"                                 // Mikan番剧id获取XPath
	GroupXPath      = "//p[@class='bangumi-info']/a[@class='magnet-link-wrap']" // Mikan番剧信息获取group字幕组id和name
	BangumiUrlXPath = "//p[@class='bangumi-info']/a[contains(@href, 'bgm.tv')]" // Mikan番剧信息中bangumi id获取XPath
)

var (
	Host              = "https://mikanani.me"
	Bucket            = "mikan"
	CacheSecond int64 = 30 * 24 * 60 * 60
)

type Mikan struct {
	cacheInit                   bool
	cacheParseMikanInfoVar      mem.Func
	cacheParseMikanBangumiIDVar mem.Func
}

func (m *Mikan) RegisterCache() {
	if anisource.Cache == nil {
		panic("需要先调用anisource.Init初始化缓存")
	}
	m.cacheInit = true
	m.cacheParseMikanInfoVar = mem.Memorized(Bucket, anisource.Cache, func(params *mem.Params, results *mem.Results) error {
		mikan, err := m.parseMikanInfo(params.Get("mikanUrl").(string))
		if err != nil {
			return err
		}
		results.Set("mikanInfo", mikan)
		return nil
	})

	m.cacheParseMikanBangumiIDVar = mem.Memorized(Bucket, anisource.Cache, func(params *mem.Params, results *mem.Results) error {
		bangumiID, err := m.parseMikanBangumiID(params.Get("mikanID").(int))
		if err != nil {
			return err
		}
		results.Set("bangumiID", bangumiID)
		return nil
	})
}

func (m Mikan) ParseCache(url string) (mikanID int, bangumiID int, err error) {
	mikan, err := m.CacheParseMikanInfo(url)
	fmt.Println(mikan)
	if err != nil {
		return
	}
	mikanID = mikan.ID
	bangumiID, err = m.CacheParseMikanBangumiID(mikanID)
	return
}

func (m Mikan) CacheParseMikanInfo(url string) (mikanInfo *MikanInfo, err error) {
	if !m.cacheInit {
		m.RegisterCache()
	}
	results := mem.NewResults("mikanInfo", &MikanInfo{})
	err = m.cacheParseMikanInfoVar(mem.NewParams("mikanUrl", url).TTL(CacheSecond), results)
	if err != nil {
		return
	}
	mikanInfo = results.Get("mikanInfo").(*MikanInfo)
	return
}

func (m Mikan) CacheParseMikanBangumiID(mikanID int) (bangumiID int, err error) {
	if !m.cacheInit {
		m.RegisterCache()
	}
	results := mem.NewResults("bangumiID", 0)
	err = m.cacheParseMikanBangumiIDVar(mem.NewParams("mikanID", mikanID).TTL(CacheSecond), results)
	if err != nil {
		return
	}
	bangumiID = results.Get("bangumiID").(int)
	return
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
	mikan, err := m.parseMikanInfo(url)
	if err != nil {
		return 0, 0, err
	}
	bangumiID, err = m.parseMikanBangumiID(mikan.ID)
	if err != nil {
		return mikan.ID, 0, err
	}
	return mikan.ID, bangumiID, nil
}

func (m Mikan) loadHtml(url string) (*html.Node, error) {
	buf := bytes.NewBuffer(nil)
	err := request.Get(&request.Param{
		Uri:     url,
		Proxy:   anisource.Proxy,
		Writer:  buf,
		Retry:   anisource.Retry,
		Timeout: anisource.Timeout,
	})
	if err != nil {
		return nil, err
	}
	doc, err := htmlquery.Parse(buf)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// parseMikanID
//  @Description: 解析网页取出mikan的id、group等信息
//  @receiver Mikan
//  @param mikanUrl string
//  @return mikan *MikanInfo
//  @return err error
//
func (m Mikan) parseMikanInfo(mikanUrl string) (mikan *MikanInfo, err error) {
	doc, err := m.loadHtml(mikanUrl)
	if err != nil {
		return
	}
	miaknLink := htmlquery.FindOne(doc, IdXPath)
	href := htmlquery.SelectAttr(miaknLink, "href")
	u, err := url.Parse(href)
	if err != nil {
		return
	}
	mikan = &MikanInfo{}
	query := u.Query()
	if query.Has("bangumiId") {
		mikan.ID, err = strconv.Atoi(query.Get("bangumiId"))
		if err != nil {
			return
		}
		mikan.SubGroupID, err = strconv.Atoi(query.Get("subgroupid"))
		if err != nil {
			return
		}
	} else {
		return nil, ParseMikanIDErr
	}

	group := htmlquery.FindOne(doc, GroupXPath)
	href = htmlquery.SelectAttr(group, "href")
	_, groupId := path.Split(href)
	mikan.PubGroupID, err = strconv.Atoi(groupId)
	if err != nil {
		return
	}
	mikan.GroupName = group.FirstChild.Data
	return
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
	doc, err := m.loadHtml(url_)
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

func init() {
	gob.Register(&MikanInfo{})
}
