package api

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/internal/constant"
	webModels "github.com/wetor/AnimeGo/internal/web/models"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
)

// ConfigFileGet godoc
//
//	@Summary		获取配置文件
//	@Description	获取AnimeGo的配置文件文本内容
//	@Tags			config
//	@Accept			json
//	@Produce		plain
//	@Success		200	{object}	webModels.Response
//	@Failure		300	{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/config/file [get]
func (a *Api) ConfigFileGet(c *gin.Context) {
	data, err := os.ReadFile(configs.ConfigFile)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail("打开配置文件失败"))
		return
	}
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.String(http.StatusOK, string(data))
}

// ConfigFilePut godoc
//
//	@Summary		写入配置文件
//	@Description	写入AnimeGo的配置文件文本内容
//	@Tags			config
//	@Accept			plain
//	@Produce		json
//	@Param			backup	query		bool	true	"备份旧的配置"
//	@Param			config	body		string	true	"配置文件文本"
//	@Success		200		{object}	webModels.Response
//	@Failure		300		{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/config/file [put]
func (a *Api) ConfigFilePut(c *gin.Context) {
	var request webModels.ConfigFilePutRequest
	if !a.checkRequest(c, &request) {
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail("读取文件参数失败"))
	}
	defer c.Request.Body.Close()

	if request.Backup {
		err = configs.BackupConfig(configs.ConfigFile, "web-backup")
		if err != nil {
			log.DebugErr(err)
			c.JSON(webModels.Fail("备份配置文件失败"))
			return
		}
	}

	err = os.WriteFile(configs.ConfigFile, body, constant.FilePerm)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail("写入配置文件失败"))
		return
	}
	// TODO 重载配置
	c.JSON(webModels.Succ("更新成功，需要重启AnimeGo以应用配置"))
}

// ConfigGet godoc
//
//	@Summary		获取设置
//	@Description	获取AnimeGo的配置文件内容
//	@Tags			config
//	@Accept			json
//	@Produce		json
//	@Param			type	query		webModels.ConfigGetRequest	true	"获取配置文件"
//	@Success		200		{object}	webModels.Response
//	@Failure		300		{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/config [get]
func (a *Api) ConfigGet(c *gin.Context) {
	var request webModels.ConfigGetRequest
	if !a.checkRequest(c, &request) {
		return
	}
	if request.Key == "all" {
		c.JSON(webModels.Succ("配置项值", a.config))
	} else if request.Key == "default" {
		c.JSON(webModels.Succ("配置项默认值", configs.DefaultConfig()))
	} else if request.Key == "comment" {
		data := make(map[string]interface{})
		err := json.Unmarshal(configs.DefaultDoc(), &data)
		if err != nil {
			log.DebugErr(err)
			c.JSON(webModels.Fail("配置项说明格式化错误"))
			return
		}
		c.JSON(webModels.Succ("配置项说明", data))
	} else if request.Key == "raw" {
		data, err := os.ReadFile(configs.ConfigFile)
		if err != nil {
			log.DebugErr(err)
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
//
//	@Summary		更新设置
//	@Description	更新AnimeGo的配置文件内容
//	@Tags			config
//	@Accept			json
//	@Produce		json
//	@Param			type	body		webModels.ConfigPutRequest	true	"更新配置文件"
//	@Success		200		{object}	webModels.Response
//	@Failure		300		{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/config [put]
func (a *Api) ConfigPut(c *gin.Context) {
	var request webModels.ConfigPutRequest
	if !a.checkRequest(c, &request) {
		return
	}
	var err error
	if request.Key == "all" || request.Key == "raw" {
		var data []byte
		if request.Key == "all" && request.Config != nil {
			data, err = configs.Config2Bytes(request.Config)
		} else if request.Key == "raw" && request.ConfigRaw != nil {
			data, err = base64.StdEncoding.DecodeString(*request.ConfigRaw)
		} else {
			c.JSON(webModels.Fail("参数错误，未传入对应数据"))
			return
		}
		if err != nil {
			log.DebugErr(err)
			c.JSON(webModels.Fail("参数格式错误"))
			return
		}

		if *request.Backup {
			err = configs.BackupConfig(configs.ConfigFile, "")
			if err != nil {
				log.DebugErr(err)
				c.JSON(webModels.Fail("备份文件 " + configs.ConfigFile + " 失败"))
				return
			}
		}
		err = os.WriteFile(configs.ConfigFile, data, constant.FilePerm)
		if err != nil {
			log.DebugErr(err)
			c.JSON(webModels.Fail("写到文件 " + configs.ConfigFile + " 失败"))
			return
		}
		c.JSON(webModels.Succ("更新成功，需要重启AnimeGo以应用配置"))
	} else {
		c.JSON(webModels.Fail("暂不支持 " + request.Key + "，目前仅支持 'all', 'raw'"))
	}
}
