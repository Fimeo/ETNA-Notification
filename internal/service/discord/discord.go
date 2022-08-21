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

func DiscordConn(bottoken string) Service {
	dg, err := discordgo.New("Bot " + bottoken)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to discord bot : %+v", err))
	}

	return Service{Session: dg}
}

func (dg *Service) SendTextMessage(channel, message string) {
	messageSend, err := dg.Session.ChannelMessageSend(channel, message)
	if err != nil {
		return
	}
	log.Print("[DEBUG] Message sent at : ", messageSend.Timestamp)
}

func (dg *Service) CreateNotificationTextChannel(username string) {
	guildID := "984028659956473867" // The server guild ID
	channelCreate, err := dg.Session.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     username + "-notification",
		Type:     discordgo.ChannelTypeGuildText,
		Position: 2,
		ParentID: "984028659956473868", // The category Notification ID
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create guild channel : %+v", err))
	}
	fmt.Printf("Channel created : %+v", channelCreate)
}

func (dg *Service) GetChannel(channelID string) *discordgo.Channel {
	channel, err := dg.Session.Channel(channelID)
	if err != nil {
		panic(fmt.Sprintf("Failed retreive channel informations : %+v", err))
	}

	return channel
}

func (dg *Service) AddTextChannelNewReadingMember(memberID string, guildChannel *discordgo.Channel) {
	members := guildChannel.PermissionOverwrites // Retrieves default members from override category
	members = append(members, &discordgo.PermissionOverwrite{
		ID:    memberID, // Add the new user notification account
		Type:  discordgo.PermissionOverwriteTypeMember,
		Deny:  0,
		Allow: 1024, // Reading only
	})

	channel, err := dg.Session.ChannelEditComplex(guildChannel.ID, &discordgo.ChannelEdit{
		PermissionOverwrites: members,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to update guild channel parent : %+v", err))
	}
	fmt.Printf("Channel edited : %+v", channel)
}

func (dg *Service) OpenSocket() {
	// Open a websocket connection to Discord and begin listening.
	err := dg.Session.Open()
	if err != nil {
		panic(fmt.Sprintf("error opening connection : %+v", err))
	}
}

func (dg *Service) Close() {
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Session.Close()
}
