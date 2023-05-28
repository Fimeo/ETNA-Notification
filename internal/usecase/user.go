package usecase

import (
	"errors"
	"log"

	"etna-notification/internal/domain"
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
)

// AuthenticateUser or close account if bad credentials. Reuse authentication is user already authenticated.
func AuthenticateUser(
	user *domain.User,
	webService service.IEtnaWebService,
	userRepository repository.IUserRepository,
	discordService service.IDiscordService) error {
	// If the user has a valid authentication cookie, there is no need to log in again.
	if user.HasValidAuthentication() {
		log.Printf("[DEBUG] user %s already authenticated, no need to login again", user.Login)
		return nil
	}

	// Perform etna web service authentication to get authenticator cookie
	authenticationCookie, err := webService.LoginCookie(user.Login, user.Password)
	if err != nil {
		wrongCredError := &service.ErrorWrongCredentials{}
		if errors.As(err, wrongCredError) {
			return UserStopAccount(user, userRepository, discordService)
		}
	}

	user.SetAuthentication(authenticationCookie)

	return err
}

// UserStopAccount defined the user account stats to close and send a message to personal channel to inform.
func UserStopAccount(user *domain.User, userRepository repository.IUserRepository, discordService service.IDiscordService) error {
	user.Status = domain.StatusClose
	_, err := userRepository.UpdateStatus(user)
	if err != nil {
		return err
	}

	_, err = discordService.SendTextMessage(user.ChannelID, "Your account has been closed. You will not receive notifications anymore.")
	return err
}
