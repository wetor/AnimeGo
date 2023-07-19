// Copyright 2022 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package re_test

import (
	"testing"

	"github.com/go-python/gpython/pytest"
	_ "github.com/wetor/AnimeGo/third_party/gpython/py"
	_ "github.com/wetor/AnimeGo/third_party/gpython/stdlib/re"
)

func TestRe(t *testing.T) {
	pytest.RunTests(t, "tests")
}
