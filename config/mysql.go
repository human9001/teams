package config

import (
	"fmt"
	"time"
)

type mySQLConfig struct {
	Host            string        `yaml:"host"     env:"MYSQL_HOST"     env-default:"localhost"`
	Port            string        `yaml:"port"     env:"MYSQL_PORT"     env-default:"3306"`
	Database        string        `yaml:"database" env:"MYSQL_DB"       env-default:"teams"`
	User            string        `yaml:"user"     env:"MYSQL_USER"     env-default:"app"`
	Password        string        `yaml:"password" env:"MYSQL_PASSWORD" env-default:"password"`
	RootPassword    string        `yaml:"root_password"  env:"MYSQL_ROOT_PASSWORD"  env-default:"password"`
	MaxOpenConns    int           `yaml:"max_open_conns"  env:"MYSQL_MAX_OPEN_CONNS"  env-default:"25"`
	MaxIdleConns    int           `yaml:"max_idle_conns"  env:"MYSQL_MAX_IDLE_CONNS"  env-default:"10"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"  env:"MYSQL_CONN_MAX_LIFETIME"  env-default:"5m"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"  env:"MYSQL_CONN_MAX_IDLE_TIME"  env-default:"2m"`
}

func (c *mySQLConfig) DSN() string {
	return fmt.Sprintf("root:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		c.RootPassword, c.Host, c.Port, c.Database,
	)
}
