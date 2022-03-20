package config

import (
	"fmt"
	"github.com/zusux/zrpc/env"
	"github.com/zusux/zrpc/net/zdb"
	"github.com/zusux/zrpc/net/zetcd"
	"github.com/zusux/zrpc/net/zredis"
	"github.com/zusux/zrpc/utils"
	"github.com/zusux/zrpc/zerr"
	"github.com/zusux/zrpc/zlog"
	"strings"
)

type toml struct{}

func (t *toml) InitToml() *Config {
	env.LoadToml()
	log := t.initLogConfig()
	mysql := t.initMysqlConfig()
	redis := t.initRedisConfig()
	publishes := t.initPublishConfig(log)
	conf := newConfig(log, mysql, redis, publishes)
	//初始化日志
	log.Zlog().WithField("init", "config").Info("开启日志成功")

	if mysql.Host != "" && mysql.Port > 0 && mysql.Username != "" && mysql.Database != "" {
		//设置日志模式
		conf.Mysql.SetLoger(conf.GetLog().Zlog())
		if mysql.Debug {
			log.Zlog().WithField("init", "mysql").Warn("mysql 开启debug模式")
		}
		_, err := conf.Mysql.NewConnection()
		if err != nil {
			log.Zlog().Error(zerr.NewZErr(zerr.MYSQL_CONNECT_ERROR, err.Error()).String())
			log.Zlog().WithField("init", "mysql").Error("mysql 连接失败")
		} else {
			log.Zlog().WithField("init", "mysql").Info("mysql 连接成功")
		}
	} else {
		log.Zlog().WithField("init", "mysql").Warn("mysql子项 未配置")
	}

	if redis.Host != "" && redis.Port > 0 {
		// redis 初始化
		log.Zlog().WithField("init", "redis").Info("初始化redis配置成功")
	} else {
		log.Zlog().WithField("init", "redis").Warn("redis 未配置")
	}

	if len(env.K.Strings("etcd.etcd_server_address")) > 0 && len(publishes) > 0 {
		log.Zlog().WithField("init", "etcd").Info("初始化Etcd配置成功")
	} else {
		log.Zlog().WithField("init", "etcd").Warn("etcd 未配置")
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

func getRegList() {
	env.K.MapKeys("zrpc")
}

func (t *toml) initPublishConfig(log *zlog.Log) (publishes Pubs) {
	var servers map[string]zetcd.Server
	env.K.Unmarshal("zrpc", &servers)
	ip, err := utils.GetLocalIP()
	if err != nil {
		log.Zlog().WithField("init", "publish").Error("获取本地ip失败, 注册未成功")
		return
	}
	var localIp = ip.String()
	log.Zlog().WithField("init", "publish").Debug(fmt.Sprintf("获取本地 ip: %s", localIp))
	for publishType, server := range servers {
		if server.Publish {
			etcd := t.initEtcdConfig(localIp, publishType, server.Port)
			publishes = append(publishes, etcd)
		}
	}
	return
}

func (t *toml) initEtcdConfig(localIp string, publishType string, port int) *zetcd.Etcd {
	return zetcd.NewEtcd(
		strings.TrimRight(env.K.String("server.name"), "/")+"/"+publishType,
		t.getServerAddr(localIp, port),
		env.K.Int64("etcd.dial_timeout"),
		env.K.Int64("etcd.dial_keep_alive"),
		env.K.Strings("etcd.etcd_server_address")...
	)
}

func (t *toml) getServerAddr(localIp string, port int) string {
	return fmt.Sprintf("%s:%d", localIp, port)
}
