package logger

import (
	"context"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//const (
//	Development = "dev"
//	Production  = "prod"
//	Test        = "test"
//)

var logger *zap.Logger
var sugar *zap.SugaredLogger

func Init() {
	config := zap.Config{
        Encoding: "json",
        Level: zap.NewAtomicLevelAt(zap.InfoLevel),
        OutputPaths: []string{"stdout"},
        EncoderConfig: zapcore.EncoderConfig{
            TimeKey:        "timestamp",
            LevelKey:      "level",
            MessageKey:    "message",
            EncodeLevel:   zapcore.LowercaseLevelEncoder,
            EncodeTime:    zapcore.TimeEncoderOfLayout(time.RFC3339),
            EncodeCaller:  zapcore.ShortCallerEncoder,
        },
    }
    
    // Добавляем базовые поля, которые будут в каждом логе
    defaultFields := []zap.Field{
		zap.String("service", "fitness-trainer"),
		zap.String("env", os.Getenv("ENV")),
	}
    
    var err error
    logger, err = config.Build(zap.Fields(
        defaultFields...,
    ))
    if err != nil {
        panic(err)
    }
    sugar = logger.Sugar()
}

func Logger() *zap.Logger {
	return logger
}

func formatFields(ctx context.Context) map[string]interface{} {
	fields := make(map[string]interface{})

	span := opentracing.SpanFromContext(ctx)

	if spanCtx, ok := span.Context().(jaeger.SpanContext); ok {
		fields["trace_id"] = spanCtx.TraceID().String()
	}

	return fields
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
