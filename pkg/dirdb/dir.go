package dirdb

import (
	"github.com/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
	"os"
	"path"
	"path/filepath"
)

type Dir struct {
	dir string
}

func Open(dir string) (*Dir, error) {
	exist, isDir := utils.IsExistDir(dir)
	if !exist {
		return nil, errors.Errorf("文件夹不存在: %s", dir)
	}
	if !isDir {
		return nil, errors.Errorf("不是文件夹: %s", dir)
	}
	return &Dir{
		dir: xpath.P(dir),
	}, nil
}

func (d Dir) Scan() ([]*File, error) {
	var files []*File
	dirs, err := os.ReadDir(d.dir)
	if err != nil {
		return nil, err
	}
	for _, info := range dirs {
		if !info.IsDir() && InExt(path.Ext(info.Name())) {
			file := path.Join(d.dir, xpath.P(info.Name()))
			files = append(files, NewFile(file))
		}
	}
	return files, nil
}

func (d Dir) ScanAll() ([]*File, error) {
	var files []*File
	err := filepath.Walk(d.dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && InExt(path.Ext(p)) {
			files = append(files, NewFile(xpath.P(p)))
		}
		return nil
	})
	if err != nil {
		return files, err
	}
	return files, nil
}
