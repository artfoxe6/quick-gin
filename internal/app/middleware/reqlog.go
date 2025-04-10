package middleware

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Log struct {
	Dir string
}

func (l Log) Write(msg []byte) (n int, err error) {

	logPath := fmt.Sprintf("%s/req%s.log", l.Dir, time.Now().Format("200601"))

	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	_, err = file.WriteString(string(msg))
	if err != nil {
		log.Fatalln(err)
	}
	return len(msg), nil
}
