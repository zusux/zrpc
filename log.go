package zrpc

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"time"
)

type zlog struct {
	filepath string
	fileName string
	age      int64
	rotation int64
	loger  *logrus.Logger
}

func NewLog( logSection string) *zlog  {
	filepath := K.String(fmt.Sprintf("%s.path",logSection))
	file := K.String(fmt.Sprintf("%s.file",logSection))
	age := K.Int64(fmt.Sprintf("%s.age",logSection))
	rotation := K.Int64(fmt.Sprintf("%s.rotation",logSection))

	zlog := &zlog{
		filepath: filepath,
		fileName: file,
		age:      age,
		rotation: rotation,
	}
	zlog.init()
	return zlog
}

func (l *zlog) init() *logrus.Logger {
	l.getLogSet()
	loger := logrus.New()
	baseLogPaht := path.Join(l.filepath, l.fileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y-%m-%d",
		rotatelogs.WithLinkName(baseLogPaht),                             // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(l.age)*time.Hour),            // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Duration(l.rotation)*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	loger.AddHook(lfHook)
	l.loger = loger
	return loger
}

//设置默认值
func (l *zlog) getLogSet() {
	if l.filepath == "" {
		l.filepath = "logs"
	}
	if l.fileName == "" {
		l.fileName = "log"
	}
	if l.age <= 0 {
		l.age = 24
	}
	if l.rotation <= 0 {
		l.rotation = 24
	}
	err := os.MkdirAll(l.filepath, 0777)
	if err != nil {
		log.Fatalln(err)
	}
}
