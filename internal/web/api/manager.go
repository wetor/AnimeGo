package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/wetor/AnimeGo/internal/models"
	webModels "github.com/wetor/AnimeGo/internal/web/models"
)

// AddItems godoc
//
//	@Summary		添加下载项
//	@Description	添加下载项到AnimeGo
//	@Tags			manager
//	@Accept			json
//	@Produce		json
//	@Param			data	body		webModels.AddItemsRequest	true	"下载项信息"
//	@Success		200		{object}	webModels.Response
//	@Failure		300		{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/manager [post]
func (a *Api) AddItems(c *gin.Context) {
	var request webModels.AddItemsRequest
	if !a.checkRequest(c, &request) {
		return
	}
	source := strings.ToLower(request.Source)

	items := make([]*models.FeedItem, 0, len(request.Data))
	for _, data := range request.Data {
		item := &models.FeedItem{
			Download: data.Torrent,
		}
		switch source {
		case "mikan":
			if name, ok := data.Info["name"]; ok {
				item.Name = name.(string)
			} else {
				c.JSON(webModels.Fail(source + " 源缺少 info.name 参数"))
				return
			}
			if url, ok := data.Info["url"]; ok {
				item.Url = url.(string)
			} else {
				c.JSON(webModels.Fail(source + " 源缺少 info.url 参数"))
				return
			}
		}
		items = append(items, item)
	}
	err := a.filterManager.Update(a.ctx, items, true, true)
	if err != nil {
		c.JSON(webModels.Fail(err.Error()))
		return
	}
	c.JSON(webModels.Succ(fmt.Sprintf("开始处理%d个下载项", len(items))))
}

func (a *Api) ListItems(c *gin.Context) {

}

func (a *Api) DeleteItems(c *gin.Context) {

}
