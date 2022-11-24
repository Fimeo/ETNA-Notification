package controller

import (
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
)

func AutoMigrateModels(userRep repository.IUserRepository, notificationRep repository.INotificationRepository) {
	userRep.Migrate()
	notificationRep.Migrate()
}

func CloseConnection(dg service.IDiscordService) {
	dg.CloseConnection()
}
