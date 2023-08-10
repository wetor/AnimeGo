package api

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/constant"
	webModels "github.com/wetor/AnimeGo/internal/web/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

func checkPluginPerm(p string) (read, write, delete bool) {
	read = true
	write = true
	delete = true
	if len(p) == 0 || p == "/" || p == "." {
		write = false
		delete = false
	}

	dir := path.Dir(p)
	isPluginRoot := false
	if len(dir) == 0 || dir == "/" || dir == "." {
		isPluginRoot = true
	}
	name := xpath.Base(p)
	// 顶层指定文件夹不可编辑和删除
	if _, ok := constant.PluginDirComment[name]; ok && isPluginRoot {
		write = false
		delete = false
	}
	// 内置插件不可编辑和删除
	if assets.IsBuiltinPlugin(name) {
		write = false
		delete = false
	}
	// README不可编辑和删除
	if name == "README.md" {
		write = false
		delete = false
	}
	return
}

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
	if len(path) == 0 || path == "/" || path == "." {
		isPluginRoot = true
	}
	pluginPath := xpath.Join(constant.PluginPath, path)
	if !utils.IsDirExist(pluginPath) {
		log.Warnf("文件夹不存在: " + pluginPath)
		c.JSON(webModels.Fail("文件夹不存在: " + path))
		return
	}
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
			Name:    f.Name(),
			IsDir:   f.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime().Unix(),
		}
		if isPluginRoot {
			file.Comment = constant.PluginDirComment[f.Name()]
		}
		file.CanRead, file.CanWrite, file.CanDelete = checkPluginPerm(xpath.Join(path, f.Name()))
		files = append(files, file)
	}
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir && !files[j].IsDir {
			return true
		} else if !files[i].IsDir && files[j].IsDir {
			return false
		} else {
			return files[i].Name < files[j].Name
		}
	})
	c.JSON(webModels.Succ("获取成功", webModels.DirResponse{
		Path:  path,
		Files: files,
	}))
}

// PluginDirPost godoc
//
//	@Summary		创建文件夹
//	@Description	创建插件文件夹中的文件夹
//	@Tags			plugin
//	@Accept			json
//	@Produce		json
//	@Param			path	body		webModels.PathRequest	true	"创建文件夹信息"
//	@Success		200		{object}	webModels.Response
//	@Failure		300		{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/plugin/manager/dir [post]
func (a *Api) PluginDirPost(c *gin.Context) {
	path := c.GetString("path")
	pluginPath := xpath.Join(constant.PluginPath, path)
	if utils.IsDirExist(pluginPath) {
		log.Warnf("文件夹已存在: " + pluginPath)
		c.JSON(webModels.Fail("文件夹已存在: " + path))
		return
	}
	err := utils.CreateMutiDir(pluginPath)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail("创建文件夹失败"))
		return
	}
	c.JSON(webModels.Succ("创建文件夹成功"))
}

// PluginDirDelete godoc
//
//	@Summary		删除文件夹
//	@Description	删除插件文件夹中指定文件文件夹
//	@Tags			plugin
//	@Accept			json
//	@Produce		json
//	@Param			path	query		string	true	"路径"
//	@Success		200		{object}	string
//	@Failure		300		{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/plugin/manager/dir [delete]
func (a *Api) PluginDirDelete(c *gin.Context) {
	path := c.GetString("path")
	pluginFile := xpath.Join(constant.PluginPath, path)
	if !utils.IsDirExist(pluginFile) {
		log.Warnf("文件夹不存在: " + pluginFile)
		c.JSON(webModels.Fail("文件夹不存在: " + path))
		return
	}
	_, w, d := checkPluginPerm(path)
	if !w || !d {
		c.JSON(webModels.Fail("禁止删除: " + path))
		return
	}
	err := os.Remove(pluginFile)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail("删除文件夹失败: " + path))
		return
	}
	c.JSON(webModels.Succ("删除文件夹成功"))
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
	if !utils.IsFileExist(pluginFile) {
		log.Warnf("文件不存在: " + pluginFile)
		c.JSON(webModels.Fail("文件不存在: " + path))
		return
	}
	r, _, _ := checkPluginPerm(path)
	if !r {
		c.JSON(webModels.Fail("禁止读取: " + path))
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

// PluginFilePost godoc
//
//	@Summary		创建或修改插件文件
//	@Description	创建或修改插件文件夹中指定文件
//	@Tags			plugin
//	@Accept			plain
//	@Produce		json
//	@Param			path	query		string	true	"路径"
//	@Param			file	body		string	true	"文件内容"
//	@Success		200		{object}	webModels.Response
//	@Failure		300		{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/plugin/manager/file [post]
func (a *Api) PluginFilePost(c *gin.Context) {
	path := c.GetString("path")
	pluginFile := xpath.Join(constant.PluginPath, path)
	create := !utils.IsFileExist(pluginFile)
	if !create {
		_, w, _ := checkPluginPerm(path)
		if !w {
			c.JSON(webModels.Fail("禁止编辑: " + path))
			return
		}
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail("读取文件参数失败"))
	}
	defer c.Request.Body.Close()

	err = os.WriteFile(pluginFile, body, constant.FilePerm)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail("写入插件文件失败"))
		return
	}
	if create {
		c.JSON(webModels.Succ("创建插件成功"))
	} else {
		c.JSON(webModels.Succ("修改插件成功"))
	}
}

// PluginFileDelete godoc
//
//	@Summary		删除插件
//	@Description	删除插件文件夹中指定插件文件
//	@Tags			plugin
//	@Accept			json
//	@Produce		json
//	@Param			path	query		string	true	"路径"
//	@Success		200		{object}	string
//	@Failure		300		{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/plugin/manager/file [delete]
func (a *Api) PluginFileDelete(c *gin.Context) {
	path := c.GetString("path")
	pluginFile := xpath.Join(constant.PluginPath, path)
	if !utils.IsFileExist(pluginFile) {
		log.Warnf("文件不存在: " + pluginFile)
		c.JSON(webModels.Fail("文件不存在: " + path))
		return
	}
	_, w, d := checkPluginPerm(path)
	if !w || !d {
		c.JSON(webModels.Fail("禁止删除: " + path))
		return
	}
	err := os.Remove(pluginFile)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail("删除文件失败: " + path))
		return
	}
	c.JSON(webModels.Succ("删除文件成功"))
}

// PluginRename godoc
//
//	@Summary		重命名文件或文件夹
//	@Description	重命名插件文件夹中的文件或文件夹
//	@Tags			plugin
//	@Accept			json
//	@Produce		json
//	@Param			path	body		webModels.NewPathRequest	true	"重命名信息"
//	@Success		200		{object}	webModels.Response
//	@Failure		300		{object}	webModels.Response
//	@Security		ApiKeyAuth
//	@Router			/api/plugin/manager/rename [put]
func (a *Api) PluginRename(c *gin.Context) {
	var request webModels.NewPathRequest
	if !a.checkRequest(c, &request) {
		return
	}
	path, err := utils.CheckPath(request.Path)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.ErrIpt("路径参数错误"))
		c.Abort()
		return
	}
	pluginPath := xpath.Join(constant.PluginPath, path)
	if !utils.IsExist(pluginPath) {
		log.Warnf("文件或文件夹不存在: " + pluginPath)
		c.JSON(webModels.Fail("文件或文件夹不存在: " + path))
		return
	}

	newPath, err := utils.CheckPath(request.NewPath)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.ErrIpt("目标路径参数错误"))
		c.Abort()
		return
	}
	newPluginPath := xpath.Join(constant.PluginPath, newPath)
	if utils.IsExist(newPluginPath) {
		log.Warnf("目标文件或文件夹已存在: " + newPluginPath)
		c.JSON(webModels.Fail("目标文件或文件夹已存在: " + newPath))
		return
	}
	_, w, _ := checkPluginPerm(path)
	if !w {
		c.JSON(webModels.Fail("禁止编辑: " + path))
		return
	}

	err = utils.Rename(pluginPath, newPluginPath)
	if err != nil {
		log.DebugErr(err)
		c.JSON(webModels.Fail("重命名失败"))
		return
	}
	c.JSON(webModels.Succ("重命名成功"))
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
	err = os.WriteFile(filename, data, constant.FilePerm)
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
