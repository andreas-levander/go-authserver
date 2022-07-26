package config

import (
	"fmt"

	"github.com/andreas-levander/go-authserver/database"
	"github.com/andreas-levander/go-authserver/tokens"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

type Env struct {
	DB database.Database
	Token tokens.Tokens
	Config *config
	Logger *zap.SugaredLogger
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

	viper.SetDefault("PORT", 8080)
	viper.SetDefault("TOKEN_TTL", 15)

	if err := viper.BindEnv("DATABASE_URL"); err != nil {
		fmt.Println("No database url found" + err.Error())
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("config not loaded: " + err.Error())
	}

	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Println("Viper unmarshal error: " + err.Error())
	} 
	return &cfg
	
}