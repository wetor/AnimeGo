// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package re implements the Python Regular Expression module.
package re

import (
	"github.com/go-python/gpython/py"
)

func Init() {
	methods := []*py.Method{
		py.MustNewMethod("match", match, 0, "match(pattern, string)"),
		py.MustNewMethod("search", search, 0, "search(pattern, string)"),
		py.MustNewMethod("sub", sub, 0, "sub(pattern, repl, string[, count=0])"),
		py.MustNewMethod("split", split, 0, "split(pattern, string[, maxsplit=0])"),
		py.MustNewMethod("findall", findAll, 0, "findall(pattern, string[, pos[, endpos]])"),
		py.MustNewMethod("compile", compile, 0, "compile(pattern)"),
	}

	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "re",
			Doc:  "Regular Expression interfaces",
		},
		Methods: methods,
	})
}

func match(self py.Object, arg py.Tuple) (py.Object, error) {
	pattern := string(arg[0].(py.String))
	return _compile(pattern).match(arg[1])
}

func search(self py.Object, arg py.Tuple) (py.Object, error) {
	pattern := string(arg[0].(py.String))
	return _compile(pattern).search(arg[1])
}

func sub(self py.Object, arg py.Tuple) (py.Object, error) {
	pattern := string(arg[0].(py.String))
	return _compile(pattern).sub(arg[1:])
}

func split(self py.Object, arg py.Tuple) (py.Object, error) {
	pattern := string(arg[0].(py.String))
	return _compile(pattern).split(arg[1:])
}

func findAll(self py.Object, arg py.Tuple) (py.Object, error) {
	pattern := string(arg[0].(py.String))
	return _compile(pattern).findAll(arg[1:])
}

func compile(self py.Object, arg py.Tuple) (py.Object, error) {
	var pattern string
	if v, ok := arg[0].(py.String); ok {
		pattern = string(v)
	} else if v, ok := arg[0].(py.Bytes); ok {
		pattern = string(v)
	}
	return _compile(pattern), nil
}
