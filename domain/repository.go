package domain

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByID(id uint) (*User, error)
	// Добавьте другие методы репозитория по необходимости
}
