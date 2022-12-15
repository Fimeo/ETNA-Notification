package usecase

import (
	"log"
	"os"

	"etna-notification/internal/domain"
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
)

// RetrieveUnreadNotificationForUser returns all unread notifications for the current user
func RetrieveUnreadNotificationForUser(user *domain.User, webService service.IEtnaWebService) ([]*domain.EtnaNotification, error) {
	// Perform etna web service authentication to get authenticator cookie
	authenticationCookie, err := webService.LoginCookie(user.Login, user.Password)
	// TODO : authentication failed count to revoke user status and send message on his personal channel
	if err != nil {
		return nil, err
	}
	// Retrieve unread notifications
	notifications, err := webService.RetrieveUnreadNotifications(authenticationCookie, user.Login)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

// SendPushNotificationForUser usecase will retrieve unread notifications from current user in etna web service (information).
// Then, compare notifications id with id already sent in notificationRepository. For each notification that is not already
// sent in user discord channel, send the notification.
func SendPushNotificationForUser(
	user *domain.User,
	notificationRepository repository.INotificationRepository,
	etnaWebService service.IEtnaWebService,
	discordService service.IDiscordService) error {
	notifications, err := RetrieveUnreadNotificationForUser(user, etnaWebService)
	if err != nil {
		return err
	}

	// If notification id was not found in notifications already sent, use discord service to send a new message in the user channel.
	for _, notification := range notifications {
		if notified, _ := notificationRepository.IsNotified(domain.BuildNotificationFromEtnaNotificationAndUser(notification, user)); notified {
			log.Print("[DEBUG] Notification already sent")
			continue
		}
		_, err := discordService.SendTextMessage(user.ChannelID, domain.BuildMessageFromEtnaNotification(notification))
		if err != nil {
			log.Printf("[ERROR] Error when trying to send discord notification to user %+v and notification %+v", user, notification)
			return err
		}
		_, err = notificationRepository.Save(domain.BuildNotificationFromEtnaNotificationAndUser(notification, user))
		if err != nil {
			return err
		}
		if err != nil {
			log.Printf("[ERROR] Error when saving notification for user %+v and notification %+v", user, notification)
			return err
		}
	}

	return nil
}

// SendErrorNotification send the message in the service.SystemErrorChannelID channel.
func SendErrorNotification(discordService service.IDiscordService, message string) {
	_, err := discordService.SendTextMessage(os.Getenv(service.SystemErrorChannelID), message)
	if err != nil {
		log.Printf("[ERROR] Error notification could not be sent %+v", err)
	}
}
