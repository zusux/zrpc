package env

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	flag "github.com/spf13/pflag"
	"github.com/zusux/zrpc/code"
	"github.com/zusux/zrpc/zerr"
	"os"
	"path"
	"runtime"
	"time"
)
var K = koanf.New(".")
func LoadToml()  {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		panic(zerr.NewZErr(code.CONFIG_FLAG_USAGED_ERROR,"config flag has usaged"))
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
			panic(zerr.NewZErr(code.CONFIG_FILE_LOADING_ERROR,err.Error()))
		}
	}
	// "time" and "type" may have been loaded from the config file, but
	// they can still be overridden with the values from the command line.
	// The bundled posflag.Provider takes a flagset from the spf13/pflag lib.
	// Passing the Koanf instance to posflag helps it deal with default command
	// line flag values that are not present in conf maps from previously loaded
	// providers.
	if err := K.Load(posflag.Provider(f, ".", K), nil); err != nil {
		panic(zerr.NewZErr(code.CONFIG_LOADING_ERROR,fmt.Sprintf("error loading config: %v", err)))
	}
}

//获取当前文件路径
func GetCurrentDir() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		panic(zerr.NewZErr(code.CONFIG_GET_CURRENT_FILE_ERROR,"Can not get current file info"))
	}
	return path.Dir(file)
}
