package utils

import "time"

import _ "time/tzdata"

func GetParisLocation() *time.Location {
	location, _ := time.LoadLocation("Europe/Paris")
	return location
}
