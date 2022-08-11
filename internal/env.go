package internal

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/zusux/zrpc/internal/util"
	"log"
	"os"
)

var K = koanf.New(".")
func LoadEnv() {
	// Load yaml config.
	f := file.Provider(getConfigFilepath())
	if err := K.Load(f, yaml.Parser()); err != nil {
		log.Fatalf("[zrpc] error loading config: %v", err)
	}
	// Watch the file and get a callback on change. The callback can do whatever,
	// like re-load the configuration.
	// File provider always returns a nil `event`.
	f.Watch(func(event interface{}, err error) {
		if err != nil {
			log.Printf("[zrpc] config file watch error: %v", err)
			return
		}
		// Throw away the old config and load a fresh copy.
		log.Println("[zrpc] config changed. Reloading ...")
		K = koanf.New(".")
		K.Load(f, yaml.Parser())
		log.Println(K.Sprint())
	})
}

func getConfigFilepath() string  {
	siteMode := os.Getenv("site_mode")
	if siteMode == ""{
		siteMode = "dev"
	}
	return fmt.Sprintf("%s/config/%s.yaml", util.GetWdDir(),siteMode)
}