package config

import (
	"fmt"
	"go-test/database"
	"go-test/tokens"

	"github.com/spf13/viper"
)

type Env struct {
	DB database.Database
	Token tokens.Tokens
	Config *config
}

type config struct {
	DB_URL string `mapstructure:"DATABASE_URL"`
	PORT string `mapstructure:"PORT"`
	TOKEN_TTL int `mapstructure:"TOKEN_TTL"`
}

func Load(path string) *config {
	viper.AddConfigPath(path)
	viper.SetConfigName("authserver")
	viper.SetConfigType("env")

	viper.SetDefault("PORT", 4000)
	viper.SetDefault("TOKEN_TTL", 15)

	viper.BindEnv("DATABASE_URL")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Viper error reading config: " + err.Error())
	}

	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Println("Viper unmarshal error: " + err.Error())
	} 
	return &cfg
	
}