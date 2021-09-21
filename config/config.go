package config

import (
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/zusux/zrpc/env"
	"github.com/zusux/zrpc/micro/zetcd"
	"github.com/zusux/zrpc/model/zdb"
	"github.com/zusux/zrpc/model/zredis"
	"github.com/zusux/zrpc/zlog"
	"log"
)

var config Config

type Config struct {
	Mysql *zdb.Mysql
	Log *zlog.Log
	Redis *zredis.Redis
	Etcd  *zetcd.Etcd
	C *map[string]interface{}
}

func InitConfig() *Config{
	c := env.LoadEnv()
	lg,err :=initLogConfig(c)
	if err != nil{
		log.Fatalln(err)
	}else{
		config.Log = lg
	}
	//初始化日志
	config.GetLog().Zlog().WithField("init","config").Info("开启日志成功")
	sql,err :=initMysqlConfig(c)
	if err != nil{
		config.GetLog().Zlog().WithField("init","mysql").Warn(err)
	}else{
		if sql.Host != "" && sql.Port >0 && sql.Username != "" && sql.Database != ""{
			config.Mysql = sql
			//设置日志模式
			config.Mysql.SetLoger(config.GetLog().Zlog())
			if sql.Debug{
				config.GetLog().Zlog().WithField("init","mysql").Warn("mysql 开启debug模式")
			}
			config.Mysql.NewConnection()
			config.GetLog().Zlog().WithField("init","mysql").Info("mysql 连接成功")
		}else{
			config.GetLog().Zlog().WithField("init","mysql").Warn("mysql子项 未配置")
		}
	}

	rds,err :=initRedisConfig(c)
	if err != nil{
		config.GetLog().Zlog().WithField("init","redis").Warn(err)
	}else{
		config.Redis = rds
		config.GetLog().Zlog().WithField("init","redis").Info("初始化redis配置成功")
	}

	etc,err :=initEtcdConfig(c)
	if err != nil{
		config.GetLog().Zlog().WithField("init","etcd").Error(err)
	}else{
		config.Etcd = etc
		config.GetLog().Zlog().WithField("init","etcd").Info("初始化Etcd配置成功")
	}
	config.C = c
	return &config
}

func GetConfig() *Config {
	return &config
}

func initLogConfig(config *map[string]interface{}) (*zlog.Log,error) {
	logInter,ok := (*config)["log"]
	var log zlog.Log
	if ok {
		err := mapstructure.Decode(logInter,&log)
		return &log,err
	}else{
		return &log,errors.New("log 未配置")
	}
}

func initMysqlConfig(config *map[string]interface{}) (*zdb.Mysql,error) {
	mysqlInter,ok := (*config)["mysql"]
	var mysql zdb.Mysql
	if ok {
		err := mapstructure.Decode(mysqlInter,&mysql)
		return &mysql,err
	}else{
		return &mysql,errors.New("mysql 未配置")
	}
}

func initRedisConfig(config *map[string]interface{}) (*zredis.Redis,error) {
	redisInter,ok := (*config)["redis"]
	var redis zredis.Redis
	if ok {
		err := mapstructure.Decode(redisInter,&redis)
		return &redis,err
	}else{
		return &redis,errors.New("redis 未配置")
	}
}

func initEtcdConfig(config *map[string]interface{}) (*zetcd.Etcd,error) {
	etcdInter,ok := (*config)["etcd"]
	var etcd zetcd.Etcd
	if ok {
		err := mapstructure.Decode(etcdInter,&etcd)
		return &etcd,err
	}else{
		return &etcd,errors.New("etcd 未配置")
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
func (c *Config) GetEtcd() *zetcd.Etcd {
	return c.Etcd
}

