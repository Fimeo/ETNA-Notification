package mysql

import (
	"database/sql"
	"log"
	"time"
)

type SQLLiteService struct {
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

func NewConnection(filepath string) (*SQLLiteService, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Printf("[ERROR] Cannot load sqlite file : %+v", err)
		return nil, err
	}

	return &SQLLiteService{DB: db}, nil
}

func (s *SQLLiteService) CreateNotification(notification Notification) error {
	_, err := s.DB.Exec(
		"INSERT INTO notification VALUES(NULL,datetime(),?, ?);",
		notification.ExternalID, notification.User)
	if err != nil {
		log.Print("[ERROR] Insert into database failed", err)
		return err
	}

	return nil
}

func (s *SQLLiteService) IsAlreadyNotified(notification Notification) (bool, error) {
	row, err := s.DB.Query(
		"SELECT * FROM notification WHERE external_id=? and user=?",
		notification.ExternalID, notification.User)
	if err != nil {
		log.Print("[ERROR] Cannot read database", err)
		return false, nil
	}
	count := 0
	for row.Next() {
		count++
	}
	return count != 0, nil
}

func (s *SQLLiteService) GetEtnaUsers() ([]EtnaUser, error) {
	rows, err := s.DB.Query("SELECT * FROM users;")
	if err != nil {
		log.Print("[ERROR] Retrieve users failed in database", err)
		return nil, err
	}

	var got []EtnaUser
	for rows.Next() {
		var r EtnaUser
		err = rows.Scan(&r.ID, &r.Time, &r.UserID, &r.ChannelID, &r.Login, &r.Password)
		if err != nil {
			return nil, err
		}
		got = append(got, r)
	}

	return got, nil
}

func (s *SQLLiteService) Close() {
	s.DB.Close()
}
