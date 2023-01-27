package api

import (
	"encoding/base64"
	"os"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/configs"
	webModels "github.com/wetor/AnimeGo/internal/web/models"
	"github.com/wetor/AnimeGo/pkg/errors"
)

// ConfigGet godoc
//  @Summary 获取设置
//  @Description 获取AnimeGo的配置文件内容
//  @Tags config
//  @Accept  json
//  @Produce  json
//  @Param type query webModels.ConfigGetRequest true "获取配置文件"
//  @Success 200 {object} webModels.Response
//  @Failure 300 {object} webModels.Response
//  @Security ApiKeyAuth
//  @Router /api/config [get]
func ConfigGet(c *gin.Context) {
	var request webModels.ConfigGetRequest
	if !checkRequest(c, &request) {
		return
	}
	if request.Key == "all" {
		c.JSON(webModels.Succ("配置项值", Config))
	} else if request.Key == "default" {
		c.JSON(webModels.Succ("配置项默认值", configs.DefaultConfig()))
	} else if request.Key == "comment" {
		data := make(map[string]interface{})
		err := jsoniter.Unmarshal(configs.DefaultDoc(), &data)
		if err != nil {
			err = errors.NewAniErrorD(err)
			zap.S().Debug(err)
			c.JSON(webModels.Fail("配置项说明格式化错误"))
			return
		}
		c.JSON(webModels.Succ("配置项说明", data))
	} else if request.Key == "raw" {
		data, err := os.ReadFile(configs.ConfigFile)
		if err != nil {
			err = errors.NewAniErrorD(err)
			zap.S().Debug(err)
			c.JSON(webModels.Fail("打开文件 " + configs.ConfigFile + " 失败"))
			return
		}
		str := base64.StdEncoding.EncodeToString(data)
		c.JSON(webModels.Succ("配置文件", str))
	} else {
		c.JSON(webModels.Fail("暂不支持 " + request.Key + "，目前仅支持 'all', 'default', 'comment', 'raw'"))
	}
}

// ConfigPut godoc
//  @Summary 更新设置
//  @Description 更新AnimeGo的配置文件内容
//  @Tags config
//  @Accept  json
//  @Produce  json
//  @Param type body webModels.ConfigPutRequest true "更新配置文件"
//  @Success 200 {object} webModels.Response
//  @Failure 300 {object} webModels.Response
//  @Security ApiKeyAuth
//  @Router /api/config [put]
func ConfigPut(c *gin.Context) {
	var request webModels.ConfigPutRequest
	if !checkRequest(c, &request) {
		return
	}
	if request.Key == "all" {
		if *request.Backup {
			err := configs.BackupConfig(configs.ConfigFile, "")
			if err != nil {
				err = errors.NewAniErrorD(err)
				zap.S().Debug(err)
				c.JSON(webModels.Fail("备份文件 " + configs.ConfigFile + " 失败"))
				return
			}
		}
		err := os.WriteFile(configs.ConfigFile, configs.Config2Bytes(request.Config), 0644)
		if err != nil {
			err = errors.NewAniErrorD(err)
			zap.S().Debug(err)
			c.JSON(webModels.Fail("写到文件 " + configs.ConfigFile + " 失败"))
			return
		}
		c.JSON(webModels.Succ("更新成功，需要重启AnimeGo以应用配置"))
	} else {
		c.JSON(webModels.Fail("暂不支持 " + request.Key + "，目前仅支持 all"))
	}
}
