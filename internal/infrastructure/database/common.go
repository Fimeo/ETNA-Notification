package database

import (
	"database/sql"
	"time"
)

type Service struct {
	DB *sql.DB
}

type Notification struct {
	ID         int
	ExternalID int
	User       string
}

type EtnaUser struct {
	ID        int
	UserID    int
	Time      time.Time
	ChannelID string
	Login     string
	Password  string
}

func (s *Service) Close() {
	s.DB.Close()
}
