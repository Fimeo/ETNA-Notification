package utils

import "time"

func GetParisLocation() *time.Location {
	location, _ := time.LoadLocation("Europe/Paris")
	return location
}
