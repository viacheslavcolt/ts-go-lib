package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _default *Logger = nil

func InitLogger() error {
	log, err := MakeConsoleLogger(nil)

	if err != nil {
		return err
	}

	_default = log

	return nil
}

type Logger struct {
	_logger *zap.Logger
}

// API
func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	l._logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zapcore.Field) {
	l._logger.Warn(msg, fields...)
}

func (l *Logger) Error(title string, err string, fields ...zapcore.Field) {
	if len(err) > 0 {
		fields = append(fields, zap.String("err", err))
	}

	l._logger.Error(title, fields...)
}

func (l *Logger) Fatal(title string, err string, fields ...zapcore.Field) {
	if len(err) > 0 {
		fields = append(fields, zap.String("err", err))
	}

	l._logger.Fatal(title, fields...)
}

func Info(msg string, fields ...zapcore.Field) {
	if _default == nil {
		return
	}

	_default.Info(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	if _default == nil {
		return
	}

	_default.Warn(msg, fields...)
}

func Error(title string, err string, fields ...zapcore.Field) {
	if _default == nil {
		return
	}

	_default.Error(title, err, fields...)
}

func Fatal(title string, err string, fields ...zapcore.Field) {
	if _default == nil {
		return
	}

	_default.Fatal(title, err, fields...)
}
