package re

import (
	"regexp"

	"github.com/go-python/gpython/py"
)

const match_doc = `The result of re.match() and re.search().
Match objects always have a boolean value of True.`

var MatchType = py.NewType("match", match_doc)

func init() {
	MatchType.Dict["group"] = py.MustNewMethod("group",
		func(self py.Object, args py.Tuple) (py.Object, error) {
			return self.(*Match).group(args)
		}, 0, match_group_doc)

	MatchType.Dict["start"] = py.MustNewMethod("start",
		func(self py.Object, args py.Tuple) (py.Object, error) {
			return self.(*Match).start(args)
		}, 0, match_start_doc)

	MatchType.Dict["end"] = py.MustNewMethod("end",
		func(self py.Object, args py.Tuple) (py.Object, error) {
			return self.(*Match).end(args)
		}, 0, match_end_doc)

	MatchType.Dict["span"] = py.MustNewMethod("span",
		func(self py.Object, args py.Tuple) (py.Object, error) {
			return self.(*Match).span(args)
		}, 0, match_span_doc)

	MatchType.Dict["groups"] = py.MustNewMethod("groups",
		func(self py.Object, args py.Tuple) (py.Object, error) {
			return self.(*Match).groups(args)
		}, 0, match_groups_doc)

	MatchType.Dict["groupdict"] = py.MustNewMethod("groupdict",
		func(self py.Object, args py.Tuple) (py.Object, error) {
			return self.(*Match).groupdict(args)
		}, 0, match_groupdict_doc)

	MatchType.Dict["expand"] = py.MustNewMethod("expand",
		func(self py.Object, args py.Object) (py.Object, error) {
			return self.(*Match).expand(args)
		}, 0, match_expand_doc)
}

type Match struct {
	regx    *regexp.Regexp
	string_ py.Object
	groups_ py.Tuple
	index_  []int
	pos     int
	endpos  int
}

func MatchNew(regx *regexp.Regexp, string_, pos, endpos py.Object) py.Object {
	slice, start, end := getSlice(string_, pos, endpos)
	str, isBytes := toString(slice)
	groups := regx.FindStringSubmatch(str)
	index := regx.FindStringSubmatchIndex(str)
	if len(groups) == 0 {
		return py.None
	}
	groupStr := make(py.Tuple, len(groups))
	if isBytes {
		for i := 0; i < len(groups); i++ {
			groupStr[i] = py.Bytes(groups[i])
		}
	} else {
		for i := 0; i < len(groups); i++ {
			groupStr[i] = py.String(groups[i])
		}
	}
	return &Match{
		regx:    regx,
		groups_: groupStr,
		index_:  index,
		pos:     start,
		endpos:  end,
	}
}

// Type of this object
func (m *Match) Type() *py.Type {
	return MatchType
}

func (m *Match) M__str__() (py.Object, error) {
	return m.M__repr__()
}

func (m *Match) M__repr__() (py.Object, error) {
	if len(m.groups_) == 1 {
		return m.groups_[0].(py.I__repr__).M__repr__()
	}
	return py.NewListFromItems(m.groups_).M__repr__()
}

func (m *Match) M__bool__() (py.Object, error) {
	return py.Bool(len(m.groups_) > 0), nil
}

const match_group_doc = `group([group1, ...]) -> str or tuple.
    Return subgroup(s) of the match by indices or names.
    For 0 returns the entire match.`

func (m *Match) group(args py.Tuple) (res py.Object, err error) {
	size := len(args)
	switch len(args) {
	case 0:
		res = m.getslice(py.Int(0), py.None)
	case 1:
		res = m.getslice(args[0], py.None)
	default:
		tuple := make(py.Tuple, size)
		for i := 0; i < size; i++ {
			index := m.getslice(args[i], py.None)
			tuple[i] = index
		}
		res = tuple
	}
	return res, nil
}

func (m *Match) getslice(index py.Object, def py.Object) (res py.Object) {
	switch t := index.(type) {
	case py.Int:
		res = m.groups_[t]
	case py.String:
		i := m.regx.SubexpIndex(string(t))
		res = m.groups_[i]
	}

	size := 0
	if l, ok := res.(py.I__len__); ok {
		res2, _ := l.M__len__()
		size = int(res2.(py.Int))
	} else if bytes, ok := res.(py.Bytes); ok {
		size = len(bytes)
	}
	if size == 0 {
		return def
	}
	return res
}

const match_start_doc = `start([group=0]) -> int.
    Return index of the start of the substring matched by group.`

func (m *Match) start(args py.Tuple) (py.Object, error) {
	var index_ py.Object = py.Int(0)
	err := py.UnpackTuple(args, nil, "start", 0, 1, &index_)
	if err != nil {
		return nil, err
	}
	index := int(index_.(py.Int))
	return py.Int(m.pos + m.index_[index*2]), nil
}

const match_end_doc = `end([group=0]) -> int.
    Return index of the end of the substring matched by group.`

func (m *Match) end(args py.Tuple) (py.Object, error) {
	var index_ py.Object = py.Int(0)
	err := py.UnpackTuple(args, nil, "end", 0, 1, &index_)
	if err != nil {
		return nil, err
	}
	index := int(index_.(py.Int))
	return py.Int(m.pos + m.index_[index*2+1]), nil
}

const match_span_doc = `span([group]) -> tuple.
    For MatchObject m, return the 2-tuple (m.start(group), m.end(group)).`

func (m *Match) span(args py.Tuple) (py.Object, error) {
	var index_ py.Object = py.Int(0)
	err := py.UnpackTuple(args, nil, "span", 0, 1, &index_)
	if err != nil {
		return nil, err
	}
	index := int(index_.(py.Int))
	return py.Tuple{py.Int(m.pos + m.index_[index*2]), py.Int(m.pos + m.index_[index*2+1])}, nil
}

const match_groups_doc = `groups([default=None]) -> tuple.
    Return a tuple containing all the subgroups of the match, from 1.
    The default argument is used for groups
    that did not participate in the match`

func (m *Match) groups(args py.Tuple) (py.Object, error) {
	var def py.Object = py.None
	if len(args) == 1 {
		def = args[0]
	}
	tuple := make(py.Tuple, len(m.groups_)-1)
	for i := 0; i < len(m.groups_)-1; i++ {
		item := m.groups_[i+1]
		size := 0
		if l, ok := item.(py.I__len__); ok {
			res, err := l.M__len__()
			if err != nil {
				return nil, err
			}
			size = int(res.(py.Int))
		} else if bytes, ok := item.(py.Bytes); ok {
			size = len(bytes)
		}
		if size == 0 {
			tuple[i] = def
		} else {
			tuple[i] = m.groups_[i+1]
		}
	}
	return tuple, nil
}

const match_groupdict_doc = `groupdict([default=None]) -> dict.
    Return a dictionary containing all the named subgroups of the match,
    keyed by the subgroup name. The default argument is used for groups
    that did not participate in the match`

func (m *Match) groupdict(args py.Tuple) (py.Object, error) {
	return nil, nil
}

const match_expand_doc = `expand(template) -> str.
    Return the string obtained by doing backslash substitution
    on the string template, as done by the sub() method.`

func (m *Match) expand(args py.Object) (py.Object, error) {
	return nil, nil
}
