package domain

import (
	"gorm.io/gorm"
)

const (
	StatusPending = "pending"
	StatusOpen    = "open"
	StatusClose   = "close"
)

type User struct {
	gorm.Model
	ChannelID      string `json:"-"`
	DiscordAccount string `json:"-"`
	Login          string `json:"login"`
	Password       string `json:"password"`
	Status         string `json:"status"`
}
