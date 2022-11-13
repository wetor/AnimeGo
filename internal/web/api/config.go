package api

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/internal/store"
	webModels "github.com/wetor/AnimeGo/internal/web/models"
	"github.com/wetor/AnimeGo/pkg/errors"
	"go.uber.org/zap"
	"os"
)

// Config godoc
//  @Summary 获取设置
//  @Description 将待下载项组合成rss发送给AnimeGo
//  @Tags config
//  @Accept  json
//  @Produce  json
//  @Param type query webModels.ConfigGetRequest true "获取设置内容"
//  @Success 200 {object} webModels.Response{data=webModels.ConfigResponse}
//  @Failure 300 {object} webModels.Response
//  @Security ApiKeyAuth
//  @Router /api/config [get]
func Config(c *gin.Context) {
	var request webModels.ConfigGetRequest
	if !checkRequest(c, &request) {
		return
	}
	if request.Key == "all" {
		resp := &webModels.ConfigResponse{
			Config: store.Config,
		}
		c.JSON(webModels.Succ("配置项", resp))
	} else if request.Key == "raw" {
		data, err := os.ReadFile(configs.ConfigFile)
		if err != nil {
			err = errors.NewAniErrorD(err)
			zap.S().Debug(err)
			c.JSON(webModels.Fail("打开文件 " + configs.ConfigFile + " 失败"))
			return
		}
		str := base64.StdEncoding.EncodeToString(data)
		c.JSON(webModels.Succ("配置文件", webModels.ConfigResponse{
			Data: str,
		}))
	} else if request.Key == "comment" {
		str := base64.StdEncoding.EncodeToString(configs.DefaultDoc())
		c.JSON(webModels.Succ("配置项说明", webModels.ConfigResponse{
			Data: str,
		}))
	} else {
		c.JSON(webModels.Fail("暂不支持 " + request.Key + "，目前仅支持 all 和 raw"))
	}
}
