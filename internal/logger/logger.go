package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents logger instance
type Logger struct {
	engine *zap.SugaredLogger
	config *zap.Config
}

var instance *Logger

// New returns new logger instance
func New() *Logger {
	if instance == nil {
		config := zap.Config{
			Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
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
			panic(err)
		}

		instance = &Logger{
			engine: l.Sugar(),
			config: &config,
		}
	}

	return instance
}

// Logger returns logger instance engine
func (l *Logger) Logger() *zap.SugaredLogger {
	return New().engine
}

// SetLevel sets logging level. Available options are "debug", "info", "warn", "error", "fatal", "panic". If different value is passed as argument, info level is used by default. Returns logger instance engine.
func (l *Logger) SetLevel(level string) *zap.SugaredLogger {
	logLevel := getLevel(level)
	instance := New()
	instance.config.Level.SetLevel(logLevel)
	return instance.engine
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
