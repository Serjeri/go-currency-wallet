package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Database   `yaml:"database"`
	HTTPServer `yaml:"http_server"`
}

type Database struct {
	Dburl string `yaml:"dburl"`
}

type HTTPServer struct {
	Address     string `yaml:"address" env-default:"localhost:8080"`
	Addressgrpc string `yaml:"addressgrpc" env-default:"localhost:50051"`
}

func MustLoad() *Config {
	configPath := "../../config/config.yml"
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(yamlFile, &cfg); err != nil {
		log.Fatalf("cannot parse YAML config: %s", err)
	}

	return &cfg
}
