package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

type PostAuthentication struct {
	ID             int      `json:"id"`
	Login          string   `json:"login"`
	Email          string   `json:"email"`
	Logas          bool     `json:"logas"`
	Groups         []string `json:"groups"`
	LoginDate      string   `json:"login_date"`
	Firstconnexion bool     `json:"firstconnexion"`
	Password       string   `json:"password"`
}

type Information struct {
	ID          int         `json:"id"`
	Message     string      `json:"message"`
	Start       time.Time   `json:"start"`
	End         interface{} `json:"end"`
	CanValidate bool        `json:"can_validate"`
	Validated   bool        `json:"validated"`
	Type        string      `json:"type"`
	Metas       `json:"metas"`
}

type Metas struct {
	Type         string `json:"type"`
	SessionID    int    `json:"session_id,omitempty"`
	ActivityType string `json:"activity_type,omitempty"`
	ActivityID   int    `json:"activity_id,omitempty"`
	Promo        string `json:"promo,omitempty"`
}

const (
	LoginURL       = "https://auth.etna-alternance.net/identity"
	InformationURL = "https://intra-api.etna-alternance.net/students/%s/informations"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w \n", err))
	}

	fmt.Println(viper.AllKeys())

	information := retrieveInformation()
	s := discordBot()
	c := sqlite()
	for _, info := range information {
		if (notificationExist(c, notification{
			ExternalID: info.ID,
			User:       viper.GetString("etna.user"),
		})) {
			fmt.Println("Notification already sent")
		} else {
			sendDiscordMessage(s, viper.GetString("discord.channel"), info.Message)
			insertIntoDatabase(c, notification{
				ExternalID: info.ID,
				User:       viper.GetString("etna.user"),
			})
		}
	}
}

func retrieveInformation() []*Information {
	payload := PostAuthentication{
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
		panic("Cannot marshal the input payload")
	}

	headers := map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7",
		"Connection":      "keep-alive",
		"Content-Type":    "application/json;charset=UTF-8",
	}

	cookies, _ := exec(http.MethodPost, LoginURL, bodyJSON, headers, nil)

	if len(cookies) == 0 {
		panic("Connection failed, no cookie in response body")
	}

	_, resp := exec(http.MethodGet, fmt.Sprintf(InformationURL, viper.GetString("etna.user")), nil, nil, cookies[0])

	var information []*Information
	err = json.Unmarshal(resp, &information)
	if err != nil {
		panic("Marshal error")
	}

	return information
}

func exec(method, url string, body []byte, headers map[string]string, cookie *http.Cookie) ([]*http.Cookie, []byte) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		fmt.Println(err)
		panic("Error during request")
	}

	for key, header := range headers {
		req.Header.Add(key, header)
	}

	if cookie != nil {
		req.AddCookie(cookie)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		panic("Error during request")
	}
	defer res.Body.Close()

	bodyJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		panic("Error during request")
	}
	fmt.Println(res.StatusCode, string(bodyJSON))

	return res.Cookies(), bodyJSON
}

func discordBot() *discordgo.Session {
	s, err := discordgo.New("Bot " + viper.GetString("discord.bot-token"))
	if err != nil {
		panic(fmt.Sprintf("Bot error : %+v", err))
	}
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot is ready")
	})
	return s
}

func sendDiscordMessage(s *discordgo.Session, channel, message string) {
	messageSend, err := s.ChannelMessageSend(channel, message)
	if err != nil {
		return
	}
	fmt.Println("Message envoyé at : ", messageSend.Timestamp)
}

func sqlite() *sql.DB {
	db, err := sql.Open("sqlite3", viper.GetString("sqlite.file"))
	if err != nil {
		panic("Cannot open sql lite file")
	}
	return db
}

type notification struct {
	ID         int
	ExternalID int
	User       string
}

func insertIntoDatabase(db *sql.DB, notification notification) {
	_, err := db.Exec("INSERT INTO notification VALUES(NULL,datetime(),?, ?);", notification.ExternalID, notification.User)
	if err != nil {
		fmt.Println(err)
		panic("Cannot insert into database")
	}
}

func notificationExist(db *sql.DB, notif notification) bool {
	row, err := db.Query(
		"SELECT * FROM notification WHERE external_id=? and user=?",
		notif.ExternalID, notif.User)
	if err != nil {
		panic("Cannot read from database")
	}
	count := 0
	for row.Next() {
		count++
	}
	return count != 0
}
