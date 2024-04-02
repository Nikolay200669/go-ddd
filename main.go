package main

import (
	"go-ddd/domain"
	"go-ddd/infrastructure"
	"go-ddd/interfaces"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Инициализация базы данных
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&domain.User{}) // Автоматическое создание таблиц

	// Инициализация репозитория
	userRepository := infrastructure.NewUserRepository(db)

	// Инициализация маршрутизатора Gin
	r := gin.Default()

	// Определение обработчиков HTTP запросов
	r.POST("/users", interfaces.CreateUser(userRepository))
	r.GET("/users/:id", interfaces.GetUserByID(userRepository))
	// ping-pong
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Запуск сервера
	r.Run(":8080")
}
