package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron"
	"github.com/spf13/viper"

	"etna-scrapping/internal/model"
	"etna-scrapping/internal/model/sqlite"
	"etna-scrapping/internal/service"
)

const (
	LoginURL       = "https://auth.etna-alternance.net/identity"
	InformationURL = "https://intra-api.etna-alternance.net/students/%s/informations"
)

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func main() {
	loadConfig()

	f := service.InitLogger()
	defer f.Close()

	log.Print("[DEBUG] Up at ", time.Now())

	s := service.DiscordConn(viper.GetString("discord.bot-token"))
	c := service.SQLLiteConn(viper.GetString("sqlite.file"))

	cr := cron.New()
	cr.AddFunc("@every 30m", func() {
		execute(s, c)
	})
	cr.Start()
	execute(s, c)

	for {
		time.Sleep(time.Second)
	}
}

func execute(s service.DiscordService, c service.SQLLiteService) {
	log.Print("[DEBUG] Execute : ", time.Now())
	information := retrieveInformation()
	for _, info := range information {
		if (c.IsAlreadyNotified(sqlite.Notification{
			ExternalID: info.ID,
			User:       viper.GetString("etna.user"),
		})) {
			log.Print("[DEBUG] Notification already sent")
		} else {
			s.SendDiscordMessage(viper.GetString("discord.channel"), info.Message)
			c.CreateNotification(sqlite.Notification{
				ExternalID: info.ID,
				User:       viper.GetString("etna.user"),
			})
		}
	}
}

func retrieveInformation() []*model.Notification {
	payload := model.Authentication{
		ID:             viper.GetInt("etna.user_id"),
		Login:          viper.GetString("etna.user"),
		Email:          viper.GetString("etna.user") + "@etna-alternance.net",
		Logas:          false,
		Groups:         []string{"student"},
		LoginDate:      time.Now().Format("2006-01-02 15-04-05"),
		Firstconnexion: false,
		Password:       viper.GetString("etna.password"),
	}

	bodyJSON, err := json.Marshal(payload)
	if err != nil {
		log.Print("[ERROR] Marshal error")
		return nil
	}

	headers := map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7",
		"Connection":      "keep-alive",
		"Content-Type":    "application/json;charset=UTF-8",
	}

	cookies, _, err := service.Exec(http.MethodPost, LoginURL, bodyJSON, headers, nil)
	if err != nil {
		return nil
	}

	if len(cookies) == 0 {
		log.Print("[ERROR] Connection failed, no cookie in response body")
		return nil
	}

	_, resp, err := service.Exec(
		http.MethodGet,
		fmt.Sprintf(InformationURL, viper.GetString("etna.user")),
		nil,
		nil,
		cookies[0])
	if err != nil {
		return nil
	}

	var notification []*model.Notification
	err = json.Unmarshal(resp, &notification)
	if err != nil {
		log.Print("[ERROR] UnMarshal error")
		return nil
	}

	return notification
}
