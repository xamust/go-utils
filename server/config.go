package server

import (
	mb "github.com/xamust/go-utils/models_bcon"
)

type Config struct {
	Name        string      `json:"name" env:"SERVER_NAME" yaml:"name" mapstructure:"name"`
	Version     string      `json:"version" env:"SERVER_VERSION" yaml:"version" mapstructure:"version"`
	Addr        string      `json:"addr" env:"SERVER_ADDRESS" yaml:"addr" mapstructure:"addr"`
	TimeoutConn mb.Duration `json:"timeout_conn" env:"SERVER_TIMEOUT_CONN" yaml:"timeout_conn" mapstructure:"timeout_conn"`
}

type EndpointParams struct {
	Url      string      `json:"url" yaml:"url" mapstructure:"url"`
	Timeout  mb.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Proxy    string      `json:"proxy" yaml:"proxy" mapstructure:"proxy"`
	Insecure bool        `json:"insecure" yaml:"insecure" mapstructure:"insecure"`
	Login    string      `json:"login" yaml:"login" mapstructure:"login"`
	Password string      `json:"password" yaml:"password" mapstructure:"password"`
}
