package config

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Environment  Environment  `yaml:"env"`
	JwtOptions   JwtOptions   `yaml:"jwt"`
	MongoOptions MongoOptions `yaml:"mongo"`
	KafkaOptions KafkaOptions `yaml:"kafka"`
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

type MongoOptions struct {
	ConnectionString string `yaml:"connectionString"`
	Database         string `yaml:"database"`
}

type KafkaOptions struct {
	Brokers []string `yaml:"brokers"`
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

func (config Config) GetMongoDb(ctx context.Context) (*mongo.Database, error) {
	mongoOptions := options.Client().ApplyURI(config.MongoOptions.ConnectionString)
	mongoOptions.TLSConfig.InsecureSkipVerify = true

	mongoClient, err := mongo.Connect(ctx, mongoOptions)
	if err != nil {
		return nil, err
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return mongoClient.Database(config.MongoOptions.Database), nil
}
