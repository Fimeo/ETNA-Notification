package repository

import (
	"etna-notification/internal/database"
	"etna-notification/internal/domain"
)

type IUserRepository interface {
	Save(user *domain.User) (*domain.User, error)
	FindAll() ([]*domain.User, error)
	Migrate() error
}

type userRepository struct {
	database.Client
}

func NewUserRepository(client database.Client) IUserRepository {
	return &userRepository{
		client,
	}
}

func (ur *userRepository) Save(user *domain.User) (*domain.User, error) {
	return user, ur.DB.Create(user).Error
}

func (ur *userRepository) FindAll() ([]*domain.User, error) {
	var users []*domain.User
	err := ur.DB.Find(&users).Error
	return users, err
}

func (ur *userRepository) Migrate() error {
	return ur.DB.AutoMigrate(&domain.User{})
}
