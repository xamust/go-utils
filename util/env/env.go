package env

import (
	"github.com/labstack/gommon/log"
	"github.com/subosito/gotenv"
	"os"
	"strings"
)

const (
	hostname = "HOSTNAME"
)

func init() {
	gotenv.Load()
}

func Env(key string) string {
	return os.Getenv(key)
}

func GetHostName(serviceName string) string {
	hn := Env(hostname)
	if strings.TrimSpace(hn) == "" {
		log.Warnf("Environment variable %s is not set", hostname)
		return serviceName
	}
	return hn
}
