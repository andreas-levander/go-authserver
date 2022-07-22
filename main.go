package main

import (
	"github.com/gin-gonic/gin"

	"go-test/routes"

	"go-test/database"

	"github.com/subosito/gotenv"
)


func main() {
	gotenv.Load()

	router := gin.Default()

	routes.PublicRoutes(router)
	routes.AdminRoutes(router)

	pool := database.ConnectToDatabase()

	defer pool.Close()

	
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}