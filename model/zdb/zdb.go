package zdb

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"time"
)

var db *gorm.DB

type Mysql struct{
	Host string
	Port int
	Username string
	Password string
	Database string
	Debug bool
	logger *logrus.Logger
}

func NewMysql(host string, port int, username string, password string, database string) *Mysql {
	return &Mysql{Host: host, Port: port, Username: username, Password: password, Database: database}
}

func (m *Mysql) NewConnection() *gorm.DB  {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", m.Username, m.Password, m.Host, m.Port, m.Database)
	var err error
	db, err = gorm.Open("mysql", dsn)

	if err != nil {
		panic("mysql连接数据库失败, error=" + err.Error())
	}
	//设置 debug 模式
	db.LogMode(m.Debug)

	if m.logger != nil{
		db.SetLogger(m.logger)
	}
	//db.Callback().Query().Register("gorm:querySql", callback)
	//db.Callback().Delete().Register("gorm:deleteSql", callback)
	//// 监听update方法
	//db.Callback().Update().Register("gorm:updateSql", callback)
	//// 监听create方法
	//db.Callback().Create().Register("gorm:insertSql", callback)
	//db.SetLogger(GetLog())
	db.DB().SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(100)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	db.DB().SetConnMaxLifetime(time.Hour)
	db.SingularTable(true)
	return db
}
// 设置日志模式
func (m *Mysql) SetLoger(log *logrus.Logger)  {
	m.logger = log
}

//获取数据库
func (m *Mysql) GetDb() *gorm.DB {
	if err := db.DB().Ping(); err != nil{
		db.Close()
		db = m.NewConnection()
	}
	return db
}

//关闭数据库
func (m *Mysql) Close() error{
	return m.GetDb().Close()
}

func (m *Mysql) GetSql (scope *gorm.Scope) string {
	sql := scope.SQL
	s := reflect.ValueOf(scope.SQLVars)
	for i := 0; i < s.Len(); i++ {
		// 每次取代一个？
		sql = strings.Replace(sql, "?", "'%v'", 1)
		// 赋值
		sql = fmt.Sprintf(sql, s.Index(i))
	}
	return sql
}