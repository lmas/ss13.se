package ss13

import (
	"html/template"
	"log"
	"time"
)

var (
	now = time.Now()

	funcmap = template.FuncMap{
		// safe_href let's us use URLs with custom protocols
		"safe_href": func(s string) template.HTMLAttr {
			return template.HTMLAttr(`href="` + s + `"`)
		},
		"inms": func(t time.Time) int64 {
			return t.Unix() * 1000
		},
		"year": func() int {
			return time.Now().Year()
		},
	}
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
