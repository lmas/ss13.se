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

	// TODO: move this block into db.go
	//tx, e := i.db.NewTransaction()
	t, e := i.db.Beginx()
	if e != nil {
		panic(e) // TODO
	}
	tx := &TX{t}

	for _, s := range servers {
		fmt.Println("Updating:", s.Title)
		// get server's db id (or create)
		id, e := tx.InsertOrSelect(s)
		if e != nil {
			continue
		}
		// create new player history point
		e = tx.AddServerPopulation(id, s, now)
		if e != nil {
			continue
		}
		// update server (urls and player stats)
		e = tx.UpdateServerStats(id, s, now)
		if e != nil {
			continue
		}
	}
	e = tx.Commit()
	if e != nil {
		panic(e) // TODO
	}

	e = i.db.RemoveOldServers(Now())
	if e != nil {
		panic(e) // TODO
	}
}

func (i *Instance) getOldServers() []*RawServerData {
	var tmp []*RawServerData
	servers, e := i.db.GetOldServers(Now())
	if e != nil {
		panic(e) // TODO
	}
	for _, old := range servers {
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
