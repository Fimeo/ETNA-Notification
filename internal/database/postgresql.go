package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

func NewDatabaseConnection() Client {
	connection, err := postgresConnection(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	if err != nil {
		panic(err)
	}

	return Client{connection}
}

func postgresConnection(username, password, host, port, database string) (*gorm.DB, error) {
	dsn := "postgresql://" + username + ":" + password + "@" + host + ":" + port + "/" + database + "?sslmode=disable"
	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[ERROR] Failed to connect to postgresql database : %+v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("[ERROR]Cannot get sql db from gorm instance : %+v", err)
		return nil, err
	}
	err = sqlDB.Ping()
	if err != nil {
		// TODO : add cron on to ping, if connection failed, send discord message in special channel
		log.Printf("[ERROR] Ping database failed : %+v", err)
		return nil, err
	}

	return db, nil
}
