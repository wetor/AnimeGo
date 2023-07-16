package api

import (
	"github.com/gin-gonic/gin"
	webModels "github.com/wetor/AnimeGo/internal/web/models"
)

func (a *Api) AddItems(c *gin.Context) {
	var request webModels.AddItemsRequest
	if !a.checkRequest(c, &request) {
		return
	}
}

func (a *Api) ListItems(c *gin.Context) {

}

func (a *Api) DeleteItems(c *gin.Context) {

}
