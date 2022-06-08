package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron"
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

type Notification struct {
	ID         int
	ExternalID int
	User       string
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
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	f := openLogFile()
	defer f.Close()

	log.SetOutput(f)

	log.Print("[DEBUG] Up at ", time.Now())

	s := discordBot()
	c := sqlite()

	cr := cron.New()
	cr.AddFunc("@every 30m", func() {
		log.Print("[DEBUG] Execute : ", time.Now())
		information := retrieveInformation()
		for _, info := range information {
			if (notificationExist(c, Notification{
				ExternalID: info.ID,
				User:       viper.GetString("etna.user"),
			})) {
				log.Print("[DEBUG] Notification already sent")
			} else {
				sendDiscordMessage(s, viper.GetString("discord.channel"), info.Message)
				insertIntoDatabase(c, Notification{
					ExternalID: info.ID,
					User:       viper.GetString("etna.user"),
				})
			}
		}
	})
	cr.Start()

	for {
		time.Sleep(time.Second)
	}
}

func discordBot() *discordgo.Session {
	s, err := discordgo.New("Bot " + viper.GetString("discord.bot-token"))
	if err != nil {
		panic(fmt.Sprintf("Bot error : %+v", err))
	}
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Print("[DEBUG] Bot is ready")
	})
	return s
}

func sqlite() *sql.DB {
	db, err := sql.Open("sqlite3", viper.GetString("sqlite.file"))
	if err != nil {
		panic("Cannot open sql lite file")
	}
	return db
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
		log.Print("[ERROR] Marshal error")
		return nil
	}

	headers := map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7",
		"Connection":      "keep-alive",
		"Content-Type":    "application/json;charset=UTF-8",
	}

	cookies, _ := exec(http.MethodPost, LoginURL, bodyJSON, headers, nil)

	if len(cookies) == 0 {
		log.Print("[ERROR] Connection failed, no cookie in response body")
		return nil
	}

	_, resp := exec(http.MethodGet, fmt.Sprintf(InformationURL, viper.GetString("etna.user")), nil, nil, cookies[0])

	var information []*Information
	err = json.Unmarshal(resp, &information)
	if err != nil {
		log.Print("[ERROR] UnMarshal error")
		return nil
	}

	return information
}

func exec(method, url string, body []byte, headers map[string]string, cookie *http.Cookie) ([]*http.Cookie, []byte) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Print("[ERROR] Request error", err)
		return nil, nil
	}

	for key, header := range headers {
		req.Header.Add(key, header)
	}

	if cookie != nil {
		req.AddCookie(cookie)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Print("[ERROR] Request error", err)
		return nil, nil
	}
	defer res.Body.Close()

	bodyJSON, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("[ERROR] Request error", err)
		return nil, nil
	}
	fmt.Println(res.StatusCode, string(bodyJSON))

	return res.Cookies(), bodyJSON
}

func sendDiscordMessage(s *discordgo.Session, channel, message string) {
	messageSend, err := s.ChannelMessageSend(channel, message)
	if err != nil {
		return
	}
	log.Print("[DEBUG] Message sent at : ", messageSend.Timestamp)
}

func insertIntoDatabase(db *sql.DB, notif Notification) {
	_, err := db.Exec("INSERT INTO notification VALUES(NULL,datetime(),?, ?);", notif.ExternalID, notif.User)
	if err != nil {
		log.Print("[ERROR] Insert into database failed", err)
	}
}

func notificationExist(db *sql.DB, notif Notification) bool {
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

func openLogFile() *os.File {
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		err := os.Mkdir("log", os.ModePerm)
		if err != nil {
			return nil
		}
	}
	f, err := os.OpenFile("log/debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return f
}
