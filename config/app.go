package config

import "time"

type appConfig struct {
	JWTSecret string        `yaml:"jwt_secret"     env:"JWT_SECRET"     env-default:"verysecretkey"`
	JWTTTL    time.Duration `yaml:"jwt_ttl"     env:"JWT_TTL"     env-default:"45m"`
}
