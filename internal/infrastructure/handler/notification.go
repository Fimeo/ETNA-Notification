package handler

import (
	"log"
	"time"

	"etna-notification/internal/application"
	"etna-notification/internal/infrastructure/database"
	"etna-notification/internal/infrastructure/etna"
)

// TODO : refacto this part into usecase and handle errors

// Next dev : create channels for new users (add isNew column)
// Create invitation links
// Rest api to add new users
// Hash password
// Use a Queue Pub sub event to send notifications

func SendNewNotifications(dependencies application.Dependencies) {
	log.Print("[DEBUG] Triggered at : ", time.Now())
	users, _ := dependencies.Users.GetEtnaUsers()
	authentications := etna.BuildAuthenticationFromUsers(users)

	for _, user := range authentications {
		cookie, _ := dependencies.Etna.Login(user)
		notifications, _ := dependencies.Etna.RetrieveNotifications(cookie, user.Login)

		for _, notification := range notifications {
			if notified, _ := (dependencies.Notification.IsAlreadyNotified(database.Notification{
				ExternalID: notification.ID,
				User:       user.Login,
			})); !notified {
				dependencies.Discord.SendTextMessage("1011372694241026098", notification.Message)
				dependencies.Notification.CreateNotification(database.Notification{
					ExternalID: notification.ID,
					User:       user.Login,
				})
			} else {
				log.Print("[DEBUG] Notification already sent")
			}
		}
	}
}
