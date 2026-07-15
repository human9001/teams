package config

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"

	"github.com/human9001/teams/pkg/helpers"
)

var applicationConfig *config

type config struct {
	MySQL mySQLConfig `yaml:"pg"`
	HTTP  httpConfig  `yaml:"http"`
	App   appConfig   `yaml:"app"`
	Cache cacheConfig `yaml:"cache"`
}

const (
	defaultConfigPath = "config.local.yaml"
	rootProjectFile   = "go.mod"
	envFIle           = "app.env"
)

func AppConfig() *config {
	return applicationConfig
}

func LoadEnvFile() error {
	rootPath, err := helpers.ProjectRoot(rootProjectFile)
	if err != nil {
		slog.Error("не удалось найти корневую директорию проекта", "file", rootProjectFile, "error", err)
		return err
	}
	err = godotenv.Load(filepath.Join(rootPath, envFIle))
	if err != nil {
		slog.Error("не удалось загрузить .env файл", "file", envFIle, "error", err)
	}
	return nil
}

func ResolveConfigPath() string {
	var cfgFlag string
	flag.StringVar(&cfgFlag, "config", "", "путь к YAML-конфигу")
	flag.Parse()

	if cfgFlag != "" {
		return cfgFlag
	}

	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		return envPath
	}

	return defaultConfigPath
}

func MustLoad(path string) {
	var cfg config

	if path != "" {
		// ReadConfig читает YAML-файл, а затем перетирает значения из env-переменных
		// Приоритет: env > yaml > env-default
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			slog.Error("не удалось загрузить конфиг", "path", path, "error", err)
		}
	}

	// Если путь не указан — читаем только из env-переменных
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		slog.Error("не удалось загрузить env конфиг", "error", err)
	}

	applicationConfig = &cfg
}
