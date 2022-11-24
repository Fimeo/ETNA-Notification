package service

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

// DiscordService is a proxy using discord-go package.
type discordService struct {
	DG *discordgo.Session
}

type IDiscordService interface {
	SendTextMessage(channelID, message string) (*discordgo.Message, error)
	CreateUserNotificationTextChannel(username string) (*discordgo.Channel, error)
	GetChannel(channelID string) (*discordgo.Channel, error)
	CloseConnection()
}

// NewDiscordService init connection with discord using discord-go package. A web socket is open.
func NewDiscordService() IDiscordService {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		log.Panicf("[ERROR] Failed to init discord connection using token : %+v", err)
	}

	err = dg.Open()
	if err != nil {
		log.Panicf("[ERROR] Failed to open discord web socket : %+v", err)
	}

	return &discordService{DG: dg}
}

func (dg *discordService) SendTextMessage(channelID, message string) (*discordgo.Message, error) {
	messageSend, err := dg.DG.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Printf("[ERROR] Failed to send discord message : %+v", err)
		return nil, err
	}

	log.Print("[DEBUG] Message sent at : ", messageSend.Timestamp)
	return messageSend, nil
}

func (dg *discordService) CreateUserNotificationTextChannel(username string) (*discordgo.Channel, error) {
	guildID := "984028659956473867" // The server guild ID
	channelCreate, err := dg.DG.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
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

func (dg *discordService) GetChannel(channelID string) (*discordgo.Channel, error) {
	channel, err := dg.DG.Channel(channelID)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve guid channel : %+v", err)
		return nil, err
	}

	return channel, nil
}

func (dg *discordService) ChannelNewReadingMember(memberID, channelID string) (*discordgo.Channel, error) {
	channel, err := dg.GetChannel(channelID)
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

	channelUpdated, err := dg.DG.ChannelEditComplex(channelID, &discordgo.ChannelEdit{
		PermissionOverwrites: members,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to add new channel member : %+v", err)
		return nil, err
	}

	log.Print("[DEBUG] New member added to channel : ", channelUpdated.Name)
	return channelUpdated, nil
}

func (dg *discordService) CloseConnection() {
	dg.DG.Close()
}
