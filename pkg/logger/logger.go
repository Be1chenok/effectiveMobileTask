package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	*zap.SugaredLogger
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	With(args ...interface{}) *zap.SugaredLogger
	Sync() error
}

func New() (Logger, error) {
	if err := createLogsDirectory(); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	conf := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout", "../../logs/app.log"},
		ErrorOutputPaths: []string{"stderr", "../../logs/error.log"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "message",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
	}

	zapLogger, err := conf.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %v", err)
	}

	sugar := zapLogger.Sugar()

	return &logger{
		sugar,
	}, nil
}

func createLogsDirectory() error {
	dir := "../../logs"
	filenames := [2]string{"app.log", "error.log"}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	for _, filename := range filenames {
		fullPath := dir + "/" + filename
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			file, err := os.Create(fullPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %v", fullPath, err)
			}
			file.Close()
		}
	}
	return nil
}
