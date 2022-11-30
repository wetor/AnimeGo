package errors

import (
	"fmt"
	"testing"
)

func TestNewAniError(t *testing.T) {
	err := NewAniError("测试错误")
	fmt.Println(err)
}

func TestNewAniErrorPanic(t *testing.T) {
	NewAniError("").TryPanic()

	NewAniErrorD(nil).TryPanic()

	NewAniError("测试错误").TryPanic()
}
