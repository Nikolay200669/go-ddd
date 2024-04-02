package main

import (
	"bytes"
	"go-ddd/domain"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"go-ddd/infrastructure"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAPI(t *testing.T) {
	// Set up a test database
	db, err := SetupTestDatabase()
	if err != nil {
		t.Fatalf("failed to set up test database: %v", err)
	}

	// Initialize UserRepository
	userRepository := infrastructure.NewUserRepository(db)

	// Initialize Gin router
	r := gin.Default()
	r.POST("/users", CreateUser(userRepository))
	r.GET("/users/:id", GetUserByID(userRepository))

	// Perform POST request to create a user
	createUserRequestBody := `{"name": "John Doe", "email": "john@example.com", "password": "password123"}`
	createUserRequest, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(createUserRequestBody))
	createUserRequest.Header.Set("Content-Type", "application/json")
	createUserRecorder := httptest.NewRecorder()
	r.ServeHTTP(createUserRecorder, createUserRequest)

	// Check response status code
	assert.Equal(t, http.StatusCreated, createUserRecorder.Code)

	// Perform GET request to get the created user
	getUserRecorder := httptest.NewRecorder()
	getUserRequest, _ := http.NewRequest("GET", "/users/1", nil)
	r.ServeHTTP(getUserRecorder, getUserRequest)

	// Check response status code
	assert.Equal(t, http.StatusOK, getUserRecorder.Code)

	// Check response body
	assert.Contains(t, getUserRecorder.Body.String(), "John Doe")
	assert.Contains(t, getUserRecorder.Body.String(), "john@example.com")
}

// SetupTestDatabase настраивает временную тестовую базу данных и возвращает подключение к ней.
func SetupTestDatabase() (*gorm.DB, error) {
	// Подключение к тестовой базе данных SQLite в памяти
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Автомиграция модели User
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateUser(ur domain.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user domain.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Создание пользователя через UserRepository
		if err := ur.CreateUser(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func GetUserByID(ur domain.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Получение пользователя по его идентификатору через UserRepository
		user, err := ur.GetUserByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
