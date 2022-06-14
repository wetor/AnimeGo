package bangumi

import (
	"GoBangumi/model"
	"fmt"
	"github.com/antchfx/htmlquery"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestNewMikan(t *testing.T) {
	bgm := model.Bangumi{}

	// ------------  获取Mikan ID
	filePath := "../../data/cache/1.html"
	doc, err := htmlquery.LoadDoc(filePath)
	if err != nil {
		panic(err)
	}
	mikan_rss_a := htmlquery.FindOne(doc, "//a[@class='mikan-rss']")
	href := htmlquery.SelectAttr(mikan_rss_a, "href")
	u, err := url.Parse(href)
	if err != nil {
		panic(err)
	}

	query := u.Query()
	if query.Has("bangumiId") {
		id, err := strconv.Atoi(query.Get("bangumiId"))
		if err != nil {
			panic(err)
		}
		bgm.SubID = id
	}
	fmt.Println(bgm)
	// ------------  获取bgm ID
	filePath = fmt.Sprintf("../../data/cache/%d.html", bgm.SubID)
	doc, err = htmlquery.LoadDoc(filePath)
	if err != nil {
		panic(err)
	}
	bangumiUrl := htmlquery.FindOne(doc, "//p[@class='bangumi-info']/a[contains(@href, 'bgm.tv')]")
	href = htmlquery.SelectAttr(bangumiUrl, "href")

	fmt.Println(href)
	hrefSplit := strings.Split(href, "/")
	bgmId, err := strconv.Atoi(hrefSplit[len(hrefSplit)-1])
	if err != nil {
		panic(err)
	}
	bgm.ID = bgmId

	fmt.Println(bgm)
}
