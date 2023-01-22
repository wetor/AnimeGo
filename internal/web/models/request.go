package models

import (
	"path"

	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
)

type SelectEpRequest struct {
	Source string `json:"source" binding:"required"`
	Rss    struct {
		Url string `json:"url" binding:"required"`
	} `json:"rss" binding:"required"`
	IsSelectEp bool     `json:"is_select_ep" default:"false"`
	EpLinks    []string `json:"ep_links"`
}

type PluginRequest struct {
	Name string `json:"name" form:"name" binding:"required"` //插件文件名
}

func (p PluginRequest) FindFile(dataPath string) (string, error) {
	file := path.Join(dataPath, "plugin", p.Name)
	ext := path.Ext(file)
	if len(ext) == 0 {
		ext = models.JSExt
	}
	script := utils.FindScript(file, ext)
	return script, nil
}

type PluginConfigUploadRequest struct {
	PluginRequest
	Data string `json:"data" binding:"required"` //base64格式的文本数据
}

type PluginConfigDownloadRequest struct {
	PluginRequest
}

type ConfigGetRequest struct {
	// Key 使用路径方式获取指定yaml key内容
	//   [暂不支持] 如 setting.save_path, advanced.download.queue_max_num
	//   all 获取所有配置项，json格式
	//   default 获取默认值配置项，json格式
	//   comment 获取所有配置项的注释文本，json格式
	//   raw 获取所有配置项，yaml文件内容，base64编码
	Key string `json:"key" form:"key" default:"all"`
}

type ConfigPutRequest struct {
	// Key 用路径方式更新指定yaml key内容
	//   [暂不支持] 如 setting/save_path, advanced/download/queue_max_num
	//   all 更新所有配置项，json格式
	Key string `json:"key" form:"key" default:"all"`
	// Backup 备份原配置文件
	Backup bool            `json:"backup" form:"backup" default:"true"`
	Config *configs.Config `json:"config"`
}

type BoltRequest struct {
	DB string `json:"db" form:"db" default:"bolt"` // bolt, bolt_sub
}

type BoltListRequest struct {
	BoltRequest
	Type   string `json:"type" form:"type" binding:"required"` // bucket, key
	Bucket string `json:"bucket" form:"bucket"`                // 当type=key时，需要此参数
}

type BoltGetRequest struct {
	BoltRequest
	Bucket string `json:"bucket" form:"bucket" binding:"required"`
	Key    string `json:"key" form:"key" binding:"required"`
}

type BoltDeleteRequest struct {
	BoltRequest
	Bucket string `json:"bucket" form:"bucket" binding:"required"`
	Key    string `json:"key" form:"key" binding:"required"`
}
