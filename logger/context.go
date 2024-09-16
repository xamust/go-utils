package logger

import (
	"context"

	"github.com/xamust/go-utils/metadata"

	"golang.org/x/exp/slog"
)

type logKey struct{}

func FromContextLogger(ctx context.Context) Logger {
	if l, ok := ctx.Value(logKey{}).(Logger); !ok {
		return DefaultLogger
	} else {
		return l
	}
}

func NewContextLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, logKey{}, logger)
}

type event struct {
	source   string
	receiver string
}

func getEvent(ctx context.Context) []slog.Attr {
	e, ok := ctx.Value(event{}).(event)
	if !ok {
		meta, _ := metadata.FromContextHeader(ctx)
		e.source = meta.Header.SourceSystem
		e.receiver = meta.Header.ReceiverSystem
	}

	attr := make([]slog.Attr, 0, 2)

	attr = append(attr, slog.String(keyEventRec, e.receiver), slog.String(keyEventSource, e.source))

	return attr
}

func NewContextEvent(ctx context.Context, source, receiver string) context.Context {
	return context.WithValue(ctx, event{}, event{
		source:   source,
		receiver: receiver,
	})
}
