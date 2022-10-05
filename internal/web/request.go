package web

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// checkRequest 绑定request结构体
//  Description 包含登录信息
//  Param c *gin.Context
//  Param data interface{} 返回结构体指针
//  Return bool
//
func checkRequest(c *gin.Context, data interface{}) bool {
	if err := c.ShouldBind(data); err != nil {
		zap.S().Error("参数错误\t", err)
		c.JSON(Fail("参数错误"))
		return false
	}
	if err := c.ShouldBindQuery(data); err != nil {
		zap.S().Error("Query参数错误\t", err)
		c.JSON(Fail("Query参数错误"))
		return false
	}
	if err := c.ShouldBindUri(data); err != nil {
		zap.S().Error("Uri参数错误\t", err)
		c.JSON(Fail("Uri参数错误"))
		return false
	}
	return true
}

type SelectEpRequest struct {
	Source string `json:"source"`
	Rss    struct {
		Url string `json:"url"`
	} `json:"rss"`
	SelectEp bool     `json:"select_ep"`
	EpLinks  []string `json:"ep_links"`
}
