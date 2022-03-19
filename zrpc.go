package zrpc

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/zusux/zrpc/config"
	"github.com/zusux/zrpc/net/zetcd"
)

func Init() {
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

func GetPublishes() []*zetcd.Etcd {
	return config.GetConfig().GetPublishes()
}
