package components_test

import (
	"github.com/gin-gonic/gin"
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/pkg/components"
	"net/http"
	"testing"
)

func TestStruct(t *testing.T) {
	type obj3 struct {
		Str4  string `attr:"字符串4" comment:"用来筛选符合条件的项目进行解析下载"`
		Str5  string `comment_key:"test_key"`
		Bool1 bool
	}
	type obj2 struct {
		Str3        string
		Obj3        []*obj3 `attr:"[]Object3数组" comment:"用来筛选符合条件的项目进行解析下载"`
		ArrayString []string
	}

	type obj struct {
		Str1 string
		Int1 int
		Str2 string
		Obj2 *obj2
	}

	val := obj{
		Str1: "这是Str1",
		Str2: "str2_value",
		Int1: 10086,
		Obj2: &obj2{
			Str3: "sssss",
			Obj3: []*obj3{
				{
					Str4:  "这是Obj3_array_1",
					Str5:  "5555",
					Bool1: true,
				},
				{
					Str4: "这是Obj3_array_2",
					Str5: "Str55566",
				},
				{
					Str4: "这是Obj3_array_3",
				},
			},
			ArrayString: []string{
				"hello",
				"world",
			},
		},
	}
	c := components.NewComponents(map[string]string{
		"test_key": `这是comment key`,
	})
	name := components.NameInfo{
		Name:        "Config",
		DisplayName: "Config",
		Comment:     "这是注释\n换行策略",
	}
	node := c.Struct2Node(name, val)
	nodes := make([]g.Node, 0)
	nodes = append(nodes, h.Action("/test"))
	nodes = append(nodes, h.Method("post"))
	nodes = append(nodes, node...)
	nodes = append(nodes, h.Input(
		h.Type("submit"),
		h.Class("btn btn-primary"),
		h.Value("提交"),
	))

	from := h.FormEl(nodes...)

	http.Handle("/", components.CreateHandler("", from))

	_ = http.ListenAndServe("localhost:8080", nil)
}

func TestSturctConfig(t *testing.T) {
	conf := configs.DefaultConfig()

	c := components.NewComponents(configs.DefaultCommentMap())
	name := components.NameInfo{
		Name:        "Config",
		DisplayName: "Config",
		Comment:     "这是注释\n换行策略",
	}
	node := c.Struct2Node(name, conf)
	nodes := make([]g.Node, 0)
	nodes = append(nodes, h.ID("myForm"))
	nodes = append(nodes, node...)
	nodes = append(nodes, h.Input(
		h.Type("submit"),
		h.Class("btn btn-primary"),
		h.Value("提交"),
	))
	from := h.FormEl(nodes...)
	http.Handle("/", components.CreateHandler("", from))

	_ = http.ListenAndServe("localhost:8080", nil)
}

func TestGinSturctConfig(t *testing.T) {
	conf := configs.DefaultConfig()

	c := components.NewComponents(configs.DefaultCommentMap())
	name := components.NameInfo{
		Name:        "Config",
		DisplayName: "Config",
		Comment:     "这是注释\n换行策略",
	}
	node := c.Struct2Node(name, conf)
	nodes := make([]g.Node, 0)
	nodes = append(nodes, h.ID("myForm"))
	nodes = append(nodes, node...)
	nodes = append(nodes, h.Input(
		h.Type("submit"),
		h.Class("btn btn-primary"),
		h.Value("提交"),
	))
	from := h.FormEl(nodes...)

	router := gin.Default()
	components.RegisterStatic(router)
	router.GET("/", func(c *gin.Context) {
		components.CreateHandler("", from).ServeHTTP(c.Writer, c.Request)
	})
	// 启动服务器
	router.Run(":8080")
}
