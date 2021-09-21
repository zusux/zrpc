package env

import (
	"io/ioutil"
	"log"
	"os"
	"gopkg.in/yaml.v3"
	"path"
	"runtime"
	"errors"
)
//解析 env.yaml 文件
func LoadEnv() *map[string]interface{}{
	var envFile  = GetCurrentDir() + "/env.yaml"
	if filePath := os.Getenv("ENV_FILE"); filePath != "" {
		envFile = filePath
	}
	f, err := ioutil.ReadFile(envFile)
	if err != nil {
		log.Fatalln("ReadFileError: ",err)
	}
	var config map[string]interface{}
	err = yaml.Unmarshal(f, &config)
	if err != nil {
		log.Fatalln("UnmarshalError:",err)
	}
	return &config
}

//获取当前文件路径
func GetCurrentDir() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New("Can not get current file info"))
	}
	return path.Dir(file)
}

