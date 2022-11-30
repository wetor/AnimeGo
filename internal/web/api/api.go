package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/internal/utils"
	webModels "github.com/wetor/AnimeGo/internal/web/models"
	"github.com/wetor/AnimeGo/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// checkRequest 绑定request结构体
//  Description 包含登录信息
//  Param c *gin.Context
//  Param data any 返回结构体指针
//  Return bool
//
func checkRequest(c *gin.Context, data any) bool {
	if err := c.ShouldBind(data); err != nil {
		zap.S().Warnf("参数错误，err: %s", errors.NewAniErrorD(err))
		c.JSON(webModels.Fail("参数错误"))
		return false
	}
	if err := c.ShouldBindQuery(data); err != nil {
		zap.S().Warnf("Query参数错误，err: %s", errors.NewAniErrorD(err))
		c.JSON(webModels.Fail("Query参数错误"))
		return false
	}
	if err := c.ShouldBindUri(data); err != nil {
		zap.S().Warnf("Uri参数错误，err: %s", errors.NewAniErrorD(err))
		c.JSON(webModels.Fail("Uri参数错误"))
		return false
	}

	key, has := c.Get("access_key")
	localKey := utils.Sha256(store.Config.WebApi.AccessKey)
	if has && key != localKey {
		zap.S().Warn(errors.NewAniError("Access key错误！"))
		c.JSON(webModels.Fail("Access key错误"))
		return false
	}
	return true
}

// Ping godoc
//  @Summary Ping
//  @Description Pong
//  @Tags web
//  @Accept  json
//  @Produce  json
//  @Success 200 {object} webModels.Response
//  @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(webModels.Succ("pong", gin.H{
		"time": time.Now().Unix(),
	}))
}

// SHA256 godoc
//  @Summary SHA256计算
//  @Description SHA256计算
//  @Tags web
//  @Accept  json
//  @Produce  json
//  @Param access_key query string true "原文本"
//  @Success 200 {object} webModels.Response{data=string}
//  @Router /sha256 [get]
func SHA256(c *gin.Context) {
	c.JSON(webModels.Succ("Access-Key", utils.Sha256(c.Query("access_key"))))
}
