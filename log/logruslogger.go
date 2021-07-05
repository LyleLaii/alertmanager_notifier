package log

import (
	"fmt"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	*logrus.Logger
}

func New(name string) Logger {
	logger := logrus.New()
	currentPath, _ := os.Getwd()
	logFilePath := "logs"
	logFileName := fmt.Sprintf("%s.log", name)

	// 日志文件
	fileName := path.Join(currentPath, logFilePath, logFileName)

	// 写入文件
	// src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	// src, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	fmt.Println("err", err)
	// }

	// 实例化
	// logger := logrus.New()

	// 设置输出
	// mw := io.MultiWriter(os.Stdout, src)
	// Logger.SetOutput(mw)
	// Logger.Out = src

	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	// 设置 rotatelogs
	logWriter, _ := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d",
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 新增 Hook
	logger.AddHook(lfHook)

	return &LogrusLogger{logger}
}

// Debug log warn infomation
func (l *LogrusLogger) Debug(module string, data interface{}) {
	l.WithFields(
		logrus.Fields{
			"module": module,
		}).Debug(data)
}

// Info log info infomation
func (l *LogrusLogger) Info(module string, data interface{}) {
	l.WithFields(
		logrus.Fields{
			"module": module,
		}).Info(data)
}

// Warn log warn infomation
func (l *LogrusLogger) Warn(module string, data interface{}) {
	l.WithFields(
		logrus.Fields{
			"module": module,
		}).Warn(data)
}

// Error log error ingomation
func (l *LogrusLogger) Error(module string, data interface{}) {
	l.WithFields(
		logrus.Fields{
			"module": module,
		}).Error(data)
}

// Panic log panic information
func (l *LogrusLogger) Panic(module string, data interface{}) {
	l.WithFields(
		logrus.Fields{
			"module": module,
		}).Panic(data)
}
