package routes

import (
	"fmt"
	"net/http"
	"os"

	"github.com/andreas-levander/go-authserver/config"
	"github.com/andreas-levander/go-authserver/database"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Public(router *gin.RouterGroup, env *config.Env) {
	public := router.Group("/public")

	{
		public.GET("/ping", ping)
		public.POST("/login", login(env))
		public.GET("/validate", validate(env))
	}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func login(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body loginRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			fmt.Fprintf(os.Stderr, "failed getting body params: %v\n", err)
			c.AbortWithStatusJSON(400, gin.H{
				"error": "missing params",
			})
			return

		}

		userDB := env.DB.GetUser(body.Username)

		var user *database.User

		if len(userDB) < 1 {
			c.AbortWithStatusJSON(404, gin.H{
				"error": "incorrect username or password",
			})
			return
		} else {
			user = userDB[0]
		}

		pErr := bcrypt.CompareHashAndPassword([]byte(user.Password_hash), []byte(body.Password))
		if pErr != nil {
			fmt.Fprintf(os.Stderr, "wrong pass: %v\n", pErr)
			c.AbortWithStatusJSON(404, gin.H{
				"error": "incorrect username or password",
			})
			return
		}
		
		c.JSON(200, gin.H{
			"token": env.Token.Create(user.Username, user.Roles, env.Config.TOKEN_TTL),
		})
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
func validate(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := env.Token.PublicKey()
		c.JSON(http.StatusOK, gin.H{ "keys": []jwk.Key{*key}})
	}
}