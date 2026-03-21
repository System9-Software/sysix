package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type PanelConfig struct {
	System    bool `yaml:"system"`
	Processes bool `yaml:"processes"`
	Ports     bool `yaml:"ports"`
	Network   bool `yaml:"network"`
}

type WebConfig struct {
	Port    int    `yaml:"port"`
	Host    string `yaml:"host"`
	Enabled bool   `yaml:"enabled"`
}

type Config struct {
	RefreshRate int         `yaml:"refresh_rate"`
	Panels      PanelConfig `yaml:"panels"`
	Web         WebConfig   `yaml:"web"`
}

func Default() *Config {
	return &Config{
		RefreshRate: 2,
		Panels: PanelConfig{
			System:    true,
			Processes: true,
			Ports:     true,
			Network:   true,
		},
		Web: WebConfig{
			Port:    8080,
			Host:    "localhost",
			Enabled: true,
		},
	}
}

func Load() (*Config, error) {
	cfg := Default()

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return cfg, nil
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
