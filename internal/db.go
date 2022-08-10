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
	logger   *logrus.Logger
	zdb  *gorm.DB
}

// GetDbBySec 指定DB
func GetDbBySec(section string) (*gorm.DB, error) {
	return GetMysql(section).GetDb()
}

func GetMysql(section string) *Mysql {
	host := fmt.Sprintf("%s.host", section)
	port := fmt.Sprintf("%s.port", section)
	username := fmt.Sprintf("%s.username", section)
	passwd := fmt.Sprintf("%s.password", section)
	database := fmt.Sprintf("%s.database", section)
	debug := fmt.Sprintf("%s.debug", section)
	prefix := fmt.Sprintf("%s.prefix", section)
	var mysql = NewMysql(
		K.String(host),
		K.Int(port),
		K.String(username),
		K.String(passwd),
		K.String(database),
		K.String(prefix),
		K.Bool(debug),
		Log(),
	)
	return mysql
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
) *Mysql {
	return &Mysql{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
		Prefix: prefix,
		Debug:    debug,
		logger: log,
	}
}

func (m *Mysql) NewConnection() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(m.logger.Writer(), "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
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
	db, err := gorm.Open(mysql.Open(m.GetDns()), &conf)
	if err != nil {
		return db, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
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
		m.logger.Info(dns)
	}
	return dns
}

//获取数据库
func (m *Mysql) GetDb() (*gorm.DB,error) {
	//m.NewConnection()
	if m.zdb == nil{
		db,err := m.NewConnection()
		if err != nil{
			zrr := NewError(zerr.MYSQL_CONNECT_ERROR,err.Error())
			m.logger.Warn(zrr.String())
		}
		m.zdb = db
	}
	mdb,err := m.zdb.DB()
	if  err != nil || mdb.Ping() != nil {
		db,err := m.NewConnection()
		if err != nil{
			zrr := NewError(zerr.MYSQL_CONNECT_ERROR,err.Error())
			m.logger.Error(zrr.String())
			m.zdb = db
		}
	}
	return m.zdb,nil
}

//关闭数据库
func (m *Mysql) Close() error {
	mdb,err := m.zdb.DB()
	if err != nil{
		return err
	}
	return mdb.Close()
}