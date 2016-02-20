package ss13

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type Instance struct {
	// Settings required by the user.
	Debug bool
	DB    *DB

	// Internal stuff
	addr   string
	router *gin.Engine
}

// Database models

type Server struct {
	ID          int
	LastUpdated time.Time
	Title       string `sql:"type:varchar(64);unique"`
	GameUrl     string
	SiteUrl     string

	PlayersCurrent int
	PlayersAvg     int
	PlayersMin     int
	PlayersMax     int

	PlayersMon int
	PlayersTue int
	PlayersWed int
	PlayersThu int
	PlayersFri int
	PlayersSat int
	PlayersSun int
}

// Check if a server's last_updated time since now is greater or equal to X hours.
func (s *Server) TimeIsGreater(hours int) bool {
	return int(time.Since(s.LastUpdated)/time.Hour) >= hours
}

// Return a formatted string of LastUpdated.
func (s *Server) Timestamp() string {
	return s.LastUpdated.String()
}

// Return a fancy string duration since LastUpdated.
func (s *Server) TimeSince() string {
	d := time.Since(s.LastUpdated)
	mins := int(d / time.Minute % 60)
	hours := int(d / time.Hour % 24)
	days := int(d / time.Hour / 24)

	tmp := fmt.Sprintf("%d minutes", mins)
	if hours > 1 {
		tmp = fmt.Sprintf("%d hours, ", hours) + tmp
	} else if hours == 1 {
		tmp = fmt.Sprintf("%d hour, ", hours) + tmp
	}
	if days > 1 {
		tmp = fmt.Sprintf("%d days, ", days) + tmp
	} else if days == 1 {
		tmp = fmt.Sprintf("%d day, ", days) + tmp
	}
	return tmp
}

type ServerPopulation struct {
	ID        int
	Timestamp time.Time
	Players   int
	ServerID  int `sql:"index;type:integer REFERENCES servers(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	Server    Server
}

// See https://github.com/jinzhu/gorm/issues/635 for why we have to manually add
// in a raw REFERENCES statement here.
// Hint: Foreign key creation is bugged when using gorm with sqlite.
