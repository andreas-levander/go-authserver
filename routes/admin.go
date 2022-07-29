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

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Roles []string `json:"roles" binding:"required"`

}

func createUser(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body CreateUserRequest
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

		fmt.Println(string(pwHash))
		fmt.Println(len(string(pwHash)))
		newUser := database.User{
			User_id: 0,
			Username: body.Username,
			Password_hash: string(pwHash),
			Roles: body.Roles,
		}

		// pErr := bcrypt.CompareHashAndPassword(pwHash, []byte(body.Password))
		// if pErr != nil {
		// 	fmt.Fprintf(os.Stderr, "wrong pass: %v\n", cErr)
		// }
		
		if err := db.AddUser(newUser); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"database error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{ "message": "added new user: " + newUser.Username })
	}
}