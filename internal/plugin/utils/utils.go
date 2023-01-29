package utils

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/errors"
)

func CheckParams(paramsSchema [][]string, params models.Object) {
	for _, schema := range paramsSchema {
		if schema[0] == "required" {
			has := false
			for key := range params {
				if key == schema[1] {
					has = true
					break
				}
			}
			if !has {
				errors.NewAniError("参数缺少: " + schema[1]).TryPanic()
			}
		}
	}

	for key := range params {
		has := false
		for _, schema := range paramsSchema {
			if key == schema[1] {
				has = true
				break
			}
		}
		if !has {
			errors.NewAniError("参数多余: " + key).TryPanic()
		}
	}
}

func CheckResult(resultSchema [][]string, result any) {
	resultMap, ok := result.(models.Object)
	if !ok {
		errors.NewAniError("返回类型错误").TryPanic()
	}

	for _, schema := range resultSchema {
		if schema[0] == "required" {
			has := false
			for key := range resultMap {
				if key == schema[1] {
					has = true
					break
				}
			}
			if !has {
				errors.NewAniError("返回值缺少: " + schema[1]).TryPanic()
			}
		}
	}

	for key := range resultMap {
		has := false
		for _, schema := range resultSchema {
			if key == schema[1] {
				has = true
				break
			}
		}
		if !has {
			errors.NewAniError("返回值多余: " + key).TryPanic()
		}
	}
}
