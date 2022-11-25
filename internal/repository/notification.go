package repository

import (
	"errors"
	"gorm.io/gorm"

	"etna-notification/internal/database"
	"etna-notification/internal/domain"
)

type INotificationRepository interface {
	Save(notification *domain.Notification) (*domain.Notification, error)
	FindAll() ([]*domain.Notification, error)
	IsNotified(notification *domain.Notification) (bool, error)
	Migrate() error
}

type notificationRepository struct {
	database.Client
}

func NewNotificationRepository(client database.Client) INotificationRepository {
	return &notificationRepository{
		client,
	}
}

func (nr *notificationRepository) Save(notification *domain.Notification) (*domain.Notification, error) {
	return notification, nr.DB.Create(notification).Error
}

func (nr *notificationRepository) FindAll() ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	err := nr.DB.Find(&notifications).Error
	return notifications, err
}

func (nr *notificationRepository) IsNotified(notification *domain.Notification) (bool, error) {
	result := nr.DB.Take(&notification, "external_id = ?", notification.ExternalID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}

	return true, result.Error
}

func (nr *notificationRepository) Migrate() error {
	return nr.DB.AutoMigrate(&domain.Notification{})
}
