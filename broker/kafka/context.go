package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type headerKey struct{}

func NewContextMetadata(ctx context.Context, headers []kafka.Header) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, headerKey{}, headers)
}

func FromContextMetadata(ctx context.Context) ([]kafka.Header, bool) {
	if m, ok := ctx.Value(headerKey{}).([]kafka.Header); !ok {
		return nil, ok
	} else {
		return m, ok
	}
}
