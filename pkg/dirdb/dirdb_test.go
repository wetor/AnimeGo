package dirdb_test

import (
	"fmt"
	"testing"

	"github.com/wetor/AnimeGo/pkg/dirdb"
)

func TestDirDB_ReadAll(t *testing.T) {
	db, err := dirdb.Open("testdata/test")
	if err != nil {
		panic(err)
	}
	files, err := db.ScanAll()
	if err != nil {
		panic(err)
	}

	data := struct {
		Path string         `json:"path"`
		ID   int            `json:"id"`
		Info map[string]any `json:"info"`
	}{}
	for i, f := range files {
		_ = f.Open()
		data.Path = f.File
		data.ID = i * 1000
		data.Info = map[string]any{
			"path": f.File,
			"id":   i * 2000,
		}
		err = f.DB.Marshal(data)
		if err != nil {
			panic(err)
		}
		err = f.DB.Unmarshal(&data)
		if err != nil {
			panic(err)
		}
		fmt.Println(data)
	}

}
