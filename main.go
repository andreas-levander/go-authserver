package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"go-test/routes"
	"go-test/tokens"

	"go-test/database"

	"github.com/subosito/gotenv"
)


func main() {
	gotenv.Load()

	router := gin.Default()

	routes.PublicRoutes(router)
	
	db := database.Connect()

	routes.Admin(router, db)

	tokens.CreateKeys()
	token := tokens.CreateToken()
	claims, ok := tokens.ValidateToken(token)

	fmt.Println(claims.User, claims.Roles)
	fmt.Println(ok)


	//pool := database.Connect()

	//env := &Env{db: pool}

	//defer pool.Close()

	
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}