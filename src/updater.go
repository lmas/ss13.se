package ss13

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

var updatedservers []string

const SERVERS_CONFIG = "./servers.json" // TODO

type RawServerData struct {
	Title     string
	Game_url  string
	Site_url  string
	Players   int
	Timestamp time.Time
}

func (i *Instance) UpdateServers() {
	reset()

	tx := i.DB.Begin()

	config, err := LoadConfig(SERVERS_CONFIG)
	if err != nil {
		fmt.Printf("Unable to load servers to poll: %s\n", err) // TODO
	} else {
		if i.Debug {
			fmt.Println("\nPolling servers...")
		}
		for _, s := range PollServers(config.PollServers, config.Timeout) {
			i.update_server(tx, s)
		}
	}

	if i.Debug {
		fmt.Println("\nScraping servers...")
	}
	for _, s := range ScrapePage() {
		i.update_server(tx, s)
	}

	if i.Debug {
		fmt.Println("\nRemoving old servers...")
	}
	RemoveOldServers(tx, Now())

	if i.Debug {
		fmt.Println("\nUpdating inactive servers...")
	}
	for _, s := range i.get_old_servers() {
		i.update_server(tx, s)
	}

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
	for _, old := range GetOldServers(i.DB, Now()) {
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

func (i *Instance) update_server(tx *gorm.DB, s *RawServerData) {
	if isupdated(s.Title) {
		return
	}

	if IsDebugging() {
		fmt.Println(s.Title)
	}

	// get server's db id (or create)
	id := InsertOrSelect(tx, s)

	// create new player history point
	AddServerPopulation(tx, id, s)

	// update server (urls and player stats)
	UpdateServerStats(tx, id, s)
}
