package ss13

import (
	"log"
	"time"
)

var (
	now = time.Now()
)

func LogError(err error) bool {
	if err != nil {
		log.Printf("WARNING: %s\n", err)
		return true
	}
	return false
}

func ResetNow() {
	now = time.Now()
}

func Now() time.Time {
	return now.UTC()
}
