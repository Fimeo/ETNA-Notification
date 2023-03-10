package controller

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron"

	"etna-notification/internal/domain"
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
	"etna-notification/internal/usecase"
)

type etnaNotificationController struct {
	DiscordService         service.IDiscordService
	EtnaWebService         service.IEtnaWebService
	UserRepository         repository.IUserRepository
	NotificationRepository repository.INotificationRepository
	CalendarRepository     repository.ICalendarEventRepository
}

type IEtnaNotificationController interface {
	StartDiscordNotificationCron()
	SendPushNotification() error
}

func NewEtnaNotificationController(repositories repository.Repositories, services service.Service) IEtnaNotificationController {
	return &etnaNotificationController{
		DiscordService:         services.IDiscordService,
		EtnaWebService:         services.IEtnaWebService,
		UserRepository:         repositories.IUserRepository,
		NotificationRepository: repositories.INotificationRepository,
		CalendarRepository:     repositories.ICalendarEventRepository,
	}
}

func (c *etnaNotificationController) StartDiscordNotificationCron() {
	callback := func() {
		err := c.SendPushNotification()
		if err != nil {
			usecase.SendErrorNotification(c.DiscordService, fmt.Sprintf("[ERROR] Something goes wrong during cron push notification: %+v", err))
			log.Fatalf("[ERROR] Something goes wrong during cron push notification: %+v", err)
			return
		}
	}
	cr := cron.New()
	err := cr.AddFunc("@every 30m", callback)
	if err != nil {
		log.Fatalf("[ERROR] Cannot create cron task : %+v", err)
	}
	cr.Start()
	go callback()
}

// SendPushNotification retrieve all registered user with a valid etna account linked to send unread notifications
// on a dedicated discord channel for each user.
func (c *etnaNotificationController) SendPushNotification() error {
	log.Print("[DEBUG] SendPushNotification triggered at : ", time.Now())
	users, err := c.UserRepository.FindAll()
	if err != nil {
		log.Printf("[ERROR] Cannot retrieve users from userRepository : %+v", err)
		return err
	}

	// Multithreading notifications retrieving and discord submission, works with a max goroutines in concurrency
	poolSize := 5
	sem := make(chan struct{}, poolSize)
	// Retrieve notifications only for user that have a linked discord channel
	for _, user := range users {
		if user.ChannelID == "" || user.Status != domain.StatusOpen {
			continue
		}
		// Shadow copy
		user := user

		// TODO : do not send error notification if intra is down or user credentials are wrong
		// TODO : split standard notification and calendar notifications.
		sem <- struct{}{}
		go func() {
			// Standard notifications
			err := usecase.SendPushNotificationForUser(user, c.NotificationRepository, c.UserRepository, c.EtnaWebService, c.DiscordService)
			if err != nil {
				usecase.SendErrorNotification(c.DiscordService, fmt.Sprintf("[ERROR] Something happens during cron push notification: %s", err.Error()))
				log.Printf("[ERROR] Something happens during cron push notification: %+v", err)
			}

			// Calendar notifications
			err = usecase.SendCalendarPushNotificationForUser(user, c.CalendarRepository, c.UserRepository, c.EtnaWebService, c.DiscordService)
			if err != nil {
				usecase.SendErrorNotification(c.DiscordService, fmt.Sprintf("[ERROR] Something happens during cron push calendar notification: %s", err.Error()))
				log.Printf("[ERROR] Something happens during cron push calendar notification: %+v", err)
			}
			<-sem
		}()
	}
	log.Print("[DEBUG] SendPushNotification end at : ", time.Now())
	return nil
}
