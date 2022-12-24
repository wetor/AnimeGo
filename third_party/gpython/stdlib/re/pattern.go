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
	repl := string(arg[0].(py.String))
	str := string(arg[1].(py.String))
	count := -1
	if len(arg) >= 3 {
		count = int(arg[2].(py.Int))
	}
	return py.String(p.regx.ReplaceAllStringFunc(str, func(s string) string {
		if count == 0 {
			return s
		}
		count--
		return p.regx.ReplaceAllString(s, repl)
	})), nil
}

func (p *Pattern) split(arg py.Tuple) (py.Object, error) {
	str := string(arg[0].(py.String))
	maxSplit := -1
	if len(arg) >= 2 {
		maxSplit = int(arg[1].(py.Int))
	}
	return py.NewListFromStrings(p.regx.Split(str, maxSplit)), nil
}

func (p *Pattern) findAll(arg py.Object) (py.Object, error) {
	str := string(arg.(py.String))
	result := p.regx.FindAllStringSubmatch(str, -1)
	tuple := make([]py.Object, len(result))
	for i, list := range result {
		tuple[i] = _match(list)
	}
	return py.NewListFromItems(tuple), nil
}

func init() {
	PatternType.Dict["match"] = py.MustNewMethod("match", func(self py.Object, arg py.Object) (py.Object, error) {
		return self.(*Pattern).match(arg)
	}, 0, `pattern`)

	PatternType.Dict["search"] = py.MustNewMethod("search", func(self py.Object, arg py.Object) (py.Object, error) {
		return self.(*Pattern).search(arg)
	}, 0, `pattern`)

	PatternType.Dict["sub"] = py.MustNewMethod("sub", func(self py.Object, args py.Tuple) (py.Object, error) {
		return self.(*Pattern).sub(args)
	}, 0, `pattern`)

	PatternType.Dict["split"] = py.MustNewMethod("split", func(self py.Object, args py.Tuple) (py.Object, error) {
		return self.(*Pattern).split(args)
	}, 0, `pattern`)

	PatternType.Dict["findall"] = py.MustNewMethod("findall", func(self py.Object, arg py.Object) (py.Object, error) {
		return self.(*Pattern).findAll(arg)
	}, 0, `pattern`)

}
