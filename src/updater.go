package ss13

import (
	"fmt"
	"time"
)

var updatedservers []string

type RawServerData struct {
	Title     string
	Game_url  string
	Site_url  string
	Players   int
	Timestamp time.Time
}

func (i *Instance) UpdateServers() {
	reset()

	tx := i.db.NewTransaction()

	if i.Debug {
		fmt.Println("\nPolling servers...")
	}
	polled, err := i.PollServers(i.Config.Servers, i.Config.UpdateTimeout)
	if !LogError(err) {
		for _, s := range polled {
			i.update_server(tx, s)
		}
	}

	if i.Debug {
		fmt.Println("\nScraping servers...")
	}
	scraped, err := i.ScrapePage()
	if !LogError(err) {
		for _, s := range scraped {
			i.update_server(tx, s)
		}
	}

	if i.Debug {
		fmt.Println("\nUpdating inactive servers...")
	}
	for _, s := range i.get_old_servers() {
		i.update_server(tx, s)
	}

	if i.Debug {
		fmt.Println("\nRemoving old servers...")
	}
	tx.RemoveOldServers(Now())

	tx.Commit()
}

func reset() {
	// If the updater is running in daemon mode we have to reset some stuff
	// each time we try to run a new update.
	updatedservers = *new([]string)
	ResetNow()
}

func isupdated(title string) bool {
	// Prevent low pop. servers, with identical name as a high pop. server,
	// from fucking with another server's history.
	for _, t := range updatedservers {
		if title == t {
			return true
		}
	}
	updatedservers = append(updatedservers, title)
	return false
}

func (i *Instance) get_old_servers() []*RawServerData {
	var tmp []*RawServerData
	for _, old := range i.db.GetOldServers(Now()) {
		s := RawServerData{
			Title:     old.Title,
			Game_url:  old.GameUrl,
			Site_url:  old.SiteUrl,
			Players:   0,
			Timestamp: old.LastUpdated,
		}
		tmp = append(tmp, &s)
	}
	return tmp
}

func (i *Instance) update_server(tx *DB, s *RawServerData) {
	if isupdated(s.Title) {
		return
	}

	if i.Debug {
		fmt.Println(s.Title)
	}

	// get server's db id (or create)
	id := tx.InsertOrSelect(s)

	// create new player history point
	tx.AddServerPopulation(id, s)

	// update server (urls and player stats)
	tx.UpdateServerStats(id, s)
}
