package ui

import (
	"github.com/gin-gonic/gin"
	h "github.com/maragudk/gomponents/html"
	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/pkg/components"
)

func PageConfig(c *gin.Context) {
	conf := configs.DefaultConfig()

	cmp := components.NewComponents(configs.DefaultCommentMap())
	name := components.NameInfo{
		Name:        "Config",
		DisplayName: "Config",
		Comment:     "这是注释\n换行策略",
	}

	from := h.FormEl(
		h.Action("/test"),
		h.Method("post"),
		cmp.Struct2Node(name, conf),
		h.Input(
			h.Type("submit"),
			h.Class("btn btn-primary"),
			h.Value("提交"),
		),
	)
	Render(c, "test", from)
}
