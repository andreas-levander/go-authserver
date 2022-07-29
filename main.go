package main

import (
	"github.com/gin-gonic/gin"

	"go-test/config"
	"go-test/routes"
	"go-test/tokens"

	"go-test/database"

	"github.com/subosito/gotenv"
)


func main() {
	gotenv.Load()

	router := gin.Default()

	env := &config.Env{DB: database.Connect(), Token: tokens.CreateKeys()}

	routes.Admin(router, env)
	routes.Public(router, env)

	
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}