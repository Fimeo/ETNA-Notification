package controller

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"

	"etna-notification/internal/domain"
	"etna-notification/internal/repository"
	"etna-notification/internal/service"
)

type registerController struct {
	discordService         service.IDiscordService
	etnaWebService         service.IEtnaWebService
	userRepository         repository.IUserRepository
	notificationRepository repository.INotificationRepository
}

type IRegisterController interface {
	CreateInvitation() (string, error)
	AddConnectHandler()
	RegisterNewUser(login, password, discordAccount string) (string, error)
}

func NewRegisterController(
	discordService service.IDiscordService,
	etnaWebService service.IEtnaWebService,
	userRepository repository.IUserRepository,
	notificationRepository repository.INotificationRepository) IRegisterController {
	return &registerController{
		discordService:         discordService,
		etnaWebService:         etnaWebService,
		userRepository:         userRepository,
		notificationRepository: notificationRepository,
	}
}

func Register(c IRegisterController) {
	c.AddConnectHandler()
}

func (c *registerController) RegisterNewUser(login, password, discordAccount string) (string, error) {
	user := &domain.User{
		DiscordAccount: discordAccount,
		Login:          login,
		Password:       password,
	}
	// 1 : check auth etna and retrieve user id
	if _, err := c.etnaWebService.LoginCookie(user); err != nil {
		return "", errors.New(fmt.Sprintf("failed to authenticate on etna web service : %+v", err))
	}
	// 2 : if success : register user
	if userFound, err := c.userRepository.FindByLogin(login); userFound != nil || err != nil {
		return "", errors.New(fmt.Sprintf("user is already registered or an error ocurred : %+v", err))
	}

	if _, err := c.userRepository.Save(user); err != nil {
		log.Printf("[ERROR]: error occurred during user register : %+v", err)
		return "", errors.New("error occurred during user register")
	}

	// 3 : create invitation link
	link, err := c.CreateInvitation()
	if err != nil {
		return "", err
	}

	return link, nil
}

func (c *registerController) CreateInvitation() (string, error) {
	log.Print("[DEBUG] Create invitation triggered at : ", time.Now())

	invitation, err := c.discordService.CreateInvitation()
	if err != nil {
		log.Fatalf("[ERROR] Something happens during new invitation: %+v", err)
	}

	return fmt.Sprintf("https://discord.gg/%s", invitation.Code), nil
}

func (c *registerController) AddConnectHandler() {
	c.discordService.GetCurrentSession().AddHandler(ConnectUser(c))
}

func ConnectUser(c *registerController) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO : send notifications on failure
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.ChannelID == "1050392032079786077" && m.Content == "/connect" {
			// 1 : test is user is registered
			user, err := c.userRepository.FindByDiscordName(fmt.Sprintf("%s#%s", m.Author.Username, m.Author.Discriminator))
			if err != nil {
				log.Printf("[ERROR] : error occured during connect user %+v : %+v", user, err)
				return
			}

			// 2 : if there is no channel linked, create a text channel for him
			if user.ChannelID != "" {
				// TODO : insert account status (credentials expired renew)
				return
			}
			textChannel, err := c.discordService.CreateUserNotificationTextChannel(user)
			if err != nil {
				log.Printf("[ERROR] : error occured during connect to create text channel %+v : %+v", user, err)
				return
			}

			user.ChannelID = textChannel.ID
			_, err = c.userRepository.Update(user)
			if err != nil {
				log.Printf("[ERROR] : error occured during connect to save the channel created %+v %+v %+v", user, textChannel, err)
				return
			}

			// Add read role to channel
			_, err = c.discordService.ChannelNewReadingMember(m.Author.ID, textChannel.ID)
			if err != nil {
				log.Printf("[ERROR] : error occured during connect to grant role privilege %+v %+v %+v", user, textChannel, err)
				return
			}

			// 3 : send message ok done
			_, err = s.ChannelMessageSendReply("1050392032079786077", "ok", &discordgo.MessageReference{ // TODO : use configuration
				MessageID: m.ID,
				ChannelID: m.ChannelID,
				GuildID:   m.GuildID,
			})
			if err != nil {
				log.Printf("[ERROR] : error occured during text reply ok : %+v", err)
				return
			}
		}
	}
}
