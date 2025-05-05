package config

import "github.com/spf13/viper"

type Config struct {
	RESTHost string `mapstructure:"rest_host"`
	RESTPort string `mapstructure:"rest_port"`
	GinMode  string `mapstructure:"gin_mode"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("cfg/")

	viper.SetConfigName("app")
	viper.SetConfigType("env")
	err := viper.MergeInConfig()
	if err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	config := &Config{}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
