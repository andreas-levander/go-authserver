package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"go-test/config"
	"go-test/routes"
	"go-test/tokens"

	"go-test/database"
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