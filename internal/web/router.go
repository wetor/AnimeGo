package web

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(r gin.IRouter) {
	r.GET("/ping", API.Ping)
	r.GET("/sha256", API.SHA256)

	apiRoot := r.Group("/api")
	apiRoot.Use(KeyAuth())
	apiRoot.POST("/rss", API.Rss)

	apiRoot.POST("/plugin/config", API.PluginConfigPost)
	apiRoot.GET("/plugin/config", API.PluginConfigGet)
	pluginManager := apiRoot.Group("/plugin/manager", CheckPath())
	pluginManager.GET("/dir", API.PluginDirGet)
	pluginManager.GET("/file", API.PluginFileGet)

	apiRoot.GET("/config", API.ConfigGet)
	apiRoot.PUT("/config", API.ConfigPut)
	apiRoot.GET("/config/file", API.ConfigFileGet)
	apiRoot.PUT("/config/file", API.ConfigFilePut)

	apiRoot.GET("/bolt", API.BoltList)
	apiRoot.GET("/bolt/value", API.Bolt)
	apiRoot.DELETE("/bolt/value", API.BoltDelete)

	apiRoot.POST("/download/manager", API.AddItems)

	wsRoot := r.Group("/websocket")
	wsRoot.Use(KeyAuth())
	wsRoot.GET("/log", WS.Log)
}
