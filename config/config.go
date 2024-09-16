package config

import "context"

type Config interface {
	Init(opts ...Option) error

	Load(context.Context, ...LoadOption) error
}

func Load(ctx context.Context, cs ...Config) (err error) {
	for _, c := range cs {
		if err = c.Init(); err != nil {
			break
		}
		if err = c.Load(ctx); err != nil {
			break
		}
	}
	return
}
