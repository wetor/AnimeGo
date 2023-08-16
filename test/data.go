package test

import (
	"io"
	"os"
	"path"
	"runtime"
)

func GetDataPath(name string, file string) string {
	_, currFile, _, _ := runtime.Caller(0)
	dir := path.Dir(currFile)
	testdata := path.Join(dir, "testdata", name, file)
	return testdata
}

func GetDataFile(name string, file string) (*os.File, error) {
	testdata := GetDataPath(name, file)
	f, err := os.Open(testdata)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func GetData(name string, file string) ([]byte, error) {
	f, err := GetDataFile(name, file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return d, nil
}
