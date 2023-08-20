package dirdb

import (
	"os"

	"github.com/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/json"
)

type DB interface {
	Open(string) error
	Unmarshal(value any) error
	Marshal(value any) error
	Get(string, any) error
	Set(string, any) error
	Close() error
}

type JsonDB struct {
	file string
}

func (d *JsonDB) Open(file string) error {
	d.file = file
	return nil
}

func (d *JsonDB) Close() error {
	return nil
}

func (d *JsonDB) Unmarshal(value any) error {
	data, err := os.ReadFile(d.file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, value)
	if err != nil {
		return err
	}
	return nil
}

func (d *JsonDB) Marshal(value any) error {
	data, err := json.MarshalIndent(value)
	if err != nil {
		return err
	}
	err = os.WriteFile(d.file, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (d *JsonDB) Get(key string, value any) error {
	return errors.New("不支持的操作")
}

func (d *JsonDB) Set(key string, value any) error {
	return errors.New("不支持的操作")
}
