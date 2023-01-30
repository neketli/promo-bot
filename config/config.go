package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		TG   `yaml:"telegram"`
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
		Token   string `env-required:"true" yaml:"token" env:"TG_TOKEN"`
		Timeout int    `env-required:"true" yaml:"timeout"`
		Mode    string `env-required:"true" yaml:"mode"`
	}
)

func New(path string) (*Config, error) {
	cfg := &Config{}

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
