package repository

import (
	"encoding/base64"
	"errors"
	"gorm.io/gorm"

	"etna-notification/internal/database"
	"etna-notification/internal/domain"
	"etna-notification/pkg/security"
)

type IUserRepository interface {
	Save(*domain.User) (*domain.User, error)
	Update(user *domain.User) (*domain.User, error)
	FindAll() ([]*domain.User, error)
	FindByDiscordName(string) (*domain.User, error)
	FindByLogin(string) (*domain.User, error)
	Migrate() error
}

type userRepository struct {
	database.Client
	security.Security
}

func NewUserRepository(client database.Client, security security.Security) IUserRepository {
	return &userRepository{
		client,
		security,
	}
}

// Save hash password user rsa public key because we need to restore the password to make
// web service authentication on etna api. This is not the most secure way, but it was the only mean
// to store hash password and restore them later. Hash if encoded into base64 to be stored in database.
func (ur *userRepository) Save(user *domain.User) (*domain.User, error) {
	encryptedPassword, err := security.Encrypt([]byte(user.Password), *ur.PublicKey)
	if err != nil {
		return nil, err
	}
	user.Password = base64.StdEncoding.EncodeToString(encryptedPassword)
	return user, ur.DB.Create(user).Error
}

func (ur *userRepository) Update(user *domain.User) (*domain.User, error) {
	return user, ur.DB.Updates(user).Error
}

// FindAll user retrieve all account registered to send notifications. Password are decrypted using rsa private key
// because we need to make authentication with clea password. The password was stored in base64, we need to decode before
// use decrypt.
func (ur *userRepository) FindAll() ([]*domain.User, error) {
	var users []*domain.User
	err := ur.DB.Find(&users).Error
	for _, user := range users {
		passDecode, err := base64.StdEncoding.DecodeString(user.Password)
		if err != nil {
			return nil, err
		}
		decryptPassword, err := security.Decrypt(passDecode, *ur.PrivateKey)
		if err != nil {
			return nil, err
		}
		user.Password = string(decryptPassword)
	}
	return users, err
}

func (ur *userRepository) FindByDiscordName(discordAccountName string) (*domain.User, error) {
	var user *domain.User
	err := ur.DB.First(&user, "discord_account = ?", discordAccountName).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New("an error occurred when find user by discord name")
	}
	return user, nil
}

func (ur *userRepository) FindByLogin(login string) (*domain.User, error) {
	var user *domain.User
	err := ur.DB.First(&user, "login = ?", login).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New("an error occurred when find user by login")
	}
	return user, nil
}

func (ur *userRepository) Migrate() error {
	return ur.DB.AutoMigrate(&domain.User{})
}
