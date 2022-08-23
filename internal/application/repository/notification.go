package repository

import (
	"etna-notification/internal/infrastructure/mysql"
)

type INotificationRepository interface {
	CreateNotification(notification mysql.Notification) error
	IsAlreadyNotified(notification mysql.Notification) (bool, error)
	Close()
}
