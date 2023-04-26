package config

import (
	"flag"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		TG   `yaml:"telegram"`
		DB
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Host string `env-required:"true" yaml:"host" env:"HTTP_HOST"`
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	TG struct {
		Token   string `env-required:"true" env:"TG_TOKEN"`
		Timeout int    `env-required:"true" yaml:"timeout"`
		Mode    string `env-required:"true" yaml:"mode"`
	}

	DB struct {
		Connection string `env-required:"true" env:"DB_CONNECTION"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	var path string
	flag.StringVar(&path, "config", "./config/config.yml", "Path to config")
	flag.Parse()

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("can't read config: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
