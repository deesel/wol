package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var instance *zap.SugaredLogger

func New(level string) (*zap.SugaredLogger, error) {
	if instance == nil {
		config := zap.Config{
			Level:       zap.NewAtomicLevelAt(getLevel(level)),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			DisableCaller: true,
			Encoding:      "json",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "time",
				LevelKey:       "level",
				NameKey:        zapcore.OmitKey,
				CallerKey:      zapcore.OmitKey,
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}

		l, err := config.Build()
		if err != nil {
			return nil, err
		}

		instance = l.Sugar()
	}

	return instance, nil
}

func Get() *zap.SugaredLogger {
	return instance
}

func getLevel(level string) zapcore.Level {
	logLevel, ok := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"fatal": zapcore.FatalLevel,
		"panic": zapcore.PanicLevel,
	}[level]

	if !ok {
		logLevel = zapcore.InfoLevel
	}

	return logLevel
}
