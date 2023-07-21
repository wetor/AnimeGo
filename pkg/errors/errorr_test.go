package errors_test

import (
	errs "errors"
	"fmt"
	"testing"

	"github.com/wetor/AnimeGo/pkg/errors"
)

func TestNewAniError(t *testing.T) {
	err := errors.NewAniError("测试错误")
	fmt.Println(err)
}

func TestNewAniErrorPanic(t *testing.T) {
	defer errors.HandleError(func(err error) {
		fmt.Println("【捕获错误】", err)

	})
	errors.NewAniError("").TryPanic()

	errors.NewAniErrorD(nil).TryPanic()

	errors.NewAniError("测试错误").TryPanic()
}

func TestNewAniErrorPanic2(t *testing.T) {
	defer errors.HandleError(func(err error) {
		fmt.Println("【捕获错误1】", err)
	})
	errors.NewAniError("").TryPanic()

	func() {
		defer errors.HandleError(func(err error) {
			fmt.Println("【捕获错误2】", err)
		})
		errors.NewAniErrorD(errs.New("测试errors.New")).TryPanic()
	}()

	errors.NewAniError("测试错误").TryPanic()
}
