package re

import (
	"fmt"
	"github.com/go-python/gpython/py"
)

var MatchType = py.NewType("match", `regular expression match object`)

type Match struct {
	groupStr []string
}

func _match(group []string) *Match {
	return &Match{
		groupStr: group,
	}
}

// Type of this object
func (m *Match) Type() *py.Type {
	return MatchType
}

func (m *Match) M__str__() (py.Object, error) {
	return py.String(fmt.Sprintf("%v", m.groupStr)), nil
}

func (m *Match) M__repr__() (py.Object, error) {
	return m.M__str__()
}

func (m *Match) M__bool__() (py.Object, error) {
	return py.Bool(len(m.groupStr) > 0), nil
}

func (m *Match) group(arg py.Object) (py.Object, error) {
	index := int(arg.(py.Int))
	return py.String(m.groupStr[index]), nil
}

func (m *Match) groups() (py.Object, error) {
	tuple := make([]py.Object, len(m.groupStr)-1)
	for i := 0; i < len(m.groupStr)-1; i++ {
		tuple[i] = py.String(m.groupStr[i+1])
	}
	return py.NewListFromItems(tuple), nil
}

func init() {
	MatchType.Dict["group"] = py.MustNewMethod("group", func(self py.Object, arg py.Object) (py.Object, error) {
		return self.(*Match).group(arg)
	}, 0, `match`)
	MatchType.Dict["groups"] = py.MustNewMethod("group", func(self py.Object) (py.Object, error) {
		return self.(*Match).groups()
	}, 0, `match`)
}
