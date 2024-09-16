package file

import (
	"context"
	"fmt"
	"os"
	"reflect"

	cfg "github.com/xamust/go-utils/config"
	ureflect "github.com/xamust/go-utils/util/reflect"

	"dario.cat/mergo"
)

type config struct {
	opts cfg.Options
	path string
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

	c.path = PathFromOption(c.opts)

	if len(c.path) < 1 {
		err = cfg.ErrPathNotFound
		c.opts.Logger.Error(context.Background(), err.Error())
	}

	return err
}

func (c *config) Load(ctx context.Context, opts ...cfg.LoadOption) error {
	path := c.path
	loadOpts := cfg.NewLoadOptions(opts...)
	if loadOpts.Context != nil {
		if v, ok := loadOpts.Context.Value(pathKey{}).(string); ok && v != "" {
			path = v
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

	fDescription, err := os.OpenFile(path, os.O_RDONLY, os.FileMode(0400))
	if err != nil {
		c.opts.Logger.Error(c.opts.Context, fmt.Sprintf("file load path %s error: %v", path, err))
		return err
	}
	defer func() {
		if err := fDescription.Close(); err != nil {
			c.opts.Logger.Error(c.opts.Context, fmt.Sprintf("close discriptor with error: %v", err))
		}
	}()

	if err := c.opts.Codec.ReadBody(fDescription, srcZero); err != nil {
		c.opts.Logger.Error(c.opts.Context, fmt.Sprintf("file load path %s error: %v", path, err))
		return err
	}

	mOpts := make([]func(*mergo.Config), 0, 1)
	mOpts = append(mOpts, mergo.WithTypeCheck)
	if loadOpts.Append {
		mOpts = append(mOpts, mergo.WithAppendSlice)
	}
	if loadOpts.Override {
		mOpts = append(mOpts, mergo.WithOverride)
	}

	err = mergo.Merge(c.opts.Dest, srcZero, mOpts...)
	return err
}
