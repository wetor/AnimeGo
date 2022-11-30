package mikan

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/pkg/errors"
	mem "github.com/wetor/AnimeGo/pkg/memorizer"
	"github.com/wetor/AnimeGo/pkg/request"
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
	Host   = "https://mikanani.me"
	Bucket = "mikan"
)

type Mikan struct {
	cacheInit                   bool
	cacheParseMikanInfoVar      mem.Func
	cacheParseMikanBangumiIDVar mem.Func
}

func (m *Mikan) RegisterCache() {
	if anidata.Cache == nil {
		errors.NewAniError("需要先调用anidata.Init初始化缓存").TryPanic()
	}
	m.cacheInit = true
	m.cacheParseMikanInfoVar = mem.Memorized(Bucket, anidata.Cache, func(params *mem.Params, results *mem.Results) error {
		mikan := m.parseMikanInfo(params.Get("mikanUrl").(string))
		results.Set("mikanInfo", mikan)
		return nil
	})

	m.cacheParseMikanBangumiIDVar = mem.Memorized(Bucket, anidata.Cache, func(params *mem.Params, results *mem.Results) error {
		bangumiID := m.parseMikanBangumiID(params.Get("mikanID").(int))
		results.Set("bangumiID", bangumiID)
		return nil
	})
}

func (m Mikan) ParseCache(url string) (mikanID int, bangumiID int) {
	mikanID = m.CacheParseMikanInfo(url).ID
	bangumiID = m.CacheParseMikanBangumiID(mikanID)
	return
}

func (m Mikan) CacheParseMikanInfo(url string) (mikanInfo *MikanInfo) {
	if !m.cacheInit {
		m.RegisterCache()
	}
	results := mem.NewResults("mikanInfo", &MikanInfo{})
	err := m.cacheParseMikanInfoVar(mem.NewParams("mikanUrl", url).
		TTL(anidata.CacheTime[Bucket]), results)
	errors.NewAniErrorD(err).TryPanic()
	mikanInfo = results.Get("mikanInfo").(*MikanInfo)
	return
}

func (m Mikan) CacheParseMikanBangumiID(mikanID int) (bangumiID int) {
	if !m.cacheInit {
		m.RegisterCache()
	}
	results := mem.NewResults("bangumiID", 0)
	err := m.cacheParseMikanBangumiIDVar(mem.NewParams("mikanID", mikanID).
		TTL(anidata.CacheTime[Bucket]), results)
	errors.NewAniErrorD(err).TryPanic()
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
func (m Mikan) Parse(url string) (mikanID int, bangumiID int) {
	mikan := m.parseMikanInfo(url)
	mikanID = mikan.ID
	bangumiID = m.parseMikanBangumiID(mikan.ID)
	return
}

func (m Mikan) loadHtml(url string) *html.Node {
	buf := bytes.NewBuffer(nil)
	err := request.GetWriter(url, buf)
	errors.NewAniErrorD(err).TryPanic()
	doc, err := htmlquery.Parse(buf)
	errors.NewAniErrorD(err).TryPanic()

	return doc
}

// parseMikanID
//  @Description: 解析网页取出mikan的id、group等信息
//  @receiver Mikan
//  @param mikanUrl string
//  @return mikan *MikanInfo
//
func (m Mikan) parseMikanInfo(mikanUrl string) (mikan *MikanInfo) {
	doc := m.loadHtml(mikanUrl)

	miaknLink := htmlquery.FindOne(doc, IdXPath)
	href := htmlquery.SelectAttr(miaknLink, "href")
	u, err := url.Parse(href)
	errors.NewAniErrorD(err).TryPanic()

	mikan = &MikanInfo{}
	query := u.Query()
	if query.Has("bangumiId") {
		mikan.ID, err = strconv.Atoi(query.Get("bangumiId"))
		errors.NewAniErrorD(err).TryPanic()

		mikan.SubGroupID, err = strconv.Atoi(query.Get("subgroupid"))
		if err != nil {
			mikan.SubGroupID = 0
		}
	} else {
		errors.NewAniError("解析Bangumi ID失败").TryPanic()
	}

	group := htmlquery.FindOne(doc, GroupXPath)
	if group == nil {
		return
	}
	href = htmlquery.SelectAttr(group, "href")
	_, groupId := path.Split(href)
	mikan.PubGroupID, err = strconv.Atoi(groupId)
	errors.NewAniErrorD(err).TryPanic()

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
func (m Mikan) parseMikanBangumiID(mikanID int) (bangumiID int) {
	url_ := fmt.Sprintf("%s/Home/bangumi/%d", Host, mikanID)
	doc := m.loadHtml(url_)

	bangumiUrl := htmlquery.FindOne(doc, BangumiUrlXPath)
	href := htmlquery.SelectAttr(bangumiUrl, "href")

	hrefSplit := strings.Split(href, "/")
	bangumiID, err := strconv.Atoi(hrefSplit[len(hrefSplit)-1])
	errors.NewAniErrorD(err).TryPanic()

	return bangumiID
}

func init() {
	gob.Register(&MikanInfo{})
}
