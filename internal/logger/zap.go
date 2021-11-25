package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logLevels = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
}

func NewZapLogger(level string) (*zap.Logger, error) {
	level = strings.ToLower(level)

	logLevel, ok := logLevels[level]
	if !ok {
		return nil, fmt.Errorf(
			"invalid log level: %s. should be one of debug|info|warn|error. case insensetive",
			level,
		)
	}

	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    newZapEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	return zapConfig.Build()
}

func newZapEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "channel",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
