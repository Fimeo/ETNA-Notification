package usecase

import (
	"log"
	"time"

	"etna-notification/internal/domain"
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
	"etna-notification/internal/service/utils"
)

// RetrieveCalendarEventForUser returns all calendar event for the current day
func RetrieveCalendarEventForUser(
	user *domain.User,
	webService service.IEtnaWebService,
	userRepository repository.IUserRepository,
	discordService service.IDiscordService) ([]*domain.EtnaCalendarEvent, error) {
	// Perform etna web service authentication to get authenticator cookie
	err := AuthenticateUser(user, webService, userRepository, discordService)
	if err != nil {
		return nil, err
	}
	// Retrieve unread notifications
	now := time.Now().In(utils.GetParisLocation())
	calendarEvents, err := webService.RetrieveCalendarEventInRange(user.GetAuthentication(), user.Login, now, now)
	if err != nil {
		return nil, err
	}

	return calendarEvents, nil
}

// SendCalendarPushNotificationForUser usecase will retrieve calendar event for the current day and for current user in etna web service.
// If the event start begins in less than 30 minutes, a notification is sent with event information. Is the event is already sent,
// do not sent notification.
func SendCalendarPushNotificationForUser(
	user *domain.User,
	eventRepository repository.ICalendarEventRepository,
	userRepository repository.IUserRepository,
	etnaWebService service.IEtnaWebService,
	discordService service.IDiscordService) error {
	events, err := RetrieveCalendarEventForUser(user, etnaWebService, userRepository, discordService)
	if err != nil {
		return err
	}

	// If notification id was not found in notifications already sent, use discord service to send a new message in the user channel.
	// Event is sent only one hour before the event start.
	for _, event := range events {
		if !event.IsNotifiable() || !event.IsInNextHour() {
			continue
		}
		if notified, _ := eventRepository.IsNotified(event.BuildCalendarEvent(user)); notified {
			continue
		}
		_, err := discordService.SendTextMessage(user.ChannelID, event.BuildCalendarMessage())
		if err != nil {
			log.Printf("[ERROR] Error when trying to send discord event notification to user %+v and event %+v", user, event.Name)
			return err
		}
		_, err = eventRepository.Save(event.BuildCalendarEvent(user))
		if err != nil {
			return err
		}
		if err != nil {
			log.Printf("[ERROR] Error when saving event notification for user %+v and event %+v", user, event)
			return err
		}
	}

	return nil
}
