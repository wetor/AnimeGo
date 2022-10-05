package web

import (
	"AnimeGo/internal/animego/feed/mikan"
	"AnimeGo/internal/models"
	"AnimeGo/internal/store"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"time": time.Now(),
		"pong": true,
	})
}

func Download(c *gin.Context) {
	var request SelectEpRequest
	if !checkRequest(c, &request) {
		return
	}
	rss := mikan.NewRss(request.Rss.Url, "")

	items := rss.Parse()
	if request.SelectEp {
		set := make(map[string]bool)
		for _, item := range request.EpLinks {
			set[item] = true
		}
		selectItems := make([]*models.FeedItem, 0, len(request.EpLinks))
		for _, item := range items {
			if _, has := set[item.Url]; has {
				selectItems = append(selectItems, item)
			}
		}
		items = selectItems
	}
	store.Process.UpdateFeed(items)
	c.JSON(Succ(fmt.Sprintf("开始处理%d个下载项", len(items))))
}
