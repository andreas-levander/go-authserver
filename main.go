package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/andreas-levander/go-authserver/config"
	"github.com/andreas-levander/go-authserver/routes"
	"github.com/andreas-levander/go-authserver/tokens"

	"github.com/andreas-levander/go-authserver/database"
)


func main() {
	cfg := config.Load("./")

	router := gin.Default()

	env := &config.Env{DB: database.Connect(cfg.DB_URL), Token: tokens.CreateKeys(), Config: cfg}

	v1 := router.Group("/v1")

	routes.Admin(v1, env)
	routes.Public(v1, env)

	
	if err := router.Run(":" + viper.GetString("PORT")); err != nil {
		fmt.Println("error running server" + err.Error())
	}

}