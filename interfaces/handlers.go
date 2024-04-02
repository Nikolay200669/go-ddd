package interfaces

import (
	"net/http"
	"strconv"

	"go-ddd/domain"

	"github.com/gin-gonic/gin"
)

func CreateUser(ur domain.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user domain.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := ur.CreateUser(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func GetUserByID(ur domain.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		user, err := ur.GetUserByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
