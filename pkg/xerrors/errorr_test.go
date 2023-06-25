package xerrors_test

import (
	errs "errors"
	"fmt"
	"testing"

	"github.com/wetor/AnimeGo/pkg/xerrors"
)

func TestNewAniError(t *testing.T) {
	err := xerrors.NewAniError("测试错误")
	fmt.Println(err)
}

func TestNewAniErrorPanic(t *testing.T) {
	defer xerrors.HandleError(func(err error) {
		fmt.Println("【捕获错误】", err)

	})
	xerrors.NewAniError("").TryPanic()

	xerrors.NewAniErrorD(nil).TryPanic()

	xerrors.NewAniError("测试错误").TryPanic()
}

func TestNewAniErrorPanic2(t *testing.T) {
	defer xerrors.HandleError(func(err error) {
		fmt.Println("【捕获错误1】", err)
	})
	xerrors.NewAniError("").TryPanic()

	func() {
		defer xerrors.HandleError(func(err error) {
			fmt.Println("【捕获错误2】", err)
		})
		xerrors.NewAniErrorD(errs.New("测试errors.New")).TryPanic()
	}()

	xerrors.NewAniError("测试错误").TryPanic()
}
