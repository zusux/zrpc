package util

import (
	"errors"
	"fmt"
	"io/ioutil"
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

func AvailablePath(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err == nil {
		isDir := fileInfo.IsDir()
		if isDir{
			f,err :=ioutil.TempFile(path,"zrpc-*.tmp")
			if err == nil{
				f.Close()
				err = os.Remove(f.Name())
				return true,nil
			}else{
				return false,err
			}
		}
		return false,errors.New(fmt.Sprintf("[zrpc][log] %s is not dir",path))
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(path,os.ModePerm)
		if err != nil{
			return false, err
		}else{
			return true, nil
		}
	}
	return false, err
}