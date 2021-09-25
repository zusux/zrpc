package zlog

import (
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"sync"
	"time"
)
var once sync.Once
var loger *logrus.Logger

type Format string
var (
	Json Format ="json"
	Text Format ="text"
)

type Log struct{
	Path string
	File string
	Age int64
	Rotation int64
	Format
}

func NewLog(path ,file ,format string,age,rotation int64) *Log{
	l := &Log{
		Path: path,
		File: file,
		Format:Format(format),
	}
	return l
}

func (l *Log) Zlog() *logrus.Logger{
	once.Do(func() {
		loger = logrus.New()
		l.initZlog()
		baseLogPaht := path.Join(l.Path, l.File)
		writer, err := rotatelogs.New(
			baseLogPaht+".%Y-%m-%d[%H]",
			rotatelogs.WithLinkName(baseLogPaht), // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(time.Duration(l.Age) * time.Hour), // 文件最大保存时间
			rotatelogs.WithRotationTime(time.Duration(l.Rotation) * time.Hour), // 日志切割时间间隔
		)
		if err != nil {
			logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
		}
		if l.Format == Text{
			lfHook := lfshook.NewHook(lfshook.WriterMap{
				logrus.DebugLevel: writer,
				logrus.InfoLevel:  writer,
				logrus.WarnLevel:  writer,
				logrus.ErrorLevel: writer,
				logrus.FatalLevel: writer,
				logrus.PanicLevel: writer,
			},&logrus.TextFormatter{
				TimestampFormat:"2006-01-02 15:03:04",
			})
			loger.AddHook(lfHook)
		}else{
			lfHook := lfshook.NewHook(lfshook.WriterMap{
				logrus.DebugLevel: writer,
				logrus.InfoLevel:  writer,
				logrus.WarnLevel:  writer,
				logrus.ErrorLevel: writer,
				logrus.FatalLevel: writer,
				logrus.PanicLevel: writer,
			},&logrus.JSONFormatter{
				TimestampFormat:"2006-01-02 15:03:04",
			})
			loger.AddHook(lfHook)
		}

	})
	return loger
}



func (l *Log) SetAge(age int64) *Log {
	l.Age = age
	return l
}
func (l *Log)SetRotation(rotation int64) *Log {
	l.Rotation = rotation
	return l
}

func (l *Log)SetPath(path string) *Log {
	l.Path = path
	return l
}

func (l *Log)SetFile(filename string) *Log {
	l.File = filename
	return l
}

func (l *Log)SetFormat(format string) *Log {
	l.Format = Format(format)
	return l
}

//设置默认值
func (l *Log) initZlog()  {
	if l.Path == ""{
		l.SetPath("logs")
	}
	if l.File == ""{
		l.SetFile("log")
	}
	err := os.MkdirAll(l.Path,0777)
	if err != nil{
		log.Fatalln(err)
	}
	if l.Age <= 0{
		l.SetAge(24)
	}
	if l.Rotation <= 0{
		l.SetRotation(24)
	}
	if l.Format == ""{
		l.SetFormat("text")
	}
}