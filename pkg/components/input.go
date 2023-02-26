package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
	"strconv"
)

func BaseLabel(name *NameInfo) []g.Node {
	return []g.Node{
		h.Label(
			g.Text(name.DisplayName),
		),
		h.A(
			h.Class("bi-info-circle"),
			g.Attr("role", "button"),
			h.DataAttr("container", "body"),
			h.DataAttr("toggle", "popover"),
			h.DataAttr("placement", "right"),
			h.DataAttr("trigger", "hover"),
			h.TabIndex("0"),
			h.DataAttr("html", "true"),
			h.DataAttr("content", name.Comment),
		),
	}
}

func BaseInput(name *NameInfo, nodes ...g.Node) g.Node {
	n := make([]g.Node, 0, len(nodes)+3)
	n = append(n, h.Class("form-group"))
	n = append(n, h.DataAttr("name", name.Last()))
	n = append(n, BaseLabel(name)...)
	n = append(n, nodes...)
	return h.Div(n...)
}

func BaseComment(comment string) g.Node {
	return g.If(len(comment) > 0, h.Div(
		h.Class("card-body"),
		h.Label(
			h.Class("text-secondary"),
			g.Raw(comment),
		),
	))
}

type InputOptions struct {
	NameInfo
	Value   any
	Default any
}

func NumberInput(opts InputOptions) g.Node {
	var defaultValue int
	if opts.Default != nil {
		defaultValue = opts.Default.(int)
	}
	return BaseInput(&opts.NameInfo, h.Input(
		h.Type("text"),
		h.ID(opts.Name),
		g.If(!opts.Hidden, h.Name(opts.Name)),
		h.Class("form-control a-input"),
		g.Attr("placeholder", strconv.Itoa(defaultValue)),
		h.Value(strconv.Itoa(opts.Value.(int))),
	))
}

func StringInput(opts InputOptions) g.Node {
	var defaultValue string
	if opts.Default != nil {
		defaultValue = opts.Default.(string)
	}
	return BaseInput(&opts.NameInfo, h.Input(
		h.Type("text"),
		h.ID(opts.Name),
		g.If(!opts.Hidden, h.Name(opts.Name)),
		h.Class("form-control a-input"),
		g.Attr("placeholder", defaultValue),
		h.Value(opts.Value.(string)),
	))
}

func BoolInput(opts InputOptions) g.Node {
	v := "false"
	if opts.Value.(bool) {
		v = "true"
	}
	return SelectInput(opts.NameInfo, v, map[string]string{
		"true":  "启用",
		"false": "禁用",
	})
}

func SelectInput(name NameInfo, value string, opts map[string]string) g.Node {
	options := make([]g.Node, 0, len(opts))
	for k, v := range opts {
		options = append(options, h.Option(
			g.If(k == value, h.Selected()),
			h.Value(k),
			g.Text(v),
		))

	}
	options = append(options, h.Class("custom-select a-input"))
	options = append(options, g.If(!name.Hidden, h.Name(name.Name)))
	options = append(options, h.ID(name.Name))
	return BaseInput(&name, h.Select(options...))
}

func StructCard(name NameInfo, nodes ...g.Node) g.Node {
	nodes = append(nodes, h.Class("card-body"))
	id := name.ID("card")
	return h.Div(
		h.Class("card mb-3"),
		h.Div(
			h.Class("card-header"),
			h.DataAttr("toggle", "collapse"),
			h.DataAttr("target", "#"+id),
			g.Attr("role", "button"),
			g.Text(name.DisplayName),
		),
		h.Div(
			h.Class("collapse"),
			h.ID(id),
			g.If(name.Comment != name.DisplayName, BaseComment(name.Comment)),
			h.Div(nodes...),
		),
	)
}

func ArrayItem(name NameInfo, index int, nodes ...g.Node) g.Node {
	var template bool
	if index < 0 {
		template = true
	}
	nodes = append(nodes, h.Class("card-body"))
	id := name.ID("item")

	return h.Li(
		g.If(template, g.Attr("hidden", "")),
		g.If(template, h.Class("list-group-item a-template")),
		g.If(!template, h.Class("list-group-item")),
		h.Div(
			h.Class("card mb-3"),
			h.Div(
				h.Class("card-header a-item-header d-flex"),
				h.Div(
					h.Class("flex-grow-1"),
					h.DataAttr("toggle", "collapse"),
					h.DataAttr("target", "#"+id),
					g.Attr("role", "button"),
					h.Span(
						g.Text(name.DisplayName),
					),
				),
				h.Button(
					h.Type("button"),
					h.Class("btn btn-outline-danger btn-sm float-right a-list-deleter"),
					g.Text("删除"),
				),
				h.Button(
					h.Type("button"),
					h.Class("btn btn-outline-info btn-sm float-right a-list-up"),
					g.Text("上移"),
				),
				h.Button(
					h.Type("button"),
					h.Class("btn btn-outline-info btn-sm float-right a-list-down"),
					g.Text("下移"),
				),
			),
			h.Div(
				h.Class("collapse a-item-body"),
				h.ID(id),
				g.If(name.Comment != name.DisplayName, BaseComment(name.Comment)),
				h.Div(nodes...),
			),
		),
	)
}

func ArrayList(name NameInfo, nodes ...g.Node) g.Node {
	lis := make([]g.Node, 0, len(nodes)+3)
	lis = append(lis, h.Class("list-group list-group-flush"))
	for _, node := range nodes {
		lis = append(lis, node)
	}
	lis = append(lis, h.Button(
		h.Type("button"),
		h.Class("btn btn-outline-success btn-sm a-list-adder"),
		g.Text("新增"),
	))
	id := name.ID("list")
	return h.Div(
		h.Class("card mb-3"),
		h.Div(
			h.Class("card-header"),
			h.DataAttr("toggle", "collapse"),
			h.DataAttr("target", "#"+id),
			g.Attr("role", "button"),
			g.Text(name.DisplayName),
		),
		h.Div(
			h.Class("collapse a-array"),
			h.ID(id),
			h.DataAttr("name", name.Name),
			g.If(name.Comment != name.DisplayName, BaseComment(name.Comment)),
			h.Ul(lis...),
		),
	)
}
