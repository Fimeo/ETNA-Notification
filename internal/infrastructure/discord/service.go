package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	Session *discordgo.Session
}

func NewConnection(token string) (*Service, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to discord bot : %+v", err)
		return nil, err
	}

	return &Service{Session: dg}, nil
}

func (s *Service) SendTextMessage(channelID, message string) (*discordgo.Message, error) {
	messageSend, err := s.Session.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Printf("[ERROR] Failed to send discord message : %+v", err)
		return nil, err
	}

	log.Print("[DEBUG] Message sent at : ", messageSend.Timestamp)
	return messageSend, nil
}

func (s *Service) CreateUserNotificationTextChannel(username string) (*discordgo.Channel, error) {
	guildID := "984028659956473867" // The server guild ID
	channelCreate, err := s.Session.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     username + "-notification",
		Type:     discordgo.ChannelTypeGuildText,
		Position: 2,
		ParentID: "984028659956473868", // The category Notification ID
	})
	if err != nil {
		log.Printf("[ERROR] Failed to create guild channel : %+v", err)
		return nil, err
	}

	log.Print("[DEBUG] Guild channel created : ", channelCreate.Name)
	return channelCreate, nil
}

func (s *Service) GetChannel(channelID string) (*discordgo.Channel, error) {
	channel, err := s.Session.Channel(channelID)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve guid channel : %+v", err)
		return nil, err
	}

	return channel, nil
}

func (s *Service) ChannelNewReadingMember(memberID, channelID string) (*discordgo.Channel, error) {
	channel, err := s.GetChannel(channelID)
	if err != nil {
		return nil, err
	}
	members := channel.PermissionOverwrites // Retrieves default members from override category
	members = append(members, &discordgo.PermissionOverwrite{
		ID:    memberID, // Add the new user notification account
		Type:  discordgo.PermissionOverwriteTypeMember,
		Deny:  0,
		Allow: 1024, // Reading only
	})

	channelUpdated, err := s.Session.ChannelEditComplex(channelID, &discordgo.ChannelEdit{
		PermissionOverwrites: members,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to add new channel member : %+v", err)
		return nil, err
	}

	log.Print("[DEBUG] New member added to channel : ", channelUpdated.Name)
	return channelUpdated, nil
}

func (s *Service) OpenSocket() error {
	// Open a websocket connection to Discord and begin listening.
	err := s.Session.Open()
	if err != nil {
		log.Printf("[ERROR] Failed to open web socket : %+v", err)
		return err
	}

	return nil
}

func (s *Service) Close() {
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	s.Session.Close()
}
