package controller

import (
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
)

func Setup(repositories repository.Repositories, services service.Service) {
	notificationCtrl := NewEtnaNotificationController(repositories, services)
	registerCtrl := NewRegisterController(repositories, services)

	notificationCtrl.StartDiscordNotificationCron()
	registerCtrl.Connect()
}
