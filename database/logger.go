package database

import (
	"context"
	log "github.com/fhmq/hmq/logger"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

// Logger 日志打印处理
type Logger struct {
}

func NewGormLogger() logger.Interface {
	return Logger{}
}

func (l Logger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Infof(msg, data...)
}

func (l Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Warnf(msg, data...)
}

func (l Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Errorf(msg, data...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil:
		sql, rows := fc()
		if rows == -1 {
			log.Errorf("%s %s [%.3fms] [rows:%v] %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Errorf("%s %s [%.3fms] [rows:%v] %s", utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		sql, rows := fc()
		if rows == -1 {
			log.Debugf("%s [%.3fms] [rows:%v] %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Debugf("%s [%.3fms] [rows:%v] %s", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
