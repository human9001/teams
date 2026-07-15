package config

import (
	"net"
	"time"
)

type httpConfig struct {
	Host    string            `yaml:"host" env:"HTTP_HOST" env-default:"localhost"`
	Port    string            `yaml:"port" env:"HTTP_PORT" env-default:"8090"`
	Options httpOptionsConfig `yaml:"options"`
}

type httpOptionsConfig struct {
	HttpReadHeaderTimeout time.Duration `yaml:"httpReadHeaderTimeout" env:"HTTP_READ_HEADER_TIMEOUT" env-default:"5s"`
	HttpReadTimeout       time.Duration `yaml:"httpReadTimeout" env:"HTTP_READ_TIMEOUT" env-default:"15s"`
	HttpWriteTimeout      time.Duration `yaml:"httpWriteTimeout" env:"HTTP_WRITE_TIMEOUT" env-default:"15s"`
	HttpIdleTimeout       time.Duration `yaml:"httpIdleTimeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
	HttpShutdownTimeout   time.Duration `yaml:"httpShutdownTimeout" env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"5s"`
}

func (c *httpConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
