package models

import "github.com/wetor/AnimeGo/configs"

type PluginResponse struct {
	Name string `json:"name"`
}

type PluginConfigResponse struct {
	PluginResponse
	Data string `json:"data"` //base64编码后的数据
}

type ConfigResponse struct {
	Config *configs.Config `json:"config,omitempty"`
	Data   string          `json:"data,omitempty"`
}

type BoltListResponse struct {
	Type   string   `json:"type"` // bucket, key
	Bucket string   `json:"bucket,omitempty"`
	Data   []string `json:"data"`
}

type BoltGetResponse struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	TTL    int64  `json:"ttl"`
	Value  any    `json:"value"`
}
