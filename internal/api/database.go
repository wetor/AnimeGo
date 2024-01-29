package api

type Database interface {
	Scan() error
	IsExist(data any) bool
	Add(data any) error
	Delete(data any) error
}
