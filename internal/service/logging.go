package service

import (
	"log"
	"os"
)

func openLogFile() *os.File {
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		err := os.Mkdir("log", os.ModePerm)
		if err != nil {
			return nil
		}
	}
	f, err := os.OpenFile("log/debug.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return f
}

func Logger() *os.File {
	f := openLogFile()
	log.SetOutput(f)

	return f
}
