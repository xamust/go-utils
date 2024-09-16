package config

import (
	"context"

	"github.com/xamust/go-utils/encoder"
	"github.com/xamust/go-utils/encoder/json"
	"github.com/xamust/go-utils/logger"
)

type Options struct {
	Dest    any
	Context context.Context
	Codec   encoder.Codec
	Logger  logger.Logger

	Funcs []func(ctx context.Context, config Config) error
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Logger:  logger.DefaultLogger,
		Codec:   json.NewCodec(),
		Context: context.Background(),
		Funcs:   make([]func(ctx context.Context, config Config) error, 0),
	}
	for _, o := range opts {
		o(&options)
	}

	return options
}

func Codec(c encoder.Codec) Option {
	return func(o *Options) {
		o.Codec = c
	}
}

func Dest(dst any) Option {
	return func(o *Options) {
		o.Dest = dst
	}
}

func Logger(l logger.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

type LoadOptions struct {
	dest    any
	Context context.Context

	Override bool
	Append   bool
}

type LoadOption func(o *LoadOptions)

func NewLoadOptions(opts ...LoadOption) LoadOptions {
	options := LoadOptions{
		Append: true,
	}

	for _, o := range opts {
		o(&options)
	}
	return options
}

func LoadDest(dest interface{}) LoadOption {
	return func(o *LoadOptions) {
		o.dest = dest
	}
}

func LoadOverride(f bool) LoadOption {
	return func(o *LoadOptions) {
		o.Override = f
	}
}

func LoadAppend(f bool) LoadOption {
	return func(o *LoadOptions) {
		o.Append = f
	}
}

func AddFunc(fs ...func(ctx context.Context, config Config) error) Option {
	return func(o *Options) {
		o.Funcs = append(o.Funcs, fs...)
	}
}
