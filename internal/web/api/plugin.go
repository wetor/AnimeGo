package api

import (
	"encoding/base64"
	"github.com/wetor/AnimeGo/assets"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/wetor/AnimeGo/internal/constant"
	webModels "github.com/wetor/AnimeGo/internal/web/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

// PluginDirGet godoc
//
//	@Summary		获取插件文件列表
//	@Description	获取插件文件夹中指定文件夹的文件列表
//	@Tags			plugin
//	@Accept			json
//	@Produce		json
//	@Param			path	query		string	true	"路径"
//	@Success		200		{object}	webModels.Response{data=webModels.DirResponse}
//	@Failure		300		{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/plugin/manager/dir [get]
func (a *Api) PluginDirGet(c *gin.Context) {
	path := c.GetString("path")
	isPluginRoot := false
	if len(path) == 0 || path == "/" {
		isPluginRoot = true
	}
	pluginPath := xpath.Join(constant.PluginPath, path)

	dirs, err := os.ReadDir(pluginPath)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail("读取文件夹失败"))
		return
	}
	files := make([]webModels.File, 0, len(dirs))
	for _, f := range dirs {
		info, err := f.Info()
		if err != nil {
			log.DebugErr(err)
			continue
		}
		file := webModels.File{
			Name:      f.Name(),
			IsDir:     f.IsDir(),
			Size:      info.Size(),
			ModTime:   info.ModTime(),
			CanRead:   true,
			CanWrite:  true,
			CanDelete: true,
		}
		if isPluginRoot {
			ok := false
			file.Comment, ok = constant.PluginDirComment[f.Name()]
			if ok {
				file.CanWrite = false
				file.CanDelete = false
			}
		}
		if assets.IsBuiltinPlugin(f.Name()) {
			file.CanWrite = false
			file.CanDelete = false
		}
		if f.Name() == "README.md" {
			file.CanWrite = false
			file.CanDelete = false
		}
		files = append(files, file)
	}
	c.JSON(webModels.Succ("获取成功", webModels.DirResponse{
		Path:  path,
		Files: files,
	}))
}

// PluginFileGet godoc
//
//	@Summary		获取插件文件内容
//	@Description	获取插件文件夹中指定文件的内容
//	@Tags			plugin
//	@Accept			json
//	@Produce		plain
//	@Param			path	query		string	true	"路径"
//	@Success		200		{object}	string
//	@Failure		300		{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/plugin/manager/file [get]
func (a *Api) PluginFileGet(c *gin.Context) {
	path := c.GetString("path")
	pluginFile := xpath.Join(constant.PluginPath, path)
	if !utils.IsExist(pluginFile) {
		log.Warnf("文件不存在: " + pluginFile)
		c.JSON(webModels.Fail("文件不存在: " + path))
		return
	}
	data, err := os.ReadFile(pluginFile)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail("打开文件失败: " + path))
		return
	}
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.String(http.StatusOK, string(data))
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
		log.DebugErr(err)
		c.JSON(webModels.Fail("配置解析错误"))
		return
	}

	filename := strings.TrimSuffix(file, xpath.Ext(file)) + ".json"
	err = os.WriteFile(filename, data, 0666)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail("写入文件失败"))
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
