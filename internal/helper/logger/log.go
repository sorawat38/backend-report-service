package logger

import (
	"github.com/CLCM3102-Ice-Cream-Shop/backend-report-service/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger      *zap.Logger
	sugarLogger *zap.SugaredLogger
)

func InitLog(cfg config.Log) {

	var config zap.Config
	if cfg.Env == "prod" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.DisableStacktrace = true
	}

	var err error
	logger, err = config.Build()
	if err != nil {
		panic("Error initializing logger")
	}
	sugarLogger = logger.Sugar()
}

// Use this function to defer closing the logger
func CloseLogger() {
	if err := logger.Sync(); err != nil {
		panic("Error closing logger")
	}
}

func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

func Infof(template string, args ...interface{}) {
	sugarLogger.Infof(template, args...)
}

func Error(msg string, fields ...zapcore.Field) {
	logger.Error(msg, fields...)
}

func Errorf(template string, args ...interface{}) {
	sugarLogger.Errorf(template, args...)
}

func Warn(msg string, fields ...zapcore.Field) {
	logger.Warn(msg, fields...)
}

func Warnf(template string, args ...interface{}) {
	sugarLogger.Warnf(template, args...)
}

func Panic(msg string, fields ...zapcore.Field) {
	logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	logger.Fatal(msg, fields...)
}
