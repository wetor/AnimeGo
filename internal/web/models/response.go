package models

import "time"

type PluginResponse struct {
	Name string `json:"name"`
}

type PluginConfigResponse struct {
	PluginResponse
	Data string `json:"data"` //base64编码后的数据
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

type DirResponse struct {
	Path  string `json:"path"`
	Files []File `json:"files"`
}

type File struct {
	IsDir     bool      `json:"is_dir"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"modify_time"`
	Comment   string    `json:"comment"`
	CanRead   bool      `json:"read"`
	CanWrite  bool      `json:"write"`
	CanDelete bool      `json:"delete"`
}
