package zetcd

type  Server struct {
	Port int `koanf:"port"`
	Publish bool `koanf:"publish"`
}


