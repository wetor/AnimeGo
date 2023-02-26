package ui

import (
	"github.com/gin-gonic/gin"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

func Render(c *gin.Context, title string, body g.Node) {
	_ = Page(title, c.Request.URL.Path, body).Render(c.Writer)
}

func Page(title, path string, body g.Node) g.Node {
	// HTML5 boilerplate document
	return c.HTML5(c.HTML5Props{
		Title:    title,
		Language: "en",
		Head: []g.Node{
			h.Link(h.Rel("stylesheet"), h.Href("https://cdn.jsdelivr.net/npm/bootstrap@4.6.2/dist/css/bootstrap.min.css")),
			h.Link(h.Rel("stylesheet"), h.Href("https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.0/font/bootstrap-icons.css")),
			h.Script(h.Src("https://cdn.jsdelivr.net/npm/jquery@3.5.1/dist/jquery.min.js")),
			h.Script(h.Src("https://cdn.jsdelivr.net/npm/popper.js@1.16.1/dist/umd/popper.min.js")),
			h.Script(h.Src("https://cdn.jsdelivr.net/npm/bootstrap@4.6.2/dist/js/bootstrap.bundle.min.js")),
			h.Link(h.Rel("stylesheet"), h.Href("/static/global.css")),
		},
		Body: []g.Node{
			h.Div(
				h.Class("card"),
				h.Div(
					h.Class("card-header"),
					g.Text(path),
				),
				h.Div(
					h.Class("card-body"),
					body,
				),
			),
			h.Script(h.Src("/static/global.js")),
		},
	})
}
