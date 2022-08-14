package internal

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/zusux/zrpc/internal/util"
	"log"
)

var K = koanf.New(".")
func LoadEnv(siteMode string) {
	// Load yaml config.
	f := file.Provider(getConfigFilepath(siteMode))
	if err := K.Load(f, yaml.Parser()); err != nil {
		log.Fatalf("[zrpc] error loading config: %v", err.Error())
	}
}

func getConfigFilepath(siteMode string) string  {
	return fmt.Sprintf("%s/config/%s.yaml", util.GetWdDir(),siteMode)
}