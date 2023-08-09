package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/wetor/AnimeGo/internal/animego/feed/rss"
	"github.com/wetor/AnimeGo/internal/models"
	webModels "github.com/wetor/AnimeGo/internal/web/models"
)

// Rss godoc
//
//	@Summary		发送下载项
//	@Description	将待下载项组合成rss发送给AnimeGo
//	@Tags			rss
//	@Accept			json
//	@Produce		json
//	@Param			rss	body		webModels.SelectEpRequest	true	"组合的rss信息"
//	@Success		200	{object}	webModels.Response
//	@Failure		300	{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/rss [post]
func (a *Api) Rss(c *gin.Context) {
	var request webModels.SelectEpRequest
	if !a.checkRequest(c, &request) {
		return
	}
	reqRss := rss.NewRss(&rss.Options{Url: request.Rss.Url})

	items, err := reqRss.Parse()
	if err != nil {
		c.JSON(webModels.Fail(err.Error()))
		return
	}
	if request.IsSelectEp {
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
	err = a.filterManager.Update(a.ctx, items, false, true)
	if err != nil {
		c.JSON(webModels.Fail(err.Error()))
		return
	}
	c.JSON(webModels.Succ(fmt.Sprintf("开始处理%d个下载项", len(items))))
}
