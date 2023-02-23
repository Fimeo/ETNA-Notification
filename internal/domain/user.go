package domain

import (
	"gorm.io/gorm"
	"net/http"
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
	authentication *http.Cookie
}

func (u *User) SetAuthentication(cookie *http.Cookie) {
	u.authentication = cookie
}

func (u *User) HasValidAuthentication() bool {
	if u.authentication == nil {
		return false
	}

	return u.authentication.Valid() == nil
}

func (u *User) GetAuthentication() *http.Cookie {
	return u.authentication
}
