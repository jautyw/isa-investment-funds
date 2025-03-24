package config

import (
	"fmt"
	"log"
	"os"

	yml "gopkg.in/yaml.v2"
)

const (
	localConfig  = "config.yaml"
	dockerConfig = "config-docker.yaml"
)

// LoadConfig from local file.
func LoadConfig() (*Config, error) {

	namespace := localConfig

	if os.Getenv("NAMESPACE") == "docker" {
		namespace = dockerConfig
	}

	log.Println("Using namespace:", namespace)

	var cfg Config
	file, err := os.ReadFile(namespace)
	if err != nil {
		return nil, fmt.Errorf("error_reading %w", err)
	}
	err = yml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error_unmarshalling_config %w", err)
	}

	return &cfg, nil
}

// Config represents the configuration fields required for the application.
type Config struct {
	Host           string `yaml:"Host"`
	User           string `yaml:"User"`
	Password       string `yaml:"Password"`
	Database       string `yaml:"Database"`
	FundTableName  string `yaml:"FundTableName"`
	OrderTableName string `yaml:"OrderTableName"`
	Port           string `yaml:"Port"`
	SSLMode        string `yaml:"SSLMode"`
}
