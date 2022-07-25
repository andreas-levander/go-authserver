package routes

import (
	"fmt"
	"go-test/database"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
)

func Admin(router *gin.Engine, db *database.DB) {
	admin := router.Group("/v1/admin")

	{
		admin.GET("/users", users(db))
		admin.POST("/createuser", createUser(db))
	}
}
func users(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		usrs := db.GetUsers()
		c.JSON(http.StatusOK, usrs)
	}
}

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Roles []string `json:"roles" binding:"required"`

}

func createUser(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body createUserRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			fmt.Fprintf(os.Stderr, "failed getting body params: %v\n", err)
			c.AbortWithStatusJSON(400, gin.H{
				"error": "missing params",
			})
			return

		}
		fmt.Println(body)
		pwHash, cErr := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "failed hashing password: %v\n", cErr)
		}

		fmt.Println(pwHash)

		pErr := bcrypt.CompareHashAndPassword(pwHash, []byte(body.Password))
		if pErr != nil {
			fmt.Fprintf(os.Stderr, "wrong pass: %v\n", cErr)
		}
		//db.AddUser()
	}
}