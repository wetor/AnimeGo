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

func GetDataFile(name string, file string) *os.File {
	testdata := GetDataPath(name, file)
	f, err := os.Open(testdata)
	if err != nil {
		panic(err)
	}
	return f
}

func GetData(name string, file string) []byte {
	f := GetDataFile(name, file)
	defer f.Close()
	d, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return d
}
