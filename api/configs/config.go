package configs

import (
	_"fmt"
	"log"
	"github.com/spf13/viper"
)

type Configs struct {
	MONGODB_URL string `mapstructure:"MONGODB_URL"`;
	PORT string `mapstructure:"PORT"`;
	SECRET string `mapstructure:"SECRET"`
}

// Load Envs
func Load() Configs {

	viper.AddConfigPath(".");
	viper.SetConfigName("../.env");
	viper.SetConfigType("env");

	viper.AutomaticEnv();

	err := viper.ReadInConfig();

	if err != nil {
		log.Fatal(err);
	}

	var configs Configs;

	viper.Unmarshal(&configs);

	return configs;
}

