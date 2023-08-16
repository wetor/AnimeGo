package models

import (
	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
	"path"
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

func (p PluginRequest) FindFile() (string, error) {
	file := p.Name
	if xpath.IsAbs(file) {
		file = xpath.Abs(xpath.P(file))
	} else {
		file = xpath.Abs(path.Join(constant.PluginPath, xpath.P(file)))
	}
	return utils.FindScript(file, ".py")
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
	Key string `json:"key" form:"key" default:"raw"`
}

type ConfigPutRequest struct {
	// Key 用路径方式更新指定yaml key内容
	//   [暂不支持] 如 setting/save_path, advanced/download/queue_max_num
	//   all 更新所有配置项，json格式
	//   raw 更新整个配置文件，base64编码
	Key string `json:"key" form:"key" default:"raw"`
	// Backup 备份原配置文件
	Backup    *bool           `json:"backup" form:"backup" default:"true"`
	Config    *configs.Config `json:"config"`
	ConfigRaw *string         `json:"config_raw"`
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

type AddItemsRequest struct {
	Source string `json:"source" binding:"required"`
	Data   []struct {
		Url  string         `json:"url" binding:"required"`
		Info map[string]any `json:"info"`
	} `json:"data" binding:"required"`
}
