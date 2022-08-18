package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	_SamplingInitial     = 100
	_SamplingTheareafter = 100

	_JsonEncoderType    = "json"
	_ConsoleEncoderType = "console"

	_InfoLevel  = zap.InfoLevel
	_DebugLevel = zap.DebugLevel
)

var DefaultOut = []string{"stderr"}

func _getEncoderConfig(isDev bool) zapcore.EncoderConfig {
	var (
		cfg zapcore.EncoderConfig
	)

	if isDev {
		cfg = zap.NewDevelopmentEncoderConfig()
	} else {
		cfg = zap.NewProductionEncoderConfig()
	}

	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.CallerKey = zapcore.OmitKey
	cfg.StacktraceKey = zapcore.OmitKey

	return cfg
}

func _makeLoggerConfig(isDev bool, level zapcore.Level, encoderType string, outPaths []string) zap.Config {
	if outPaths == nil {
		outPaths = DefaultOut
	}

	return zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: isDev,
		Sampling: &zap.SamplingConfig{
			Initial:    _SamplingInitial,
			Thereafter: _SamplingTheareafter,
		},
		Encoding:         encoderType,
		EncoderConfig:    _getEncoderConfig(isDev),
		OutputPaths:      outPaths,
		ErrorOutputPaths: outPaths,
	}
}

func _makeLogger(isDev bool, level zapcore.Level, encoderType string, outPaths []string) (*Logger, error) {
	var (
		_logger *zap.Logger
		err     error
	)

	if _logger, err = _makeLoggerConfig(isDev, level, encoderType, outPaths).Build(); err != nil {
		return nil, err
	}

	return &Logger{
		_logger: _logger,
	}, nil
}

// api of this module
func MakeJsonLogger(outPaths []string) (*Logger, error) {
	return _makeLogger(false, _InfoLevel, _JsonEncoderType, outPaths)
}

func MakeDevLogger(outPaths []string) (*Logger, error) {
	return _makeLogger(true, _DebugLevel, _ConsoleEncoderType, outPaths)
}

func MakeConsoleLogger(outPaths []string) (*Logger, error) {
	return _makeLogger(false, _InfoLevel, _ConsoleEncoderType, outPaths)
}
