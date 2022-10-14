package errors

import (
	"fmt"
	"testing"
)

func TestNewAniError(t *testing.T) {
	err := NewAniError("测试错误")
	fmt.Println(err)
}
