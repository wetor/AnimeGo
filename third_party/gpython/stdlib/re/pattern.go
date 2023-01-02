package re

import (
	"fmt"
	"github.com/go-python/gpython/py"
	"regexp"
)

var PatternType = py.NewType("pattern", `regular expression pattern object`)

type Pattern struct {
	pattern string
	regx    *regexp.Regexp
}

// Type of this object
func (p *Pattern) Type() *py.Type {
	return PatternType
}

func (p *Pattern) M__str__() (py.Object, error) {
	return py.String(fmt.Sprintf("%v", p.pattern)), nil
}

func (p *Pattern) M__repr__() (py.Object, error) {
	return p.M__str__()
}

func _compile(pattern string) *Pattern {
	return &Pattern{
		pattern: pattern,
		regx:    regexp.MustCompile(pattern),
	}
}

func (p *Pattern) match(arg py.Object) (py.Object, error) {
	str := string(arg.(py.String))
	return _match(p.regx.FindStringSubmatch(str)), nil
}

func (p *Pattern) search(arg py.Object) (py.Object, error) {
	str := string(arg.(py.String))
	return _match(p.regx.FindStringSubmatch(str)), nil
}

func (p *Pattern) sub(arg py.Tuple) (py.Object, error) {
	var method func(s string) string
	var err error

	switch val := arg[0].(type) {
	case py.String:
		method = func(s string) string {
			return p.regx.ReplaceAllString(s, string(val))
		}
	case *py.Function:
		method = func(s string) string {
			ret, e := val.M__call__(py.Tuple{_match([]string{s})}, nil)
			err = e
			if ret == nil {
				return ""
			}
			return string(ret.(py.String))
		}
	}

	str := string(arg[1].(py.String))
	count := -1
	if len(arg) >= 3 {
		count = int(arg[2].(py.Int))
		if count == 0 {
			count = -1
		}
	}
	return py.String(p.regx.ReplaceAllStringFunc(str, func(s string) string {
		if count == 0 {
			return s
		}
		count--
		return method(s)
	})), err
}

func (p *Pattern) split(args py.Tuple) (py.Object, error) {
	str := string(args[0].(py.String))
	maxSplit := -1
	if len(args) >= 2 {
		maxSplit = int(args[1].(py.Int))
		if maxSplit == 0 {
			maxSplit = -1
		}
	}
	return py.NewListFromStrings(p.regx.Split(str, maxSplit)), nil
}

func (p *Pattern) findAll(args py.Tuple) (py.Object, error) {
	str := string(args[0].(py.String))
	pos := 0
	end := len(str)
	if len(args) >= 2 {
		pos = int(args[1].(py.Int))
	}
	if len(args) >= 3 {
		end = int(args[2].(py.Int))
	}
	result := p.regx.FindAllStringSubmatch(str[pos:end], -1)
	tuple := make([]py.Object, len(result))
	for i, list := range result {
		if len(list) > 1 {
			t := make(py.Tuple, len(list))
			for j, s := range list {
				t[j] = py.String(s)
			}
			tuple[i] = t
		} else if len(list) == 1 {
			tuple[i] = py.String(list[0])
		}
	}
	return py.NewListFromItems(tuple), nil
}

func init() {
	PatternType.Dict["match"] = py.MustNewMethod("match", func(self py.Object, arg py.Object) (py.Object, error) {
		return self.(*Pattern).match(arg)
	}, 0, `match(string)`)

	PatternType.Dict["search"] = py.MustNewMethod("search", func(self py.Object, arg py.Object) (py.Object, error) {
		return self.(*Pattern).search(arg)
	}, 0, `search(string)`)

	PatternType.Dict["sub"] = py.MustNewMethod("sub", func(self py.Object, args py.Tuple) (py.Object, error) {
		return self.(*Pattern).sub(args)
	}, 0, `sub(repl, string[, count=0])`)

	PatternType.Dict["split"] = py.MustNewMethod("split", func(self py.Object, args py.Tuple) (py.Object, error) {
		return self.(*Pattern).split(args)
	}, 0, `split(string[, maxsplit=0])`)

	PatternType.Dict["findall"] = py.MustNewMethod("findall", func(self py.Object, args py.Tuple) (py.Object, error) {
		return self.(*Pattern).findAll(args)
	}, 0, `findall(string[, pos[, endpos]])`)

}
