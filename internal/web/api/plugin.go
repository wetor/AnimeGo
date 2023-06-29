package api

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/wetor/AnimeGo/internal/animego/feed/rss"
	"github.com/wetor/AnimeGo/internal/models"
	webModels "github.com/wetor/AnimeGo/internal/web/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/xerrors"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

// Rss godoc
//
//	@Summary		发送下载项
//	@Description	将待下载项组合成rss发送给AnimeGo
//	@Tags			plugin
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
	err = a.filterManager.Update(a.ctx, items, nil, false, true)
	if err != nil {
		c.JSON(webModels.Fail(err.Error()))
		return
	}
	c.JSON(webModels.Succ(fmt.Sprintf("开始处理%d个下载项", len(items))))
}

// PluginConfigPost godoc
//
//	@Summary		发送插件配置
//	@Description	将当前插件的配置发送给AnimeGo并保存
//	@Description	插件名为不包含 'plugin' 的路径
//	@Description	插件名可以忽略'.py'后缀；插件名也可以使用上层文件夹名，会自动加载文件夹内部的 'main.py'
//	@Description	如设置为 'plugin/test'，会依次尝试加载 'plugin/test/main.py', 'plugin/test.py'
//	@Tags			plugin
//	@Accept			json
//	@Produce		json
//	@Param			plugin	body		webModels.PluginConfigUploadRequest	true	"插件信息，data为base64编码后的json文本"
//	@Success		200		{object}	webModels.Response
//	@Failure		300		{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/plugin/config [post]
func (a *Api) PluginConfigPost(c *gin.Context) {
	var request webModels.PluginConfigUploadRequest
	if !a.checkRequest(c, &request) {
		return
	}
	file, err := request.FindFile()
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail(err.Error()))
		return
	}

	data, err := base64.StdEncoding.DecodeString(request.Data)
	if err != nil {
		err = xerrors.NewAniErrorD(err)
		log.DebugErr(err)
		c.JSON(webModels.Fail(err.Error()))
		return
	}

	filename := strings.TrimSuffix(file, xpath.Ext(file)) + ".json"
	err = os.WriteFile(filename, data, 0666)
	if err != nil {
		err = xerrors.NewAniErrorD(err)
		log.DebugErr(err)
		c.JSON(webModels.Fail(err.Error()))
		return
	}
	c.JSON(webModels.Succ("写入插件配置文件成功", webModels.PluginResponse{
		Name: request.Name,
	}))
}

// PluginConfigGet godoc
//
//	@Summary		获取插件配置
//	@Description	从AnimeGo中获取当前插件的配置
//	@Description	插件名为不包含 'plugin' 的路径
//	@Description	插件名可以忽略'.js'后缀；插件名也可以使用上层文件夹名，会自动寻找文件夹内部的 'main.js' 或 'plugin.js'
//	@Description	如传入 'test'，会依次尝试寻找 'plugin/test/main.js', 'plugin/test/plugin.js', 'plugin/test.js'
//	@Tags			plugin
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	true	"插件信息"
//	@Success		200		{object}	webModels.Response
//	@Failure		300		{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/plugin/config [get]
func (a *Api) PluginConfigGet(c *gin.Context) {
	var request webModels.PluginConfigDownloadRequest
	if !a.checkRequest(c, &request) {
		return
	}
	file, err := request.FindFile()
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail(err.Error()))
		return
	}
	filename := strings.TrimSuffix(file, ".js") + ".json"

	data, err := os.ReadFile(filename)
	if err != nil {
		err = xerrors.NewAniErrorD(err)
		log.DebugErr(err)
		c.JSON(webModels.Fail("打开文件 " + filename + " 失败"))
		return
	}
	str := base64.StdEncoding.EncodeToString(data)
	c.JSON(webModels.Succ("读取插件配置文件成功", webModels.PluginConfigResponse{
		PluginResponse: webModels.PluginResponse{
			Name: request.Name,
		},
		Data: str,
	}))
}
