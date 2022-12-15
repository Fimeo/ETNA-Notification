package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Time           time.Time `json:"-"`
	ChannelID      string    `json:"-"`
	DiscordAccount string    `json:"-"`
	Login          string    `json:"login"`
	Password       string    `json:"password"`
	Status         string    `json:"status"`
}
