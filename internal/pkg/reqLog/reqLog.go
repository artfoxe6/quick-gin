package reqLog

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Log struct {
	Dir string
}

// Write 按照日期分开存储日志，方便排查
func (l Log) Write(msg []byte) (n int, err error) {

	logPath := fmt.Sprintf("%s/req%s.log", l.Dir, time.Now().Format("20060102"))
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
