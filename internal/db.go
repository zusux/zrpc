package internal

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zusux/zrpc/internal/zerr"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)


type Mysql struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Prefix   string
	Debug    bool
	Logger   *logrus.Logger
	Db  *gorm.DB
	MaxIdleConns int //最大空闲连接数
	MaxOpenConns int //最大连接数
	ConnMaxLifetime int //分钟
	SectionLog string
}

func NewMysql(
	host string,
	port int,
	username string,
	password string,
	database string,
	prefix string,
	debug bool,
	log *logrus.Logger,
	sectionLog string,
) *Mysql {
	return &Mysql{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
		Prefix: prefix,
		Debug:    debug,
		Logger: log,
		SectionLog: sectionLog,
	}
}


func (m *Mysql) NewConnection() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(m.Logger.Writer(), "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,   // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:      false,   // 禁用彩色打印
		},
	)
	conf := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: m.Prefix,
			SingularTable: true,
		},
		Logger: newLogger,
	}
	return gorm.Open(mysql.Open(m.GetDns()), &conf)
}

func (m *Mysql) SetLogger(logger *logrus.Logger) {
	m.Logger = logger
}

func (m *Mysql) GetLevel()  logger.LogLevel {
	if m.Debug{
		return logger.Info
	}
	return logger.Warn
}

func (m *Mysql) GetDns() string {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.Database,
	)
	if m.Debug{
		m.Logger.Info(dns)
	}
	return dns
}

//GetDb 获取数据库
func (m *Mysql) GetDb() (db *gorm.DB,err error) {
	if m.Db == nil{
		db,err = m.NewConnection()
		if err != nil{
			zrr := NewError(zerr.MYSQL_CONNECT_ERROR,err.Error())
			m.Logger.Errorf("[zrpc][db] connect error:%s",zrr.String())
			return db, err
		}
		m.Db = db
	}
	return m.Db,err
}

// Close 关闭数据库
func (m *Mysql) Close() error {
	mdb,err := m.Db.DB()
	if err != nil{
		m.Logger.Errorf("[zrpc][Db] close error:%s",err.Error())
		return err
	}
	return mdb.Close()
}

func (m *Mysql) getMaxIdleConns() int {
	if m.MaxIdleConns == 0{
		m.MaxIdleConns = 1
	}
	return m.MaxIdleConns
}

func (m *Mysql) getMaxOpenConns() int {
	if m.MaxOpenConns == 0{
		m.MaxOpenConns = 2
	}
	return m.MaxOpenConns
}

func (m *Mysql) getConnMaxLifetime() int {
	if m.ConnMaxLifetime == 0{
		m.ConnMaxLifetime = 60
	}
	return m.ConnMaxLifetime
}