package main

import (
	"github.com/gin-gonic/gin"
	h "github.com/maragudk/gomponents/html"
	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/pkg/components"
)

func main() {
	conf := configs.DefaultConfig()

	c := components.NewComponents(configs.DefaultCommentMap())
	name := components.NameInfo{
		Name:        "Config",
		DisplayName: "Config",
		Comment:     "这是注释\n换行策略",
	}

	from := h.FormEl(
		h.Action("/test"),
		h.Method("post"),
		c.Struct2Node(name, conf),
		h.Input(
			h.Type("submit"),
			h.Class("btn btn-primary"),
			h.Value("提交"),
		),
	)

	router := gin.Default()
	components.RegisterStatic(router)
	router.GET("/", func(c *gin.Context) {
		components.CreateHandler("", from).ServeHTTP(c.Writer, c.Request)
	})
	// 启动服务器
	router.Run(":8080")
}
