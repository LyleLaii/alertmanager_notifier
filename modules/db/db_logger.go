package db

import (
	"alertmanager_notifier/config"
	"alertmanager_notifier/log"
	"context"
	"fmt"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

var (
	infoStr      = "%s "
	warnStr      = "%s "
	errStr       = "%s "
	traceStr     = "%s [%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s [%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s [%.3fms] [rows:%v] %s"
)

type DBLogger struct {
	logger        log.Logger
	loglevel      config.Level
	SlowThreshold time.Duration
}


func NewDBLogger(logger log.Logger, runConf *config.RunningConfig) gormlogger.Interface {
	return &DBLogger{logger: logger, loglevel: runConf.Level.Level(), SlowThreshold: time.Duration(100) * time.Millisecond}
}

func (l *DBLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	//newlogger.loglevel = level
	return &newlogger
}

func (l DBLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.loglevel <= config.Info {
		l.logger.Info("gorm", fmt.Sprintf(infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)))
	}
}

func (l DBLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.loglevel <= config.Warn {
		l.logger.Warn("gorm", fmt.Sprintf(warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)))
	}
}

func (l DBLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.loglevel <= config.Error {
		l.logger.Error("gorm", fmt.Sprintf(errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)))
	}
}

func (l DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
		elapsed := time.Since(begin)

		if err != nil {
			sql, rows := fc()
			if rows == -1 {
				l.logger.Error("gorm", fmt.Sprintf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql))
			} else {
				l.logger.Error("gorm", fmt.Sprintf(traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql))
			}
		} else if elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.loglevel >= config.Warn {
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				l.logger.Warn("gorm", fmt.Sprintf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql))
			} else {
				l.logger.Warn("gorm", fmt.Sprintf(traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql))
			}
		} else {
			sql, rows := fc()
			if rows == -1 {
				l.logger.Debug("gorm", fmt.Sprintf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql))
			} else {
				l.logger.Debug("gorm", fmt.Sprintf(traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql))
			}
		}


}
