package api

type Database interface {
	IsExist(data any) bool
	Add(data any) error
}
