package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

//获取当前文件路径
func GetCurrentDir() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Can not get current file info"))
	}
	return filepath.Dir(file)
}

func GetExecutableDir() string {
	fp, err := os.Executable()
	if err != nil{
		panic(errors.New(fmt.Sprintf("executable path not find: %v", err)))
	}
	return filepath.Dir(fp)
}

func GetWdDir() string {
	fd, err := os.Getwd()
	if err != nil{
		panic(errors.New(fmt.Sprintf("pwd dir not find: %v", err)))
	}
	return fd
}
