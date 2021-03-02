package utils

import (
	"os"
	"strconv"
	"time"
)

// loggs errors
func LogError(err string, file string, line int) {

	// define log path
	currentTime := time.Now()
	path := "./error-logs/" + currentTime.Format("01-02-2006") + ".log"

	// check existance
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Create(path)
	}

	// write to file
	f, _ := os.OpenFile(path, os.O_APPEND, 0644)
	defer f.Close()
	f.WriteString("[" + file + ", " + strconv.Itoa(line) + "] " + err + "\n")
}
