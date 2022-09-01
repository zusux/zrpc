package zrpc

import (
	"github.com/go-redis/redis/v8"
	"github.com/knadh/koanf"
	"github.com/sirupsen/logrus"
	"github.com/zusux/zrpc/internal"
	"gorm.io/gorm"
	"log"
	"sync"
)

var once sync.Once
var K *koanf.Koanf
var server  *internal.Server
var db *gorm.DB
var m sync.Mutex
var zLog *logrus.Logger
var red *redis.Client


func Init()  {
	once.Do(func() {
		server = InitConfig()
		zLog = server.Log()
		var err error
		db,err = server.GetDb()
		if err != nil{
			zLog.Errorf("[zrpc] GetDb error:%s", err)
		}
		red,err = server.GetRedis("default_redis")
		if err != nil{
			zLog.Errorf("[zrpc] GetRedis error:%s", err)
		}
		K = internal.K
	})
}

func InitConfig() *internal.Server {
	s, err := internal.NewServer()
	if err != nil {
		log.Fatalf("[zrpc] init server error: %v",err)
	}
	return s
}



// GetDb 系统默认Db
func GetDb() *gorm.DB {
	return db.Session(&gorm.Session{})
}

// NewDb 获取指定Db
func NewDb(section string) (*gorm.DB,error) {
	return server.GetDbBySec(section)
}

// NewLog 获取日志
func NewLog(section string) *logrus.Logger {
	return server.GetLog(section)
}
// Log 获取日志
func Log() *logrus.Logger {
	return zLog
}

// Redis redis
func Redis() *redis.Client {
	return red
}

func GetConf() *internal.Server {
	return server
}

func NewError(code int, message string) *internal.Error {
	return internal.NewError(code,message)
}