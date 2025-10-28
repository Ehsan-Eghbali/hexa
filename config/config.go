package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig        `mapstructure:"app"`
	Databases []DatabaseConfig `mapstructure:"databases"`
}

func Get() *Config {
	return instance
}

type AppConfig struct {
	Debug bool `mapstructure:"debug"`
	Port  int  `mapstructure:"port"`
}

type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

var instance *Config

// LoadConfig initializes Viper and loads the configuration
func LoadConfig() error {

	configName := "config"

	RunEnv := os.Getenv("BART_ENV")
	fmt.Println(RunEnv)
	if RunEnv == "local" {
		configName = "local-config"
	}
	fmt.Println("Loading config from", configName)
	viper.SetConfigName(configName) // name of the config file (config.yaml)
	viper.SetConfigType("yaml")     // file type
	viper.AddConfigPath(".")        // path to look for the config file in

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Set up environment variable reading
	viper.AutomaticEnv() // Automatically override config values with environment variables

	// Now you can unmarshal the config
	if err := viper.Unmarshal(&instance); err != nil {
		return err
	}

	return nil
}
