package logging

import (
	"log"
	"os"
)

type loggerService struct {
	f *os.File
}

type Logger interface {
	CloseFile()
}

func InitLogger() Logger {
	f := openLogFile()
	log.SetOutput(f)

	return &loggerService{f: f}
}

func (l loggerService) CloseFile() {
	l.f.Close()
}

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
