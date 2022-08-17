package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
