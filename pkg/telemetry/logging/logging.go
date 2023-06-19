package logging

import "context"

var (
	l Logger = NullLogger{}
)

func SetDefaultLogger(logger Logger) {
	if logger == nil {
		logger = NullLogger{}
	}
	l = logger
}

type loggerKey struct{}

// Logger define the interface for logging in this service
type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})

	Warn(args ...interface{})
	Warnf(template string, args ...interface{})

	Error(args ...interface{})
	Errorf(template string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})

	WithFields(fields map[string]interface{}) Logger
	WithField(k string, v interface{}) Logger
}

// IntoContext return a new context with the logger injected
func IntoContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerKey{}).(Logger); ok {
		return l
	}

	return l
}

func WithFields(ctx context.Context, fields map[string]interface{}) (context.Context, Logger) {
	logger := FromContext(ctx).WithFields(fields)
	return IntoContext(ctx, logger), logger
}

func WithField(ctx context.Context, k string, v interface{}) (context.Context, Logger) {
	logger := FromContext(ctx).WithField(k, v)
	return IntoContext(ctx, logger), logger
}

func Copy(dst context.Context, src context.Context) context.Context {
	return IntoContext(dst, FromContext(src))
}

type NullLogger struct{}

func (n NullLogger) Info(_ ...interface{}) {}

func (n NullLogger) Infof(_ string, _ ...interface{}) {}

func (n NullLogger) Warn(_ ...interface{}) {}

func (n NullLogger) Warnf(_ string, _ ...interface{}) {}

func (n NullLogger) Error(_ ...interface{}) {}

func (n NullLogger) Errorf(_ string, _ ...interface{}) {}

func (n NullLogger) Fatal(_ ...interface{}) {}

func (n NullLogger) Fatalf(_ string, _ ...interface{}) {}

func (n NullLogger) WithField(_ string, _ interface{}) Logger { return n }

func (n NullLogger) WithFields(_ map[string]interface{}) Logger { return n }
