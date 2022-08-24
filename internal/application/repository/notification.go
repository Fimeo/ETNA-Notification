package repository

import (
	"etna-notification/internal/infrastructure/database"
)

type INotificationRepository interface {
	CreateNotification(notification database.Notification) error
	IsAlreadyNotified(notification database.Notification) (bool, error)
}

type Database interface {
	Close()
}
