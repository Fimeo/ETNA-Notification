package controller

import (
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
)

type Controllers struct {
	IRegisterController
	IEtnaNotificationController
	ICalendarController
}

func InitControllers(repositories repository.Repositories, services service.Service) Controllers {
	notificationCtrl := NewEtnaNotificationController(repositories, services)
	registerCtrl := NewRegisterController(repositories, services)
	calendarCtrl := NewCalendarController(repositories, services)

	return Controllers{
		IRegisterController:         registerCtrl,
		IEtnaNotificationController: notificationCtrl,
		ICalendarController:         calendarCtrl,
	}
}

func Setup(controllers Controllers) {
	controllers.IEtnaNotificationController.StartDiscordNotificationCron()

	controllers.IRegisterController.Connect()
	controllers.IRegisterController.Stop()
}
