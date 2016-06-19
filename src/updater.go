package ss13

import (
	"fmt"
	"time"
)

var now = time.Now()

func Now() time.Time {
	return now.UTC()
}

type RawServerData struct {
	Title     string
	Game_url  string
	Site_url  string
	Players   int
	Timestamp time.Time
}

func (i *Instance) UpdateServers() {
	now = time.Now()
	servers := make(map[string]*RawServerData)
	addServer := func(s *RawServerData) {
		if _, exists := servers[s.Title]; exists {
			fmt.Println("Already existed:", s.Title)
			return
		}
		servers[s.Title] = s
	}

	fmt.Println("polling")
	polled := i.PollServers(i.Config.Servers, i.Config.UpdateTimeout)
	for _, s := range polled {
		addServer(s)
	}

	fmt.Println("scraping")
	scraped, e := i.ScrapePage()
	if e != nil {
		Log("Error scraping servers: %s", e)
	} else {
		for _, s := range scraped {
			addServer(s)
		}
	}

	fmt.Println("oldies")
	for _, s := range i.getOldServers() {
		addServer(s)
	}

	tx := i.db.NewTransaction()
	for _, s := range servers {
		fmt.Println("Updating:", s.Title)
		// get server's db id (or create)
		id := tx.InsertOrSelect(s)
		// create new player history point
		tx.AddServerPopulation(id, s, now)
		// update server (urls and player stats)
		tx.UpdateServerStats(id, s, now)
	}
	tx.RemoveOldServers(Now())
	tx.Commit()
}

func (i *Instance) getOldServers() []*RawServerData {
	// TODO: there's some bug that makes this func return all servers?
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
