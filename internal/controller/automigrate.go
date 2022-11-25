package controller

import (
	"etna-notification/internal/repository"
)

func AutoMigrateModels(userRep repository.IUserRepository, notificationRep repository.INotificationRepository) {
	userRep.Migrate()
	notificationRep.Migrate()
}
