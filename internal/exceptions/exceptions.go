package exceptions

import "github.com/pkg/errors"

type ExistError interface {
	Exist() bool
}

func IsExist(err error) bool {
	if e, ok := err.(ExistError); ok && e.Exist() {
		return true
	} else if e, ok = errors.Cause(err).(ExistError); ok && e.Exist() {
		return true
	}
	return false
}

type NotFoundError interface {
	NotFound() bool
}

func IsNotFound(err error) bool {
	if e, ok := err.(NotFoundError); ok && e.NotFound() {
		return true
	} else if e, ok = errors.Cause(err).(NotFoundError); ok && e.NotFound() {
		return true
	}
	return false
}

type ParseFailedError interface {
	ParseFailed() bool
}

func IsParseFailed(err error) bool {
	if e, ok := err.(ParseFailedError); ok && e.ParseFailed() {
		return true
	} else if e, ok = errors.Cause(err).(ParseFailedError); ok && e.ParseFailed() {
		return true
	}
	return false
}
