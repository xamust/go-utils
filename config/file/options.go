package file

import cfg "github.com/xamust/go-utils/config"

type pathKey struct{}

func Path(p string) cfg.Option {
	return cfg.SetOption(pathKey{}, p)
}

func LoadPath(path string) cfg.LoadOption {
	return cfg.SetLoadOption(pathKey{}, path)
}

func PathFromOption(opts cfg.Options) string {
	return opts.Context.Value(pathKey{}).(string)
}

func PathFromLoadOption(opts cfg.LoadOptions) string {
	return opts.Context.Value(pathKey{}).(string)
}
