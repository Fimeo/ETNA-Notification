package usecase

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"time"

	"etna-notification/internal/domain"
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
)

// CreateServerInvitation returns the discord link with the server invitation. The invitation link
// has 1 day validity.
func CreateServerInvitation(discordService service.IDiscordService, channelID string) (*string, error) {
	log.Print("[DEBUG] Create invitation triggered at : ", time.Now())

	invitation, err := discordService.CreateInvitation(channelID)
	if err != nil {
		log.Printf("[ERROR] cannot create server invitation: %+v", err)
		return nil, fmt.Errorf("cannot create server invitation: %w", err)
	}

	invitationLink := fmt.Sprintf("https://discord.gg/%s", invitation.Code)
	return &invitationLink, nil
}

// CheckEtnaAccountValidity checks if the login password tuple is a valid etna account.
// Returns error if authentication failed
func CheckEtnaAccountValidity(etnaWebService service.IEtnaWebService, login, password string) error {
	// Check if the etna account if valid by making authentication
	if _, err := etnaWebService.LoginCookie(login, password); err != nil {
		log.Printf("[DEBUG]: user have bad credentials are register : %s", login)
		return fmt.Errorf("failed to authenticate on etna web service : %w", err)
	}

	return nil
}

// RegisterNewUser usecase will save the user into userRepository. The user objet was updated.
func RegisterNewUser(userRepository repository.IUserRepository, user *domain.User) error {
	// Authentication is valid, if user already exists update password
	var userFound *domain.User
	var err error
	if userFound, err = userRepository.FindByLogin(user.Login); err != nil {
		log.Printf("[ERROR]: user registration has failed on find user by login  : %+v", err)
		return fmt.Errorf("user registration has failed on find user by login : %w", err)
	}
	// If user was found in userRepository, update the password and discord account fields only
	if userFound != nil {
		userFound.Password = user.Password
		userFound.DiscordAccount = user.DiscordAccount
		_, err := userRepository.Update(user)
		if err != nil {
			log.Printf("[ERROR]: error occurred update user on register new user : %+v", err)
			return err
		}
	} else {
		if _, err := userRepository.Save(user); err != nil {
			log.Printf("[ERROR]: error occurred during user register : %+v", err)
			return err
		}
	}

	return nil
}

// CreatePersonalChannel
// TODO : insert all pending notifications into database to avoid massive assignment for user with hundreds of notifications
func CreatePersonalChannel(userRepository repository.IUserRepository, discordService service.IDiscordService, m *discordgo.MessageCreate) error {
	if m.Author.ID == discordService.Session().State.User.ID {
		return nil
	}
	discordName := fmt.Sprintf("%s#%s", m.Author.Username, m.Author.Discriminator)

	// Check if message was sent in service.ConnectChannelID and use /connect slash command.
	if m.ChannelID == os.Getenv(service.ConnectChannelID) && m.Content == "/connect" {
		// Check if user is register in userRepository, else ask to make register step
		user, err := userRepository.FindByDiscordName(discordName)
		if err != nil {
			discordService.SendTextMessageReply("error occurred during connect", m.Message)
			log.Printf("[ERROR] error occurred during connect user %s : %+v", discordName, err)
			return fmt.Errorf("[ERROR] error occurred during connect user %s : %w", discordName, err)
		}
		if user == nil {
			discordService.SendTextMessageReply("etna account was not found for you, please do register step", m.Message)
			log.Printf("[INFO] account was not found for %s on connect", discordName)
			return nil
		}

		// If the user has a channel linked, step over
		if user.ChannelID != "" {
			discordService.SendTextMessageReply("channel is already created", m.Message)
			log.Printf("[INFO] channel is already created for user %s : %+v", user.Login, err)
			return nil
		}

		// The user has no channel linked, create a user notification text channel and update the user.
		textChannel, err := discordService.CreateUserNotificationTextChannel(user, m.GuildID)
		if err != nil {
			discordService.SendTextMessageReply("error occurred, cannot create the personal channel", m.Message)
			log.Printf("[ERROR] error occurred, cannot create the personal channel %s : %+v", user.Login, err)
			return fmt.Errorf("[ERROR] error occurred, cannot create the personal channel %s : %w", user.Login, err)
		}

		user.ChannelID = textChannel.ID
		_, err = userRepository.Update(user)
		if err != nil {
			discordService.SendTextMessageReply("error occurred, cannot link the personal channel", m.Message)
			log.Printf("[ERROR] error occured during on save personal channel for user %s %s %+v", user.Login, textChannel.ID, err)
			return fmt.Errorf("[ERROR] error occured during on save personal channel for user %s %s %w", user.Login, textChannel.ID, err)
		}

		// Add read role to channel
		_, err = discordService.ChannelNewReadingMember(m.Author.ID, textChannel.ID)
		if err != nil {
			discordService.SendTextMessageReply("error occurred, cannot grant access to personal channel", m.Message)
			log.Printf("[ERROR] : error occured during connect to grant role privilege %s %s %+v", user.Login, textChannel.ID, err)
			return fmt.Errorf("[ERROR] : error occured during connect to grant role privilege %s %s %w", user.Login, textChannel.ID, err)
		}

		// Everything if ok, send confirmation message
		discordService.SendTextMessageReply("personal channel created, you will receive notifications soon", m.Message)
	} else {
		discordService.SendTextMessageReply("command not found, use `/connect`", m.Message)
	}
	return nil
}
