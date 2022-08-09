package zrpc

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/zusux/zrpc/utils"
	"log"
	"os"
)
func init(){
	LoadEnv()
}
var K = koanf.New(".")
func LoadEnv() {
	// Load yaml config.
	envFilepath := os.Getenv("env.config")
	if envFilepath == "" {
		//默认路径
		envFilepath = utils.GetWdDir()  + "/env.yaml"
	}
	f := file.Provider(envFilepath)
	if err := K.Load(f, yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	// Watch the file and get a callback on change. The callback can do whatever,
	// like re-load the configuration.
	// File provider always returns a nil `event`.
	f.Watch(func(event interface{}, err error) {
		if err != nil {
			log.Printf("watch error: %v", err)
			return
		}

		// Throw away the old config and load a fresh copy.
		log.Println("config changed. Reloading ...")
		K = koanf.New(".")
		K.Load(f, yaml.Parser())
		K.Print()
	})
}