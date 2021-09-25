package main

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/zusux/zrpc/config"
	"github.com/zusux/zrpc/micro/zetcd"
)

func main()  {
	config.InitConfig()
}

func GetDb() *gorm.DB {
	return config.GetConfig().GetMysql().GetDb()
}

func GetLog() *logrus.Logger {
	return config.GetConfig().GetLog().Zlog()
}

func GetConf() *config.Config {
	return config.GetConfig()
}

func GetEtcd() *zetcd.Etcd {
	return config.GetConfig().GetEtcd()
}
