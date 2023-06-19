package zap

import (
	"fmt"

	"github.com/ldhk/tonton-be/tonton-be/pkg/telemetry/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	logging.SetDefaultLogger(New(LocalConfig()))
}

type Logger struct {
	*zap.SugaredLogger
}

// WithFields adds a variadic number of fields to the logging context and return new logger.
// Example:
//
//		logger.WithFields(Fields{
//			"hello": "world",
//			"error", errors.New("http timeout"),
//	    	"count", 42,
//		})
func (l *Logger) WithFields(fields map[string]interface{}) logging.Logger {
	var args []interface{}
	for k, v := range fields {
		args = append(args, k)
		args = append(args, v)
	}

	return &Logger{SugaredLogger: l.SugaredLogger.With(args...)}
}

// WithField adds a field to the logging context and return new logger.
// Example:
//
//	logger.WithField("count", 42)
func (l *Logger) WithField(k string, v interface{}) logging.Logger {
	return l.WithFields(map[string]interface{}{k: v})
}

func New(c zap.Config, opts ...zap.Option) *Logger {
	z, err := c.Build(opts...)
	if err != nil {
		panic(fmt.Errorf("logging/zap: build logger from config: %s", err))
	}

	return &Logger{SugaredLogger: z.Sugar()}
}

func ReleaseConfig() zap.Config {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	return cfg
}

func LocalConfig() zap.Config {
	cfg := ReleaseConfig()
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return cfg
}
