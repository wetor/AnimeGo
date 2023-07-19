package re

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/go-python/gpython/parser"
	"github.com/go-python/gpython/py"
)

const pattern_doc = `Compiled regular expression objects`

var PatternType = py.NewType("pattern", pattern_doc)

func init() {
	PatternType.Dict["match"] = py.MustNewMethod("match",
		func(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
			return self.(*Pattern).match(args, kwargs)
		}, 0, pattern_match_doc)

	PatternType.Dict["fullmatch"] = py.MustNewMethod("fullmatch",
		func(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
			return self.(*Pattern).fullmatch(args, kwargs)
		}, 0, pattern_fullmatch_doc)

	PatternType.Dict["search"] = py.MustNewMethod("search",
		func(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
			return self.(*Pattern).search(args, kwargs)
		}, 0, pattern_search_doc)

	PatternType.Dict["split"] = py.MustNewMethod("split",
		func(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
			return self.(*Pattern).split(args, kwargs)
		}, 0, pattern_split_doc)

	PatternType.Dict["findall"] = py.MustNewMethod("findall",
		func(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
			return self.(*Pattern).findall(args, kwargs)
		}, 0, pattern_findall_doc)

	PatternType.Dict["finditer"] = py.MustNewMethod("finditer",
		func(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
			return self.(*Pattern).finditer(args, kwargs)
		}, 0, pattern_finditer_doc)

	PatternType.Dict["sub"] = py.MustNewMethod("sub",
		func(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
			return self.(*Pattern).sub(args, kwargs)
		}, 0, pattern_sub_doc)

	PatternType.Dict["subn"] = py.MustNewMethod("subn",
		func(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
			return self.(*Pattern).subn(args, kwargs)
		}, 0, pattern_subn_doc)

}

type Pattern struct {
	pattern string
	regx    *regexp.Regexp
}

func (p *Pattern) Type() *py.Type {
	return PatternType
}

func (p *Pattern) M__str__() (py.Object, error) {
	return py.String(fmt.Sprintf("%v", p.pattern)), nil
}

func (p *Pattern) M__repr__() (py.Object, error) {
	return p.M__str__()
}

// _compile
func PatternNew(pattern py.Object, flags py.Object) *Pattern {
	flagStr := ""
	if flags != nil {
		flag := int(flags.(py.Int))
		if flag&SRE_FLAG_IGNORECASE > 0 {
			flagStr += "i"
		}
		if flag&SRE_FLAG_MULTILINE > 0 {
			flagStr += "m"
		}
		if flag&SRE_FLAG_DOTALL > 0 {
			flagStr += "s"
		}
		if len(flagStr) > 0 {
			flagStr = "(?" + flagStr + ")"
		}
	}
	str, _ := toString(pattern)
	patternStr := flagStr + str

	return &Pattern{
		pattern: patternStr,
		regx:    regexp.MustCompile(patternStr),
	}
}

const pattern_match_doc = `match(string[, pos[, endpos]]) -> match object or None.
    Matches zero or more characters at the beginning of the string`

func (p *Pattern) match(args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var String py.Object
	var Pos py.Object
	var EndPos py.Object
	var pattern py.Object
	kwlist := []string{"string", "pos", "endpos", "pattern"}
	err = py.ParseTupleAndKeywords(args, kwargs, "|Onn$O:match", kwlist, &String, &Pos, &EndPos, &pattern)
	if err != nil {
		return nil, err
	}
	String, err = fixStringParam(String, pattern, "pattern")
	return MatchNew(p.regx, String, Pos, EndPos), nil
}

const pattern_fullmatch_doc = `fullmatch(string[, pos[, endpos]]) -> match object or None.
    Matches against all of the string`

func (p *Pattern) fullmatch(args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var String py.Object
	var Pos py.Object
	var EndPos py.Object
	var pattern py.Object
	kwlist := []string{"string", "pos", "endpos", "pattern"}
	err = py.ParseTupleAndKeywords(args, kwargs, "|Onn$O:match", kwlist, &String, &Pos, &EndPos, &pattern)
	if err != nil {
		return nil, err
	}
	String, err = fixStringParam(String, pattern, "pattern")
	p.regx.Longest()
	return MatchNew(p.regx, String, Pos, EndPos), nil
}

const pattern_search_doc = `search(string[, pos[, endpos]]) -> match object or None.
    Scan through string looking for a match, and return a corresponding
    match object instance. Return None if no position in the string matches.`

func (p *Pattern) search(args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var String py.Object
	var Pos py.Object
	var EndPos py.Object
	var pattern py.Object
	kwlist := []string{"string", "pos", "endpos", "pattern"}
	err = py.ParseTupleAndKeywords(args, kwargs, "|Onn$O:search", kwlist, &String, &Pos, &EndPos, &pattern)
	if err != nil {
		return nil, err
	}
	String, err = fixStringParam(String, pattern, "pattern")
	return MatchNew(p.regx, String, Pos, EndPos), nil
}

const pattern_split_doc = `split(string[, maxsplit = 0])  -> list.
    Split string by the occurrences of pattern.`

func (p *Pattern) split(args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var String py.Object
	var MaxSplit py.Object
	var String2 py.Object
	kwlist := []string{"string", "maxsplit", "source"}
	err = py.ParseTupleAndKeywords(args, kwargs, "|On$O:split", kwlist, &String, &MaxSplit, &String2)
	if err != nil {
		return nil, err
	}
	String, err = fixStringParam(String, String2, "source")
	if err != nil {
		return nil, err
	}
	str, isBytes := toString(String)
	maxSplit := -1
	if MaxSplit != nil {
		maxSplit = int(MaxSplit.(py.Int))
		if maxSplit == 0 {
			maxSplit = -1
		} else {
			maxSplit++
		}
	}
	result := p.regx.Split(str, maxSplit)
	if isBytes {
		items := make(py.Tuple, len(result))
		for i, item := range result {
			items[i] = py.Bytes(item)
		}
		return py.NewListFromItems(items), nil
	}
	return py.NewListFromStrings(result), nil
}

const pattern_findall_doc = `findall(string[, pos[, endpos]]) -> list.
   Return a list of all non-overlapping matches of pattern in string.`

func (p *Pattern) findall(args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var String py.Object
	var Pos py.Object
	var EndPos py.Object
	var String2 py.Object
	kwlist := []string{"string", "pos", "endpos", "source"}
	err = py.ParseTupleAndKeywords(args, kwargs, "|Onn$O:findall", kwlist, &String, &Pos, &EndPos, &String2)
	if err != nil {
		return nil, err
	}
	String, err = fixStringParam(String, String2, "source")
	str, isBytes := toString(String)
	pos := 0
	end := len(str)
	if Pos != nil {
		pos = int(Pos.(py.Int))
	}
	if EndPos != nil {
		end = int(EndPos.(py.Int))
	}
	result := p.regx.FindAllStringSubmatch(str[pos:end], -1)
	tuple := make([]py.Object, len(result))
	for i, list := range result {
		if len(list) > 1 {
			list = list[1:]
			t := make(py.Tuple, len(list))
			if isBytes {
				for j, s := range list {
					t[j] = py.Bytes(s)
				}
			} else {
				for j, s := range list {
					t[j] = py.String(s)
				}
			}
			tuple[i] = t
		} else if len(list) == 1 {
			if isBytes {
				tuple[i] = py.Bytes(list[0])
			} else {
				tuple[i] = py.String(list[0])
			}
		}
	}
	return py.NewListFromItems(tuple), nil
}

const pattern_finditer_doc = `finditer(string[, pos[, endpos]]) -> iterator.
    Return an iterator over all non-overlapping matches for the
    RE pattern in string. For each match, the iterator returns a
    match object.`

func (p *Pattern) finditer(args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	//PyObject* scanner;
	//PyObject* search;
	//PyObject* iterator;
	//
	//scanner = pattern_scanner(pattern, args, kw);
	//if (!scanner)
	//return NULL;
	//
	//search = PyObject_GetAttrString(scanner, "search");
	//Py_DECREF(scanner);
	//if (!search)
	//return NULL;
	//
	//iterator = PyCallIter_New(search, Py_None);
	//Py_DECREF(search);
	//
	//return iterator;
	return nil, nil
}

func (p *Pattern) scanner(args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	//var String py.Object
	//var Pos py.Object
	//var EndPos py.Object
	//var String2 py.Object
	//kwlist := []string{"string", "pos", "endpos", "source"}
	//err = py.ParseTupleAndKeywords(args, kwargs, "|Onn$O:scanner", kwlist, &String, &Pos, &EndPos, &String2)
	//if err != nil {
	//	return nil, err
	//}
	//String, err = fixStringParam(String, String2, "source")
	//str := string(String.(py.String))
	//pos := 0
	//end := len(str)
	//if Pos != nil {
	//	pos = int(Pos.(py.Int))
	//}
	//if EndPos != nil {
	//	end = int(EndPos.(py.Int))
	//}
	//result := p.regx.FindAllStringSubmatch(str[pos:end], -1)
	//tuple := make([]py.Object, len(result))
	//for i, list := range result {
	//	if len(list) > 1 {
	//		t := make(py.Tuple, len(list))
	//		for j, s := range list {
	//			t[j] = py.String(s)
	//		}
	//		tuple[i] = t
	//	} else if len(list) == 1 {
	//		tuple[i] = py.String(list[0])
	//	}
	//}
	//return py.NewListFromItems(tuple), nil
	return nil, nil
}

const pattern_sub_doc = `sub(repl, string[, count = 0]) -> newstring.
    Return the string obtained by replacing the leftmost non-overlapping
    occurrences of pattern in string by the replacement repl.`

func (p *Pattern) sub(args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var repl py.Object
	var String py.Object
	var Count py.Object
	kwlist := []string{"repl", "string", "count"}
	err = py.ParseTupleAndKeywords(args, kwargs, "OO|n:sub", kwlist, &repl, &String, &Count)
	if err != nil {
		return nil, err
	}

	res, _, err = p.subx(repl, String, Count)
	return res, err
}

const pattern_subn_doc = `subn(repl, string[, count = 0]) -> (newstring, number of subs)
    Return the tuple (new_string, number_of_subs_made) found by replacing
    the leftmost non-overlapping occurrences of pattern with the
    replacement repl.`

func (p *Pattern) subn(args py.Tuple, kwargs py.StringDict) (res py.Object, err error) {
	var repl py.Object
	var String py.Object
	var Count py.Object
	kwlist := []string{"repl", "string", "count"}
	err = py.ParseTupleAndKeywords(args, kwargs, "OO|n:sub", kwlist, &repl, &String, &Count)
	if err != nil {
		return nil, err
	}
	var num int
	res, num, err = p.subx(repl, String, Count)
	return py.Tuple{res, py.Int(num)}, err
}

func (p *Pattern) subx(repl, String, Count py.Object) (res py.Object, num int, err error) {
	var method func(s string) string
	str, isBytes := toString(String)

	switch val := repl.(type) {
	case *py.Function:
		method = func(s string) string {
			ret, e := py.Call(val, py.Tuple{&Match{
				groups_: py.Tuple{py.String(s)},
			}}, nil)
			err = e
			if ret == nil {
				return ""
			}
			return string(ret.(py.String))
		}
	default:
		replStr, _ := toString(repl)
		b := bytes.NewBufferString(replStr)
		o, err := parser.DecodeEscape(b, false)
		if err == nil {
			replStr = o.String()
		}
		repl_ := replStr
		//repl_, err := strconv.Unquote("\"" + replStr + "\"")
		//if err != nil {
		//	repl_ = replStr
		//}
		method = func(s string) string {
			return p.regx.ReplaceAllString(s, repl_)
		}
	}
	count := -1
	num = 0
	if Count != nil {
		count = int(Count.(py.Int))
		if count == 0 {
			count = -1
		}
	}
	result := p.regx.ReplaceAllStringFunc(str, func(s string) string {
		if count == 0 {
			return s
		}
		count--
		num++
		return method(s)
	})
	if isBytes {
		return py.Bytes(result), num, err
	}
	return py.String(result), num, err
}
