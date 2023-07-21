// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package re_test

import (
	"testing"

	"github.com/go-python/gpython/pytest"
	"github.com/wetor/AnimeGo/third_party/gpython/stdlib/re"
)

func TestRe(t *testing.T) {
	re.Init()
	pytest.RunScript(t, "testdata/test.py")
}
