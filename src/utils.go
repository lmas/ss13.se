package ss13

import (
	"log"
	"time"
)

var (
	now            = time.Now()
	debugging bool = false
)

func log_error(err error) {
	if err != nil {
		log.Panic("WARNING ", err)
	}
}

func check_error(err error) {
	if err != nil {
		log.Fatal("ERROR ", err)
	}
}

func SetDebug(val bool) {
	debugging = val
}

func IsDebugging() bool {
	return debugging
}

func ResetNow() {
	now = time.Now()
}

func Now() time.Time {
	return now
}
