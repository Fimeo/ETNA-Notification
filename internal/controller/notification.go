package controller

import (
	"log"
	"time"

	"github.com/robfig/cron"

	"etna-notification/internal/repository"
	"etna-notification/internal/service"
	"etna-notification/internal/usecase"
)

type etnaNotificationController struct {
	discordService         service.IDiscordService
	etnaWebService         service.IEtnaWebService
	userRepository         repository.IUserRepository
	notificationRepository repository.INotificationRepository
}

type IEtnaNotificationController interface {
	SendPushNotification() error
}

func NewEtnaNotificationController(
	discordService service.IDiscordService,
	etnaWebService service.IEtnaWebService,
	userRepository repository.IUserRepository,
	notificationRepository repository.INotificationRepository) IEtnaNotificationController {
	return &etnaNotificationController{
		discordService:         discordService,
		etnaWebService:         etnaWebService,
		userRepository:         userRepository,
		notificationRepository: notificationRepository,
	}
}

func StartPushNotificationCron(c IEtnaNotificationController) {
	cr := cron.New()
	err := cr.AddFunc("@every 30m", func() {
		err := c.SendPushNotification()
		if err != nil {
			log.Fatalf("[ERROR] Something happens during cron push notification: %+v", err)
			return
		}
	})
	if err != nil {
		log.Fatalf("[ERROR] Cannot create cron task : %+v", err)
	}
	err = c.SendPushNotification()
	if err != nil {
		log.Fatalf("[ERROR] Something happens during cron push notification: %+v", err)
	}
}

func (c *etnaNotificationController) SendPushNotification() error {
	log.Print("[DEBUG] SendPushNotification triggered at : ", time.Now())
	users, err := c.userRepository.FindAll()
	if err != nil {
		log.Printf("[ERROR] Cannot retrieve users from userRepository : %+v", err)
		return err
	}

	for _, user := range users {
		if user.ChannelID == "" {
			continue
		}
		user := user
		go func() {
			err := usecase.SendPushNotificationForUser(user, c.notificationRepository, c.etnaWebService, c.discordService)
			if err != nil {
				if err != nil {
					log.Printf("[ERROR] Something happens during cron push notification: %+v", err)
				}
			}
		}()
	}
	log.Print("[DEBUG] SendPushNotification end at : ", time.Now())
	return nil
}
