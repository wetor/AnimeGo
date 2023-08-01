package dirdb

import (
	"reflect"

	"github.com/wetor/AnimeGo/pkg/xpath"
)

type File struct {
	DB   DB
	Dir  string
	File string
}

func NewFile(file string) *File {
	db := reflect.New(reflect.TypeOf(DefaultDB)).Elem()
	db.Set(reflect.ValueOf(DefaultDB))
	return &File{
		File: file,
		Dir:  xpath.Dir(file),
		DB:   db.Interface().(DB),
	}
}

func (f *File) Open() error {
	return f.DB.Open(f.File)
}

func (f *File) Close() error {
	return f.DB.Close()
}
