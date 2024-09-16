package env

import (
	"context"
	"os"
	"testing"

	cfg "github.com/xamust/go-utils/config"
	"github.com/xamust/go-utils/config/file"

	"github.com/stretchr/testify/assert"
)

func Test_EnvParse(t *testing.T) {
	configServerFirst := &struct {
		Server string `env:"SERVER"`
	}{}

	err := os.Setenv("SERVER", ":5201")
	assert.Nil(t, err)

	err = cfg.Load(
		context.Background(),
		NewConfig(
			cfg.Dest(configServerFirst),
		),
	)
	assert.Nil(t, err)
}

func Test_Merge(t *testing.T) {
	urls := []string{"http://localhost0.com", "http://localhost1.com", "http://localhost2.com"}
	config := &struct {
		Addrs []string `env:"ADDRS" json:"addrs"`
	}{
		Addrs: []string{urls[0]},
	}

	err := os.Setenv("ADDRS", urls[1])
	assert.Nil(t, err)
	err = cfg.Load(
		context.Background(),
		NewConfig(
			cfg.Dest(config),
		),
		file.NewConfig(
			cfg.Dest(config),
			file.Path("config.json"),
		),
	)
	assert.Nil(t, err)
	assert.Equal(t, urls, config.Addrs)
}
