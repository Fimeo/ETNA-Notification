package domain

import "time"

type User struct {
	ID             int       `json:"-"`
	Time           time.Time `json:"-"`
	ChannelID      string    `json:"-"`
	DiscordAccount string    `json:"-"`
	Login          string    `json:"login"`
	Password       string    `json:"password"`
}
