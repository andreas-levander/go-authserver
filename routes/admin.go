package routes

import (
	"go-test/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Admin(router *gin.Engine, db *database.DB) {
	admin := router.Group("/v1/admin")

	{
		admin.GET("/users", users(db))
	}
}
func users(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		usrs := db.GetUsers()
		c.JSON(http.StatusOK, usrs)
	}
}