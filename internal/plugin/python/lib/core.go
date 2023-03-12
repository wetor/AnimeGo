package lib

import (
	"github.com/go-python/gpython/py"
	"gopkg.in/yaml.v3"

	"github.com/wetor/AnimeGo/internal/animego/anidata/mikan"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/feed/rss"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/try"
	"github.com/wetor/AnimeGo/pkg/utils"
)

func InitAnimeGo() {
	methods := []*py.Method{
		py.MustNewMethod("loads", loads, 0,
			`loads(s, type='json') -> dict
type support 'json' and 'yaml'`),
		py.MustNewMethod("dumps", dumps, 0,
			`dumps(obj, type='json') -> str
type support 'json' and 'yaml'`),
		py.MustNewMethod("parse_mikan", parseMikan, 0,
			`parse_mikan(url) -> dict
id: Mikan ID
sub_group_id: Mikan subgroup ID
pub_group_id: Mikan pubgroup ID
group_name: 字幕组名`),
		py.MustNewMethod("parse_mikan_rss", parseMikanRss, 0,
			`parse_mikan_rss(rss_data) -> dict`),
	}

	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "core",
			Doc:  "AnimeGo Core Module",
		},
		Methods: methods,
	})
}

func loads(self py.Object, args py.Tuple) (py.Object, error) {
	var content []byte
	encodng := "json"
	if data, ok := args[0].(py.Bytes); ok {
		content = data
	} else if data, ok := args[0].(py.String); ok {
		content = []byte(data)
	} else {
		return nil, py.ExceptionNewf(py.TypeError, "a str or bytes is required")
	}

	if len(args) > 1 {
		encodng = string(args[1].(py.String))
	}

	var err error
	result := map[string]any{}

	switch encodng {
	case "json":
		err = json.Unmarshal(content, &result)
	case "yaml":
		err = yaml.Unmarshal(content, &result)
	default:
		return nil, py.ExceptionNewf(py.TypeError, "only 'json' or 'yaml' is supported")
	}
	if err != nil {
		return nil, err
	}
	return plugin.Value2PyObject(result), nil
}

func dumps(self py.Object, args py.Tuple) (py.Object, error) {
	var obj py.StringDict
	encodng := "json"
	if data, ok := args[0].(py.StringDict); ok {
		obj = data
	} else {
		obj = args[0].Type().GetDict()
		delete(obj, "__module__")
	}
	if len(args) > 1 {
		encodng = string(args[1].(py.String))
	}

	object := plugin.PyObject2Value(obj)
	var result []byte
	var err error

	switch encodng {
	case "json":
		result, err = json.Marshal(object)
	case "yaml":
		result, err = yaml.Marshal(object)
	default:
		return nil, py.ExceptionNewf(py.TypeError, "only 'json' or 'yaml' is supported")
	}
	if err != nil {
		return nil, err
	}
	return plugin.Value2PyObject(string(result)), nil
}

func parseMikan(self py.Object, arg py.Object) (py.Object, error) {
	var info *mikan.MikanInfo
	var err error
	try.This(func() {
		info = anisource.Mikan().CacheParseMikanInfo(string(arg.(py.String)))
	}).Catch(func(e try.E) {
		err = e.(error)
	})
	if err != nil {
		return nil, err
	}
	obj := map[string]any(utils.StructToMap(info))

	return plugin.Value2PyObject(obj), nil
}

func parseMikanRss(self py.Object, arg py.Object) (py.Object, error) {
	var items []*models.FeedItem
	var err error
	try.This(func() {
		items = rss.NewRss(&rss.Options{Raw: string(arg.(py.String))}).Parse()
	}).Catch(func(e try.E) {
		err = e.(error)
	})
	if err != nil {
		return nil, err
	}
	return plugin.Value2PyObject(items), nil
}
