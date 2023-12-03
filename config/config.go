package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Environment Environment `yaml:"env"`
	JwtOptions  JwtOptions  `yaml:"jwt"`
}

type Environment string

const (
	EnvironmentDev        Environment = "dev"
	EnvironmentProduction Environment = "prod"
)

type JwtOptions struct {
	Secret   string `yaml:"secret"`
	Lifetime int    `yaml:"lifetime"`
}

func GetFromYaml() (*Config, error) {
	config := &Config{}

	file, err := os.Open("../config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
