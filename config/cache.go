package config

import "time"

type cacheConfig struct {
	Addr     string        `yaml:"addr"     env:"CACHE_ADDR"     env-default:"127.0.0.1"`
	Password string        `yaml:"password"     env:"CACHE_PASSWORD"     env-default:""`
	Db       int           `yaml:"db"     env:"CACHE_DB"     env-default:"0"`
	TTL      time.Duration `yaml:"ttl"     env:"CACHE_TTL"     env-default:"5m"`
}
