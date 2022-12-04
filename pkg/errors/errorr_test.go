package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewAniError(t *testing.T) {
	err := NewAniError("测试错误")
	fmt.Println(err)
}

func TestNewAniErrorPanic(t *testing.T) {
	defer HandleError(func(err error) {
		fmt.Println("【捕获错误】", err)
		panic(err)

	})
	NewAniError("").TryPanic()

	NewAniErrorD(nil).TryPanic()

	NewAniError("测试错误").TryPanic()
}

func TestNewAniErrorPanic2(t *testing.T) {
	defer HandleError(func(err error) {
		fmt.Println("【捕获错误1】", err)
	})
	NewAniError("").TryPanic()

	func() {
		defer HandleError(func(err error) {
			fmt.Println("【捕获错误2】", err)
		})
		NewAniErrorD(errors.New("测试errors.New")).TryPanic()
	}()

	NewAniError("测试错误").TryPanic()
}
