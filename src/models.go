package ss13

import (
	"fmt"
	"time"
)

// Database models

type Server struct {
	ID          int
	LastUpdated time.Time `db:"last_updated"`
	Title       string
	GameUrl     string `db:"game_url"`
	SiteUrl     string `db:"site_url"`

	PlayersCurrent int `db:"players_current"`
	PlayersAvg     int `db:"players_avg"`
	PlayersMin     int `db:"players_min"`
	PlayersMax     int `db:"players_max"`

	PlayersMon int `db:"players_mon"`
	PlayersTue int `db:"players_tue"`
	PlayersWed int `db:"players_wed"`
	PlayersThu int `db:"players_thu"`
	PlayersFri int `db:"players_fri"`
	PlayersSat int `db:"players_sat"`
	PlayersSun int `db:"players_sun"`
}

// Check if a server's last_updated time since now is greater or equal to X hours.
func (s *Server) TimeIsGreater(hours int) bool {
	return int(time.Since(s.LastUpdated)/time.Hour) >= hours
}

// Return a formatted string of LastUpdated.
func (s *Server) Timestamp() string {
	return s.LastUpdated.Format("2006-01-02 15:04 MST")
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
	ServerID  int `db:"server_id"`
}
