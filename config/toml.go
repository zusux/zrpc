package config

import (
	"fmt"
	"github.com/zusux/zrpc/code"
	"github.com/zusux/zrpc/env"
	"github.com/zusux/zrpc/micro/zetcd"
	"github.com/zusux/zrpc/model/zdb"
	"github.com/zusux/zrpc/model/zredis"
	"github.com/zusux/zrpc/zerr"
	"github.com/zusux/zrpc/zlog"
)

type toml struct {}

func (t *toml) InitToml() *Config {
	env.LoadToml()
	log := t.initLogConfig()
	mysql := t.initMysqlConfig()
	redis := t.initRedisConfig()
	etcd := t.initEtcdConfig()
	conf := newConfig(log, mysql, redis, etcd)
	//初始化日志
	conf.GetLog().Zlog().WithField("init", "config").Info("开启日志成功")

	if mysql.Host != "" && mysql.Port > 0 && mysql.Username != "" && mysql.Database != "" {
		//设置日志模式
		conf.Mysql.SetLoger(conf.GetLog().Zlog())
		if mysql.Debug {
			conf.GetLog().Zlog().WithField("init", "mysql").Warn("mysql 开启debug模式")
		}
		_,err := conf.Mysql.NewConnection()
		if err != nil{
			conf.GetLog().Zlog().Error(zerr.NewZErr(code.MYSQL_CONNECT_ERROR,err.Error()).String())
			conf.GetLog().Zlog().WithField("init", "mysql").Error("mysql 连接失败")
		}else{
			conf.GetLog().Zlog().WithField("init", "mysql").Info("mysql 连接成功")
		}
	} else {
		conf.GetLog().Zlog().WithField("init", "mysql").Warn("mysql子项 未配置")
	}

	if redis.Host != "" && redis.Port > 0 {
		// redis 初始化
		conf.GetLog().Zlog().WithField("init", "redis").Info("初始化redis配置成功")
	} else {
		conf.GetLog().Zlog().WithField("init", "redis").Warn("redis 未配置")
	}

	if len(etcd.EtcdServerAddress) > 0 {
		conf.GetLog().Zlog().WithField("init", "etcd").Info("初始化Etcd配置成功")
	} else {
		conf.GetLog().Zlog().WithField("init", "etcd").Warn("etcd 未配置")
	}
	return conf
}

func (t *toml) initLogConfig() *zlog.Log {
	var log = zlog.NewLog(
		env.K.String("log.path"),
		env.K.String("log.file"),
		env.K.String("log.format"),
		env.K.Int64("log.age"),
		env.K.Int64("log.rotation"),
	)
	return log
}

func (t *toml) initMysqlConfig() *zdb.Mysql {
	var mysql = zdb.NewMysql(
		env.K.String("mysql.host"),
		env.K.Int("mysql.port"),
		env.K.String("mysql.username"),
		env.K.String("mysql.password"),
		env.K.String("mysql.database"),
		env.K.Bool("mysql.debug"),
	)
	return mysql
}

func (t *toml) initRedisConfig() *zredis.Redis {
	var redis = zredis.NewRedis(
		env.K.String("redis.host"),
		env.K.Int("redis.port"),
		env.K.String("redis.auth"),
	)
	return redis
}

func (t *toml) initEtcdConfig() *zetcd.Etcd {
	var etcd = zetcd.NewEtcd(
		env.K.String("server.name"),
		t.getGrpcServerAddress(),
		env.K.Int64("etcd.dial_timeout"),
		env.K.Int64("etcd.dial_keep_alive"),
		env.K.Strings("etcd.etcd_server_address")...
	)
	return etcd
}

func (t *toml) getGrpcServerAddress() string {
	return fmt.Sprintf("%s:%d", env.K.String("server.grpc.host"), env.K.Int("server.grpc.port"))
}
