package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
)

var (
	logger      *zap.Logger
	sugarLogger *zap.SugaredLogger
)

// Config 配置文件
type Config struct {
	Path  string // 日志位置，如果为空，则输出到屏幕
	Debug bool   // 调试模式
}

// init 默认初始化
func init() {
	logger = NewLogger()
	sugarLogger = logger.Sugar()
}

// ConfigLogger 配置日志
func ConfigLogger(config Config) {
	logger = NewLogger(config)
	sugarLogger = logger.Sugar()
}

// NewLogger 新建ZAP日志
func NewLogger(configs ...Config) *zap.Logger {

	var config Config
	if len(configs) == 0 {
		config = Config{}
	} else {
		config = configs[0]
	}

	// 调整时间格式是格式化的
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 设置输出普通的Log Encoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	var writeSyncers []zapcore.WriteSyncer

	// 设置屏幕输出
	writeSyncers = append(writeSyncers, os.Stdout)

	// 设置文件输出
	if config.Path != "" {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   path.Join(config.Path, "hiot.log"),
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     90,
			Compress:   false,
		}
		writeSyncers = append(writeSyncers, zapcore.AddSync(lumberJackLogger))
	}

	writeSyncer := zapcore.NewMultiWriteSyncer(writeSyncers...)

	level := zapcore.InfoLevel
	if config.Debug {
		level = zapcore.DebugLevel
	}
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// zap.AddCaller 输出哪里报错
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Debugf(template string, args ...interface{}) {
	sugarLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	sugarLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sugarLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	sugarLogger.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	sugarLogger.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	sugarLogger.Fatalf(template, args...)
}
