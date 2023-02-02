package controller

import (
	"fmt"
	"log"
	"sync"
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
	}
}

func (c *etnaNotificationController) StartDiscordNotificationCron() {
	cr := cron.New()
	err := cr.AddFunc("@every 1m", func() {
		err := c.SendPushNotification()
		if err != nil {
			usecase.SendErrorNotification(c.DiscordService, fmt.Sprintf("[ERROR] Something goes wrong during cron push notification: %+v", err))
			log.Fatalf("[ERROR] Something goes wrong during cron push notification: %+v", err)
			return
		}
	})
	if err != nil {
		log.Fatalf("[ERROR] Cannot create cron task : %+v", err)
	}
	cr.Start()
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

	// Wait group for multithreading notifications retrieving and discord submission
	var wg sync.WaitGroup
	// Retrieve notifications only for user that have a linked discord channel
	for _, user := range users {
		if user.ChannelID == "" || user.Status != domain.StatusOpen {
			continue
		}
		// Shadow copy
		user := user

		wg.Add(1)
		go func() {
			// TODO : catch error in channel list
			err := usecase.SendPushNotificationForUser(user, c.NotificationRepository, c.EtnaWebService, c.DiscordService)
			if err != nil {
				usecase.SendErrorNotification(c.DiscordService, fmt.Sprintf("[ERROR] Something happens during cron push notification: %s", err.Error()))
				log.Printf("[ERROR] Something happens during cron push notification: %+v", err)
			}
			defer wg.Done()
		}()
	}
	wg.Wait()
	log.Print("[DEBUG] SendPushNotification end at : ", time.Now())
	return nil
}
