package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PublicRoutes(router *gin.Engine) {
	public := router.Group("/v1/public")

	{
		public.GET("/ping", ping)
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}