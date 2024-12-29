package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	MongoDB MongoDBConfig
	HTTP    HTTPConfig
}

type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type HTTPConfig struct {
	Port string `mapstructure:"port"`
}

func Load() (*Config, error) {
	//viper.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	viper.SetDefault("mongodb.uri", "mongodb+srv://nicolaskesseli:A1WPWe2U85VC0A0X@cluster0.xdias.mongodb.net/")
	viper.SetDefault("mongodb.database", "run_management")
	viper.SetDefault("http.port", "8080")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
