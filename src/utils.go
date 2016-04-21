package ss13

import (
	"log"
	"time"
)

var (
	now = time.Now()
)

func log_error(err error) bool {
	if err != nil {
		log.Printf("WARNING: %s\n", err)
		return true
	}
	return false
}

func check_error(err error) {
	if err != nil {
		log.Fatal("ERROR ", err)
	}
}

func ResetNow() {
	now = time.Now()
}

func Now() time.Time {
	return now.UTC()
}
