package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(username, password, host, port, database string) (*Service, error) {
	connStr := "postgresql://" + username + ":" + password + "@" + host + ":" + port + "/" + database + "?sslmode=disable"
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to postgresql database : %+v", err)
		return nil, err
	}

	return &Service{DB: db}, nil
}
