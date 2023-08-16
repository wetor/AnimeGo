package dirdb

import (
	"path"
	"reflect"
)

type File struct {
	DB   DB
	Dir  string
	Ext  string
	File string
}

func NewFile(file string) *File {
	db := reflect.New(reflect.TypeOf(DefaultDB)).Elem()
	db.Set(reflect.ValueOf(DefaultDB))
	return &File{
		File: file,
		Dir:  path.Dir(file),
		Ext:  path.Ext(file),
		DB:   db.Interface().(DB),
	}
}

func (f *File) Open() error {
	return f.DB.Open(f.File)
}

func (f *File) Close() error {
	return f.DB.Close()
}
