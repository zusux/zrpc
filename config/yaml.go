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

type yaml struct {}

func (y *yaml) InitYaml() *Config{
	var conf  *Config
	c := env.LoadYaml()
	lg,err :=y.initLogConfig(c)
	if err != nil{
		log.Fatalln(err)
	}else{
		conf.Log = lg
	}
	//初始化日志
	conf.GetLog().Zlog().WithField("init","config").Info("开启日志成功")
	sql,err :=y.initMysqlConfig(c)
	if err != nil{
		conf.GetLog().Zlog().WithField("init","mysql").Warn(err)
	}else{
		if sql.Host != "" && sql.Port >0 && sql.Username != "" && sql.Database != ""{
			conf.Mysql = sql
			//设置日志模式
			conf.Mysql.SetLoger(conf.GetLog().Zlog())
			if sql.Debug{
				conf.GetLog().Zlog().WithField("init","mysql").Warn("mysql 开启debug模式")
			}
			conf.Mysql.NewConnection()
			conf.GetLog().Zlog().WithField("init","mysql").Info("mysql 连接成功")
		}else{
			conf.GetLog().Zlog().WithField("init","mysql").Warn("mysql子项 未配置")
		}
	}

	rds,err :=y.initRedisConfig(c)
	if err != nil{
		conf.GetLog().Zlog().WithField("init","redis").Warn(err)
	}else{
		conf.Redis = rds
		conf.GetLog().Zlog().WithField("init","redis").Info("初始化redis配置成功")
	}

	etc,err :=y.initEtcdConfig(c)
	if err != nil{
		conf.GetLog().Zlog().WithField("init","etcd").Error(err)
	}else{
		conf.Etcd = etc
		conf.GetLog().Zlog().WithField("init","etcd").Info("初始化Etcd配置成功")
	}
	conf.C = c
	return conf
}

func (y *yaml)GetConfig() *Config {
	return config
}

func (y *yaml)initLogConfig(config *map[string]interface{}) (*zlog.Log,error) {
	logInter,ok := (*config)["log"]
	var log zlog.Log
	if ok {
		err := mapstructure.Decode(logInter,&log)
		return &log,err
	}else{
		return &log,errors.New("log 未配置")
	}
}

func (y *yaml)initMysqlConfig(config *map[string]interface{}) (*zdb.Mysql,error) {
	mysqlInter,ok := (*config)["mysql"]
	var mysql zdb.Mysql
	if ok {
		err := mapstructure.Decode(mysqlInter,&mysql)
		return &mysql,err
	}else{
		return &mysql,errors.New("mysql 未配置")
	}
}

func (y *yaml)initRedisConfig(config *map[string]interface{}) (*zredis.Redis,error) {
	redisInter,ok := (*config)["redis"]
	var redis zredis.Redis
	if ok {
		err := mapstructure.Decode(redisInter,&redis)
		return &redis,err
	}else{
		return &redis,errors.New("redis 未配置")
	}
}

func (y *yaml)initEtcdConfig(config *map[string]interface{}) (*zetcd.Etcd,error) {
	etcdInter,ok := (*config)["etcd"]
	var etcd zetcd.Etcd
	if ok {
		err := mapstructure.Decode(etcdInter,&etcd)
		return &etcd,err
	}else{
		return &etcd,errors.New("etcd 未配置")
	}
}

