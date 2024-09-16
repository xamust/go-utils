package logger

import (
	"context"
)

var (
	// DefaultLogger variable
	DefaultLogger Logger = NewLogger(nil)
	// DefaultLevel used by logger
	DefaultLevel Level = InfoLevel
)

type Logger interface {
	Init(opts ...Option) error

	Log(ctx context.Context, level Level, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Fatal(ctx context.Context, msg string, args ...any)

	String() string
	Fields(fields map[string]interface{}) Logger
}

func NewLogger(cfg *Config, o ...Option) Logger {
	level := InfoLevel
	msgSize := defaultMaxBytesMessage
	if cfg != nil {
		if lvl, err := GetLevel(cfg.Level); err == nil {
			level = lvl
		}
		msgSize = cfg.MaxMsgSize
		if msgSize < (1 << 7) {
			msgSize = defaultMaxBytesMessage
		}
	}

	o = append([]Option{WithLevel(level), MaxBytesMessage(msgSize)}, o...)
	l := NewSlogLogger(o...)

	if err := l.Init(); err != nil {
		panic("problem with init logger")
	}

	return l
}
