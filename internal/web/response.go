package web

type PluginResponse struct {
	Name string `json:"name"`
}

type PluginConfigResponse struct {
	PluginResponse
	Data string `json:"data"` //base64编码后的数据
}
