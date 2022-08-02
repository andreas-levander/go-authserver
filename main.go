package main

import (
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

	
	router.Run(":" + viper.GetString("PORT")) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}