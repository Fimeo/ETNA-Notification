package service

import (
	"database/sql"
	"log"

	"etna-notification/internal/model/sqlite"
)

type SQLLiteService struct {
	DB *sql.DB
}

func SQLLiteConn(filepath string) SQLLiteService {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic("Cannot open sql lite file")
	}

	return SQLLiteService{DB: db}
}

func (s *SQLLiteService) CreateNotification(notification sqlite.Notification) {
	_, err := s.DB.Exec(
		"INSERT INTO notification VALUES(NULL,datetime(),?, ?);",
		notification.ExternalID, notification.User)
	if err != nil {
		log.Print("[ERROR] Insert into database failed", err)
	}
}

func (s *SQLLiteService) IsAlreadyNotified(notification sqlite.Notification) bool {
	row, err := s.DB.Query(
		"SELECT * FROM notification WHERE external_id=? and user=?",
		notification.ExternalID, notification.User)
	if err != nil {
		panic("Cannot read from database")
	}
	count := 0
	for row.Next() {
		count++
	}
	return count != 0
}
