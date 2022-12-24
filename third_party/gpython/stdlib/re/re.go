// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package re implements the Python Regular Expression module.
package re

import (
	"github.com/go-python/gpython/py"
)

func init() {
	methods := []*py.Method{
		py.MustNewMethod("match", match, 0, "match"),
		py.MustNewMethod("search", search, 0, "match"),
		py.MustNewMethod("sub", sub, 0, "match"),
		py.MustNewMethod("split", split, 0, "match"),
		py.MustNewMethod("findall", findAll, 0, "match"),
		py.MustNewMethod("compile", compile, 0, "match"),
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
	return _compile(pattern).findAll(arg[1])
}

func compile(self py.Object, arg py.Object) (py.Object, error) {
	pattern := string(arg.(py.String))
	return _compile(pattern), nil
}
