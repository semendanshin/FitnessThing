package logger

import (
	"os"

	"go.uber.org/zap"
)

//const (
//	Development = "dev"
//	Production  = "prod"
//	Test        = "test"
//)

var logger *zap.Logger
var sugar *zap.SugaredLogger

func Init() {
	logger, _ = zap.NewDevelopment()
	sugar = logger.Sugar()
}

func Logger() *zap.Logger {
	return logger
}

func Debug(msg string, args ...interface{}) {
	if logger == nil {
		panic("logger not initialized")
	}

	sugar.Debugw(msg, args...)
}

func Debugf(msg string, args ...interface{}) {
	if logger == nil {
		panic("logger not initialized")
	}

	sugar.Debugf(msg, args...)
}

func Info(msg string, args ...interface{}) {
	if logger == nil {
		panic("logger not initialized")
	}

	sugar.Infow(msg, args...)
}

func Infof(msg string, args ...interface{}) {
	if logger == nil {
		panic("logger not initialized")
	}

	sugar.Infof(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	if logger == nil {
		panic("logger not initialized")
	}

	sugar.Warnw(msg, args...)
}

func Warnf(msg string, args ...interface{}) {
	if logger == nil {
		panic("logger not initialized")
	}

	sugar.Warnf(msg, args...)
}

func Error(msg string, args ...interface{}) {
	if logger == nil {
		panic("logger not initialized")
	}

	sugar.Errorw(msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	if logger == nil {
		panic("logger not initialized")
	}

	sugar.Errorf(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	if logger == nil {
		panic("logger not initialized")
	}

	sugar.Errorw(msg, args...)
	os.Exit(1)
}

func Fatalf(msg string, args ...interface{}) {
	if logger == nil {
		panic("logger not initialized")
	}

	sugar.Errorf(msg, args...)
	os.Exit(1)
}

func Panic(msg string, args ...interface{}) {
	if logger == nil {
		panic("logger not initialized")
	}

	sugar.Errorw(msg, args...)
	panic(msg)
}

func Panicf(msg string, args ...interface{}) {
	if logger == nil {
		panic("logger not initialized")
	}

	sugar.Errorf(msg, args...)
	panic(msg)
}
