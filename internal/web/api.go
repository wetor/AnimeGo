package web

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wetor/AnimeGo/internal/animego/feed/rss"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"time": time.Now().Unix(),
		"pong": true,
	})
}

func Rss(c *gin.Context) {
	var request SelectEpRequest
	if !checkRequest(c, &request) {
		return
	}
	rss := rss.NewRss(request.Rss.Url, "")

	items, _ := rss.Parse()
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
	go store.Process.UpdateFeed(items)
	c.JSON(Succ(fmt.Sprintf("开始处理%d个下载项", len(items))))
}

func PluginConfigPost(c *gin.Context) {
	var request PluginConfigUploadRequest
	if !checkRequest(c, &request) {
		return
	}
	if err := request.CheckName(); err != nil {
		zap.S().Debug(err)
		c.JSON(Fail(err.Error()))
		return
	}

	data, err := base64.StdEncoding.DecodeString(request.Data)
	if err != nil {
		err = errors.NewAniErrorD(err)
		zap.S().Debug(err)
		c.JSON(Fail(err.Error()))
		return
	}
	filename := strings.TrimSuffix(request.Name, ".js") + ".json"
	file := path.Join(store.Config.DataPath, "plugin", filename)
	err = os.WriteFile(file, data, 0666)
	if err != nil {
		err = errors.NewAniErrorD(err)
		zap.S().Debug(err)
		c.JSON(Fail(err.Error()))
		return
	}
	c.JSON(Succ(fmt.Sprintf("写入插件配置文件成功，%s", filename)))
}

func PluginConfigGet(c *gin.Context) {
	var request PluginConfigDownloadRequest
	if !checkRequest(c, &request) {
		return
	}
	if err := request.CheckName(); err != nil {
		zap.S().Debug(err)
		c.JSON(Fail(err.Error()))
		return
	}
	filename := strings.TrimSuffix(request.Name, ".js") + ".json"
	file := path.Join(store.Config.DataPath, "plugin", filename)

	data, err := os.ReadFile(file)
	if err != nil {
		err = errors.NewAniErrorD(err)
		zap.S().Debug(err)
		c.JSON(Fail(err.Error()))
		return
	}
	str := base64.StdEncoding.EncodeToString(data)
	c.JSON(Succ("读取插件配置文件成功", PluginConfigResponse{
		PluginResponse: PluginResponse{
			Name: filename,
		},
		Data: str,
	}))
}
