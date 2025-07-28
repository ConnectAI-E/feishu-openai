package utils

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"time"
)

type MyLogWriter struct {
}

func (writer MyLogWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().UTC().Format("2006-01-02T15:04:05.999Z") + string(bytes))
}

func CloseLogger(logger *lumberjack.Logger) {
	err := logger.Close()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("logger closed")
	}
}
