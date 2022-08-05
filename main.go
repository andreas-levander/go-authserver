package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/andreas-levander/go-authserver/config"
	"github.com/andreas-levander/go-authserver/database"
	"github.com/andreas-levander/go-authserver/routes"
	"github.com/andreas-levander/go-authserver/tokens"
	logger "github.com/andreas-levander/go-authserver/util"
)

func main() {
	cfg := config.Load("./")

	router := gin.Default()

	logger := logger.Load()
	defer logger.Sync()

	keys, err := tokens.CreateKeys()
	if err != nil {
		logger.Fatal(err)
	}

	env := &config.Env{DB: database.Connect(cfg.DB_URL), Token: keys, Config: cfg, Logger: logger}

	v1 := router.Group("/v1")

	routes.Admin(v1, env)
	routes.Public(v1, env)

	if err := router.Run(":" + viper.GetString("PORT")); err != nil {
		logger.Fatal("error running server: " + err.Error())
	}

}
