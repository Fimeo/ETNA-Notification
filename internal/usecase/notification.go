package usecase

import (
	"log"
	"os"

	"etna-notification/internal/domain"
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
)

// RetrieveUnreadNotificationForUser returns all unread notifications for the current user
func RetrieveUnreadNotificationForUser(
	user *domain.User,
	webService service.IEtnaWebService,
	userRepository repository.IUserRepository,
	discordService service.IDiscordService) ([]*domain.EtnaNotification, error) {
	// Perform etna web service authentication to get authenticator cookie
	err := AuthenticateUser(user, webService, userRepository, discordService)
	if err != nil {
		return nil, err
	}
	// Retrieve unread notifications
	notifications, err := webService.RetrieveUnreadNotifications(user.GetAuthentication(), user.Login)
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
	userRepository repository.IUserRepository,
	etnaWebService service.IEtnaWebService,
	discordService service.IDiscordService) error {
	notifications, err := RetrieveUnreadNotificationForUser(user, etnaWebService, userRepository, discordService)
	if err != nil {
		return err
	}

	// If notification id was not found in notifications already sent, use discord service to send a new message in the user channel.
	for _, notification := range notifications {
		if notified, _ := notificationRepository.IsNotified(notification.BuildNotification(user)); notified {
			continue
		}
		_, err := discordService.SendTextMessage(user.ChannelID, notification.BuildNotificationMessage())
		if err != nil {
			log.Printf("[ERROR] Error when trying to send discord notification to user %+v and notification %+v", user, notification)
			return err
		}
		_, err = notificationRepository.Save(notification.BuildNotification(user))
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
