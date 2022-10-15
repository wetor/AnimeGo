package web

import (
	"AnimeGo/internal/store"
	"AnimeGo/internal/utils"
	"AnimeGo/pkg/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"path"
	"strings"
)

// checkRequest 绑定request结构体
//  Description 包含登录信息
//  Param c *gin.Context
//  Param data interface{} 返回结构体指针
//  Return bool
//
func checkRequest(c *gin.Context, data interface{}) bool {
	if err := c.ShouldBind(data); err != nil {
		zap.S().Warnf("参数错误，err: %s", errors.NewAniError(err.Error()))
		c.JSON(Fail("参数错误"))
		return false
	}
	if err := c.ShouldBindQuery(data); err != nil {
		zap.S().Warnf("Query参数错误，err: %s", errors.NewAniError(err.Error()))
		c.JSON(Fail("Query参数错误"))
		return false
	}
	if err := c.ShouldBindUri(data); err != nil {
		zap.S().Warnf("Uri参数错误，err: %s", errors.NewAniError(err.Error()))
		c.JSON(Fail("Uri参数错误"))
		return false
	}

	key, has := c.Get("access_key")
	localKey := utils.Sha256(store.Config.WebApi.AccessKey)
	if has && key != localKey {
		zap.S().Warn(errors.NewAniError("Access key错误！"))
		c.JSON(Fail("Access key错误"))
		return false
	}
	return true
}

type SelectEpRequest struct {
	Source string `json:"source" binding:"required"`
	Rss    struct {
		Url string `json:"url" binding:"required"`
	} `json:"rss" binding:"required"`
	IsSelectEp bool     `json:"is_select_ep" binding:"required"`
	EpLinks    []string `json:"ep_links"`
}

type PluginRequest struct {
	Name string `json:"name" form:"name" binding:"required"` //插件文件名
}

// CheckName 检查插件名
func (p PluginRequest) CheckName() error {
	dir, name := path.Split(p.Name)
	if len(dir) == 0 {
		return errors.NewAniError("插件子目录不能为空")
	}
	if len(name) == 0 || path.Ext(name) != ".js" {
		return errors.NewAniError("插件文件名错误")
	}
	if strings.Index(p.Name, "..") >= 0 {
		return errors.NewAniError("路径中不允许出现'..'")
	}

	if !utils.IsExist(path.Join(store.Config.DataPath, "plugin", p.Name)) {
		return errors.NewAniError("插件不存在")
	}
	return nil
}

type PluginConfigUploadRequest struct {
	PluginRequest
	Data string `json:"data" binding:"required"` //base64格式的文本数据
}

type PluginConfigDownloadRequest struct {
	PluginRequest
}
