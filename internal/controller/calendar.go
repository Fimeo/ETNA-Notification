package controller

import (
	"time"

	"etna-notification/internal/domain"
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
)

type calendarController struct {
	EtnaWebService service.IEtnaWebService
	UserRepository repository.IUserRepository
}

type ICalendarController interface {
	GetCalendarEvent() []*domain.EtnaCalendarEvent
}

func NewCalendarController(repositories repository.Repositories, services service.Service) ICalendarController {
	return &calendarController{
		EtnaWebService: services.IEtnaWebService,
		UserRepository: repositories.IUserRepository,
	}
}

func (c *calendarController) GetCalendarEvent() []*domain.EtnaCalendarEvent {
	user, _ := c.UserRepository.FindByLogin("lefevr_a")
	cookie, _ := c.EtnaWebService.LoginCookie(user.Login, user.Password)
	events, _ := c.EtnaWebService.RetrieveCalendarEventInRange(cookie, user.Login, time.Now().Add(-time.Hour*24*2), time.Now().Add(time.Hour*24*2))
	return events
}
