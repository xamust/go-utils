package env

import (
	"context"
	"reflect"

	cfg "github.com/xamust/go-utils/config"
	ureflect "github.com/xamust/go-utils/util/reflect"

	"dario.cat/mergo"
	"github.com/caarlos0/env/v10"
)

type config struct {
	opts cfg.Options
}

func NewConfig(opts ...cfg.Option) cfg.Config {
	options := cfg.NewOptions(opts...)
	return &config{
		opts: options,
	}
}

func (c *config) Init(opts ...cfg.Option) (err error) {
	for _, o := range opts {
		o(&c.opts)
	}

	return err
}

func (c *config) Load(ctx context.Context, opts ...cfg.LoadOption) error {
	for _, f := range c.opts.Funcs {
		if err := f(ctx, c); err != nil {
			return err
		}
	}

	va := reflect.ValueOf(c.opts.Dest)
	if va.Kind() != reflect.Pointer {
		return cfg.ErrNotPtr
	}

	srcZero, err := ureflect.NewZeroValue(c.opts.Dest)
	if err != nil {
		return err
	}

	if err = env.Parse(srcZero); err != nil {
		return err
	}

	loadOpts := cfg.NewLoadOptions(opts...)
	mOpts := make([]func(*mergo.Config), 0, 1)
	mOpts = append(mOpts, mergo.WithTypeCheck)
	if loadOpts.Append {
		mOpts = append(mOpts, mergo.WithAppendSlice)
	}
	if loadOpts.Override {
		mOpts = append(mOpts, mergo.WithOverride)
	}

	if err := mergo.Merge(c.opts.Dest, srcZero, mOpts...); err != nil {
		return err
	}

	return nil
}
