package ss13_se

import (
	"html/template"
	"net/url"
	"time"
)

type ServerEntry struct {
	ID      string    `db:"id"`
	Title   string    `db:"title"`
	SiteURL string    `db:"site_url"`
	GameURL string    `db:"game_url"`
	Time    time.Time `db:"time"`
	Players int       `db:"players"`
}

func (e ServerEntry) IsZero() bool {
	return e.ID == ""
}

func (e ServerEntry) LastUpdated() string {
	return e.Time.Format("2006-01-02 15:04 MST")
}
func (e ServerEntry) ByondURL() template.URL {
	u, err := url.Parse(e.GameURL)
	if err != nil {
		return ""
	}

	if u.Scheme != "byond" {
		return ""
	}

	return template.URL(u.String())
}

type ServerPoint struct {
	Time     time.Time `db:"time"`
	ServerID string    `db:"server_id"`
	Players  int       `db:"players"`
}

func (p ServerPoint) IsZero() bool {
	return p.ServerID == "" && p.Time.IsZero()
}

type Storage interface {
	Open() error

	SaveServers([]ServerEntry) error
	GetServer(string) (ServerEntry, error)
	GetServers() ([]ServerEntry, error)
	RemoveServers([]ServerEntry) error

	SaveServerHistory([]ServerPoint) error
	GetServerHistory(int) ([]ServerPoint, error)
	GetSingleServerHistory(string, int) ([]ServerPoint, error)
}
