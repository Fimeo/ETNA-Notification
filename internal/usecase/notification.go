package usecase

import (
	"log"

	"etna-notification/internal/domain"
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
)

func SendPushNotificationForUser(
	user *domain.User,
	notificationRep repository.INotificationRepository,
	etnaWS service.IEtnaWebService,
	discordS service.IDiscordService) error {
	authenticationCookie, err := etnaWS.LoginCookie(domain.BuildAuthenticationFromUser(user))
	if err != nil {
		return err
	}
	notifications, err := etnaWS.RetrievePendingNotifications(authenticationCookie, user.Login)
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		if notified, _ := notificationRep.IsNotified(domain.BuildNotificationFromEtnaNotificationAndUser(notification, user)); !notified {
			// Send notification is this case
			// TODO : factory to build formatted message notification
			_, err := discordS.SendTextMessage("1011372694241026098", notification.Message)
			if err != nil {
				log.Printf("[ERROR] Error when trying to send discord notification to user %+v and notification %+v", user, notification)
				continue
			}
			_, err = notificationRep.Save(domain.BuildNotificationFromEtnaNotificationAndUser(notification, user))
			if err != nil {
				return err
			}
			if err != nil {
				log.Printf("[ERROR] Error when saving notification for user %+v and notification %+v", user, notification)
				continue
			}
		} else {
			log.Print("[DEBUG] Notification already sent")
		}
	}

	return nil
}
