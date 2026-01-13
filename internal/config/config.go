package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)


type Site struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Config struct {
	Interval time.Duration `yaml:"interval"`
	Sites    []Site        `yaml:"sites"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config file %q: %w", path, err)
	}

	if len(cfg.Sites) == 0 {
		return nil, fmt.Errorf("config file %q: sites list is empty", path)
	}

	if cfg.Interval <= 0 {
		return nil, fmt.Errorf("config file %q: interval must be greater than zero", path)
	}

	return &cfg, nil
}