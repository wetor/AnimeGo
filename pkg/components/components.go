package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
	"net/http"
)

func CreateHandler(title string, body g.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Rendering a Node is as simple as calling Render and passing an io.Writer
		_ = Page(title, r.URL.Path, body).Render(w)
	}
}

func Page(title, path string, body g.Node) g.Node {
	// HTML5 boilerplate document
	return c.HTML5(c.HTML5Props{
		Title:    title,
		Language: "en",
		Head: []g.Node{
			Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/bootstrap@4.6.2/dist/css/bootstrap.min.css")),
			Link(Rel("stylesheet"), Href("https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.0/font/bootstrap-icons.css")),
			Script(Src("https://cdn.jsdelivr.net/npm/jquery@3.5.1/dist/jquery.min.js")),
			Script(Src("https://cdn.jsdelivr.net/npm/popper.js@1.16.1/dist/umd/popper.min.js")),
			Script(Src("https://cdn.jsdelivr.net/npm/bootstrap@4.6.2/dist/js/bootstrap.bundle.min.js")),
			Link(Rel("stylesheet"), Href("/static/global.css")),
		},
		Body: []g.Node{
			Container(
				Prose(body),
			),
			Script(Src("/static/global.js")),
		},
	})
}

type PageLink struct {
	Path string
	Name string
}

func Container(children ...g.Node) g.Node {
	return Div(Class("max-w-7xl mx-auto px-2 sm:px-6 lg:px-8"), g.Group(children))
}

func Prose(children ...g.Node) g.Node {
	return Div(Class("prose"), g.Group(children))
}
