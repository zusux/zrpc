package env

import (
	"fmt"
	"log"
	"os"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	flag "github.com/spf13/pflag"
	"time"
)
var K = koanf.New(".")
func LoadToml()  {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
	envFilepath := os.Getenv("env.path")
	if envFilepath == ""{
		envFilepath = GetCurrentDir()+"/"+"env.toml"
	}
	// Path to one or more config files to load into koanf along with some config params.
	f.StringSlice("conf", []string{envFilepath}, "path to one or more .toml config files")
	f.String("time", time.Now().Format("2006-01-02 15:04:05"), "app start time")
	f.Parse(os.Args[1:])

	// Load the config files provided in the commandline.
	cFiles, _ := f.GetStringSlice("conf")
	for _, c := range cFiles {
		if err := K.Load(file.Provider(c), toml.Parser()); err != nil {
			log.Fatalf("error loading file: %v", err)
		}
	}
	// "time" and "type" may have been loaded from the config file, but
	// they can still be overridden with the values from the command line.
	// The bundled posflag.Provider takes a flagset from the spf13/pflag lib.
	// Passing the Koanf instance to posflag helps it deal with default command
	// line flag values that are not present in conf maps from previously loaded
	// providers.
	if err := K.Load(posflag.Provider(f, ".", K), nil); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
}
