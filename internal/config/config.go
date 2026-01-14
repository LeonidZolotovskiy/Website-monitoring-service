package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Site struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Config struct {
	Interval time.Duration `yaml:"interval" env:"CHECK_INTERVAL"`
	Sites    []Site        `yaml:"sites"`

	Port     int    `yaml:"port" env:"APP_PORT"`
	LogLevel string `yaml:"log_level" env:"LOG_LEVEL"`

	DB DBConfig `yaml:"database"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     int    `yaml:"port" env:"DB_PORT"`
	User     string `yaml:"user" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Name     string `yaml:"name" env:"DB_NAME"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE"`
}

func Load(path string) (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("cannot read env: %w", err)
	}

	if len(cfg.Sites) == 0 {
		return nil, fmt.Errorf("sites list is empty")
	}

	if cfg.Interval <= 0 {
		return nil, fmt.Errorf("interval must be greater than zero")
	}

	return &cfg, nil
}
