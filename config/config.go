package config

import (
	"github.com/zusux/zrpc/net/zdb"
	"github.com/zusux/zrpc/net/zetcd"
	"github.com/zusux/zrpc/net/zredis"
	"github.com/zusux/zrpc/zlog"
)

var conf *Config

type Pubs []*zetcd.Etcd

func (p *Pubs) Register()  {
	for _,v :=range *p{
		v.Register()
	}
}
func (p *Pubs) UnRegister()  {
	for _,v :=range *p{
		v.UnRegister()
	}
}

type Config struct {
	Mysql     *zdb.Mysql
	Log       *zlog.Log
	Redis     *zredis.Redis
	Pubs  Pubs
	C         *map[string]interface{}
}

func InitConfig() {
	initConfigByToml()
}

func GetConfig() *Config {
	return conf
}

func initConfigByToml() {
	t := &toml{}
	conf = t.InitToml()
}

func newConfig(log *zlog.Log, mysql *zdb.Mysql, redis *zredis.Redis, publishes []*zetcd.Etcd) *Config {
	return &Config{
		Log:       log,
		Mysql:     mysql,
		Redis:     redis,
		Pubs: publishes,
	}
}

func (c *Config) GetMysql() *zdb.Mysql {
	return c.Mysql
}
func (c *Config) GetRedis() *zredis.Redis {
	return c.Redis
}
func (c *Config) GetLog() *zlog.Log {
	return c.Log
}
func (c *Config) GetPubs() Pubs {
	return c.Pubs
}
