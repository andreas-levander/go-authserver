package config

import (
	"fmt"
	"go-test/database"
	"go-test/tokens"

	"github.com/spf13/viper"
)

type Env struct {
	DB *database.DB
	Token *tokens.Keys
	Config *config
}

type config struct {
	DB_URL string `mapstructure:"DATABASE_URL"`
	PORT string `mapstructure:"PORT"`
	TOKEN_TTL int `mapstructure:"TOKEN_TTL"`
}

func Setup() *config {
	viper.SetDefault("PORT", 4000)
	viper.SetDefault("TOKEN_TTL", 15)

	viper.BindEnv("DATABASE_URL")
	viper.AutomaticEnv()

	all := viper.AllKeys()
	fmt.Println(all)
	fmt.Println(viper.Get("token_ttl"))

	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Println("Viper unmarshal error: " + err.Error())
	} 
	return &cfg
	
}