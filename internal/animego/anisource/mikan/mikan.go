package mikan

import (
	"bytes"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"golang.org/x/net/html"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/pkg/request"
	"github.com/wetor/AnimeGo/pkg/log"
	mem "github.com/wetor/AnimeGo/pkg/memorizer"
)

var Set = wire.NewSet(
	NewMikan,
)

type Mikan struct {
	cacheInit                   bool
	cacheParseMikanInfoVar      mem.Func
	cacheParseMikanBangumiIDVar mem.Func

	*Options
}

func NewMikan(opts *Options) *Mikan {
	return &Mikan{
		Options: opts,
	}
}

func (a *Mikan) Name() string {
	return "Mikan"
}

func (a *Mikan) RegisterCache() {
	a.cacheInit = true
	a.cacheParseMikanInfoVar = mem.Memorized(constant.MikanBucket, a.Cache, func(params *mem.Params, results *mem.Results) error {
		mikan, err := a.parseMikanInfo(params.Get("mikanUrl").(string))
		if err != nil {
			return err
		}
		results.Set("mikanInfo", mikan)
		return nil
	})

	a.cacheParseMikanBangumiIDVar = mem.Memorized(constant.MikanBucket, a.Cache, func(params *mem.Params, results *mem.Results) error {
		bangumiID, err := a.parseMikanBangumiID(params.Get("mikanID").(int))
		if err != nil {
			return err
		}
		results.Set("bangumiID", bangumiID)
		return nil
	})
}

func (a *Mikan) ParseCache(url any) (entity any, err error) {
	mikan, err := a.CacheParseMikanInfo(url.(string))
	if err != nil {
		return nil, errors.Wrap(err, "解析Mikan信息失败")
	}
	bangumiID, err := a.cacheParseMikanBangumiID(mikan.(*MikanInfo).ID)
	if err != nil {
		return nil, errors.Wrap(err, "解析Mikan BangumiID失败")
	}
	return &Entity{
		MikanID:   mikan.(*MikanInfo).ID,
		BangumiID: bangumiID,
	}, nil
}

// Parse
//
//	@Description: 通过mikan剧集的url，解析两次网页，分别获取到mikanID和bangumiID
//	@receiver Mikan
//	@param url string mikan剧集的url
//	@return mikanID int
//	@return bangumiID int
func (a *Mikan) Parse(url any) (entity any, err error) {
	mikan, err := a.parseMikanInfo(url.(string))
	if err != nil {
		return nil, errors.Wrap(err, "解析Mikan信息失败")
	}
	bangumiID, err := a.parseMikanBangumiID(mikan.ID)
	if err != nil {
		return nil, errors.Wrap(err, "解析Mikan BangumiID失败")
	}
	return &Entity{
		MikanID:   mikan.ID,
		BangumiID: bangumiID,
	}, nil
}

func (a *Mikan) CacheParseMikanInfo(url string) (mikanInfo any, err error) {
	if !a.cacheInit {
		a.RegisterCache()
	}
	results := mem.NewResults("mikanInfo", &MikanInfo{})
	err = a.cacheParseMikanInfoVar(mem.NewParams("mikanUrl", url).
		TTL(a.CacheTime), results)
	if err != nil {
		return nil, err
	}
	mikanInfo = results.Get("mikanInfo").(*MikanInfo)
	return mikanInfo, nil
}

func (a *Mikan) cacheParseMikanBangumiID(mikanID int) (bangumiID int, err error) {
	if !a.cacheInit {
		a.RegisterCache()
	}
	results := mem.NewResults("bangumiID", 0)
	err = a.cacheParseMikanBangumiIDVar(mem.NewParams("mikanID", mikanID).
		TTL(a.CacheTime), results)
	if err != nil {
		return 0, err
	}
	bangumiID = results.Get("bangumiID").(int)
	return bangumiID, nil
}

func (a *Mikan) loadHtml(url string) (*html.Node, error) {
	buf := bytes.NewBuffer(nil)
	err := request.GetWriter(url, buf)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrRequest{Name: a.Name()})
	}
	doc, err := htmlquery.Parse(buf)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrMikanParseHTML{})
	}
	return doc, nil
}

// parseMikanID
//
//	@Description: 解析网页取出mikan的id、group等信息
//	@receiver Mikan
//	@paraa *MikanUrl string
//	@return mikan *MikanInfo
func (a *Mikan) parseMikanInfo(mikanUrl string) (mikan *MikanInfo, err error) {
	doc, err := a.loadHtml(mikanUrl)
	if err != nil {
		return nil, err
	}
	miaknLink := htmlquery.FindOne(doc, constant.MikanIdXPath)
	href := htmlquery.SelectAttr(miaknLink, "href")
	u, err := url.Parse(href)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrMikanParseHTML{Message: "MikanUrl"})
	}

	mikan = &MikanInfo{}
	// 解析url中的MikanID
	query := u.Query()
	if query.Has("bangumiId") {
		mikan.ID, err = strconv.Atoi(query.Get("bangumiId"))
		if err != nil {
			log.DebugErr(err)
			return nil, errors.WithStack(&exceptions.ErrMikanParseHTML{Message: "MikanID"})
		}
		mikan.SubGroupID, err = strconv.Atoi(query.Get("subgroupid"))
		if err != nil {
			mikan.SubGroupID = 0
		}
	} else {
		err = errors.WithStack(&exceptions.ErrMikanParseHTML{Message: "MikanID"})
		log.DebugErr(err)
		return nil, err
	}

	// 解析字幕组信息
	group := htmlquery.FindOne(doc, constant.MikanGroupXPath)
	if group != nil {
		href = htmlquery.SelectAttr(group, "href")
		_, groupId := path.Split(href)
		mikan.PubGroupID, err = strconv.Atoi(groupId)
		if err != nil {
			log.DebugErr(err)
			return nil, errors.WithStack(&exceptions.ErrMikanParseHTML{Message: "PubGroupID"})
		}
		mikan.GroupName = group.FirstChild.Data
	}
	return mikan, nil
}

// parseMikanBangumiID
//
//	@Description: 解析网页取出bangumiID
//	@receiver Mikan
//	@paraa *MikanID int
//	@return bangumiID int
func (a *Mikan) parseMikanBangumiID(mikanID int) (bangumiID int, err error) {
	url_ := fmt.Sprintf("%s/Home/bangumi/%d", constant.MikanHost, mikanID)
	doc, err := a.loadHtml(url_)
	if err != nil {
		return 0, err
	}

	bangumiUrl := htmlquery.FindOne(doc, constant.MikanBangumiUrlXPath)
	href := htmlquery.SelectAttr(bangumiUrl, "href")

	hrefSplit := strings.Split(href, "/")
	bangumiID, err = strconv.Atoi(hrefSplit[len(hrefSplit)-1])
	if err != nil {
		log.DebugErr(err)
		return 0, errors.WithStack(&exceptions.ErrMikanParseHTML{Message: "BangumiID"})
	}
	return bangumiID, nil
}

// Check interface is satisfied
var _ api.AniDataParse = &Mikan{}
