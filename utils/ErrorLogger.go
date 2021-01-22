package utils

import (
	"os"
	"time"
)

func LogError(err string) {
	currentTime := time.Now()
	path := "./error-logs/" + currentTime.Format("01-02-2006") + ".log"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Create(path)
	}
	f, _ := os.OpenFile(path, os.O_APPEND, 0644)
	defer f.Close()
	f.WriteString(err + "\n")
}
