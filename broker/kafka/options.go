package kafka

import (
	"github.com/xamust/go-utils/encoder"
	"github.com/xamust/go-utils/encoder/json"
	"github.com/xamust/go-utils/logger"
)

type PublishOptions struct {
	codec  encoder.Codec
	logger logger.Logger
}

type PublishOption func(o *PublishOptions)

func NewOptions(opts ...PublishOption) PublishOptions {
	options := PublishOptions{
		codec:  json.NewCodec(),
		logger: logger.DefaultLogger,
	}
	for _, o := range opts {
		o(&options)
	}

	return options
}

func Codec(c encoder.Codec) PublishOption {
	return func(o *PublishOptions) {
		o.codec = c
	}
}

func Logger(l logger.Logger) PublishOption {
	return func(o *PublishOptions) {
		o.logger = l
	}
}
