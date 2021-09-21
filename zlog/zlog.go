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
	JSON Format ="JSON"
	Text Format ="Text"
)

type Log struct{
	Path string
	File string
	Age time.Duration
	Rotation time.Duration
	Format
}

func NewLog(path string,format Format) *Log{
	l := &Log{
		Path: path,
		Format:format,
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
			rotatelogs.WithMaxAge(l.Age), // 文件最大保存时间
			rotatelogs.WithRotationTime(l.Rotation), // 日志切割时间间隔
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
			},&logrus.TextFormatter{})
			loger.AddHook(lfHook)
		}else{
			lfHook := lfshook.NewHook(lfshook.WriterMap{
				logrus.DebugLevel: writer,
				logrus.InfoLevel:  writer,
				logrus.WarnLevel:  writer,
				logrus.ErrorLevel: writer,
				logrus.FatalLevel: writer,
				logrus.PanicLevel: writer,
			},&logrus.JSONFormatter{})
			loger.AddHook(lfHook)
		}

	})
	return loger
}



func (l *Log) SetAge(age time.Duration) *Log {
	l.Age = time.Hour * time.Duration(age)
	return l
}
func (l *Log)SetRotation(rotation time.Duration) *Log {
	l.Rotation = time.Hour * time.Duration(rotation)
	return l
}

func (l *Log)SetFile(filename string) *Log {
	l.File = filename
	return l
}

func (l *Log)SetFormat(format Format) *Log {
	l.Format = format
	return l
}

//设置默认值
func (l *Log) initZlog()  {
	if l.Path == ""{
		l.Path = "logs"
	}
	if l.File == ""{
		l.SetFile("log")
	}
	err := os.MkdirAll(l.Path,0777)
	if err != nil{
		log.Fatalln(err)
	}
	if l.Age <= 0{
		l.SetAge(time.Hour*24)
	}
	if l.Rotation <= 0{
		l.SetRotation(time.Hour*24)
	}
	if l.Format == ""{
		l.SetFormat(Text)
	}
}