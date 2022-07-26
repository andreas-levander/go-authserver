package routes

import (
	"net/http"

	"github.com/andreas-levander/go-authserver/config"
	"github.com/andreas-levander/go-authserver/database"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
)

func Admin(router *gin.RouterGroup, env *config.Env) {
	admin := router.Group("/admin", authMiddleware(env))

	{
		admin.GET("/users", users(env))
		admin.POST("/createuser", createUser(env))
	}
}
func authMiddleware(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header["Authorization"]
		if len(header) != 1 || len(header[0]) < 7 || header[0][:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		token := header[0][7:]
		if claims, ok := env.Token.Validate(token); ok && slices.Contains(claims.Roles, "admin") {
			c.Set("user", claims.User)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
	}
}

//for testing purposes
func users(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		usrs := env.DB.GetUsers()
		c.JSON(http.StatusOK, usrs)
	}
}

type CreateUserRequest struct {
	Username string   `json:"username" binding:"required"`
	Password string   `json:"password"`
	Roles    []string `json:"roles" binding:"required"`
}

func createUser(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body CreateUserRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			env.Logger.Debugf("failed getting body params %w", err)
			c.AbortWithStatusJSON(400, gin.H{
				"error": "missing params",
			})
			return

		}
		pwHash, cErr := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if cErr != nil {
			env.Logger.Errorf("failed generating password hash: %w", cErr)
			return
		}

		newUser := database.User{
			User_id:       0,
			Username:      body.Username,
			Password_hash: string(pwHash),
			Roles:         body.Roles,
		}

		if err := env.DB.AddUser(newUser); err != nil {
			env.Logger.Errorf("database error: %w", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"database error": err.Error(),
			})
			return
		}
		if user, ok := c.Get("user"); ok {
			env.Logger.Infof("%s added new user: %s", user, newUser.Username)
		}

		c.JSON(http.StatusOK, gin.H{"message": "added new user: " + newUser.Username})
	}
}
