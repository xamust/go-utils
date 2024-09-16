package file

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	cfg "github.com/xamust/go-utils/config"
	"github.com/xamust/go-utils/server"

	"github.com/stretchr/testify/assert"
)

func Test_FileCfg(t *testing.T) {
	pathFirst := "./config.json"
	configServerFirst := &struct {
		Server server.Config `json:"server_first"`
	}{}
	configServerSecond := &struct {
		Server server.Config `json:"server_second"`
	}{}

	err := cfg.Load(
		context.Background(),
		NewConfig(
			cfg.Dest(configServerFirst),
			Path(pathFirst),
		),
		NewConfig(
			cfg.Dest(configServerSecond),
			Path(pathFirst),
		),
	)
	assert.Nil(t, err)

	file, err := os.ReadFile(pathFirst)
	assert.Nil(t, err)
	assert.NotNil(t, file)

	configServer := &struct {
		ServerF server.Config `json:"server_first"`
		ServerS server.Config `json:"server_second"`
	}{}
	err = json.Unmarshal(file, configServer)
	assert.Nil(t, err)

	assert.Equal(t, configServer.ServerF, configServerFirst.Server)
	assert.Equal(t, configServer.ServerS, configServerSecond.Server)
}
