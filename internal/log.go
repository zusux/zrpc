package internal

import (
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/zusux/zrpc/internal/util"
	"path"
	"time"
)

type ZLog struct {
	Path     string
	File     string
	Age      int64
	Rotation int64
	Log   *logrus.Logger
	Initialize bool
}

func NewLog(path, file string, age, rotation int64) *ZLog {
	log := &ZLog{
		Path:     path,
		File:     file,
		Age:      age,
		Rotation: rotation,
	}
	log.init()
	return log
}

func (l *ZLog) init() {
	l.Log = logrus.New()
	baseLogPath := l.GetFilePath()
	writer, err := rotatelogs.New(
		baseLogPath+".%Y-%m-%d",
		rotatelogs.WithLinkName(baseLogPath),                                  // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(l.getAge())*time.Hour),            // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Duration(l.GetRotation())*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("[zrpc] rotatelogs error! config local file system logger %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	l.Initialize = true
	l.Log.AddHook(lfHook)
}

func (l *ZLog) GetFilePath() string {
	if l.Path == "" {
		l.Path = "logs"
	}
	if l.File == "" {
		l.File = "log"
	}
	baseLogPath := path.Join(l.Path, l.File)
	if !l.CheckOrCreateDIr(baseLogPath,"config"){
		workingDir := util.GetWdDir()
		if !l.CheckOrCreateDIr(workingDir,"working"){
			panic(errors.New("[zrpc][log] log path is not available"))
		}
		l.Path = workingDir
		baseLogPath = path.Join(l.Path, l.File)
		l.Log.Warnf("[zrpc][log] switch log path to working dir: %s",workingDir)
	}
	return baseLogPath
}

func (l *ZLog) CheckOrCreateDIr(baseLogPaht string, where string) bool {
	ok,err := util.AvailablePath(baseLogPaht)
	if err == nil && ok{
		return true
	}
	l.Log.Warnf("[zrpc][log] %s path is not available, exsit: %v, err: %v", where,ok,err)
	return false
}

//设置默认值
func (l *ZLog) getAge() int64 {
	if l.Age <= 0 {
		l.Age = 24
	}
	return l.Age
}
func (l *ZLog) GetRotation() int64 {
	if l.Rotation <= 0 {
		l.Rotation = 24
	}
	return l.Rotation
}
