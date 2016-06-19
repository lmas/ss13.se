package ss13

import "time"

var updatedServers []string

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

	polled := i.PollServers(i.Config.Servers, i.Config.UpdateTimeout)
	for _, s := range polled {
		i.updateServer(tx, s)
	}

	scraped, e := i.ScrapePage()
	if e != nil {
		Log("Error scraping servers: %s", e)
	} else {
		for _, s := range scraped {
			i.updateServer(tx, s)
		}
	}

	for _, s := range i.getOldServers() {
		i.updateServer(tx, s)
	}

	tx.RemoveOldServers(Now())
	tx.Commit()
}

func reset() {
	// Have to reset some stuff between each update.
	updatedServers = *new([]string)
	ResetNow()
}

func isUpdated(title string) bool {
	// Prevent low pop. servers, with identical name as a high pop. server,
	// from fucking with another server's history.
	for _, t := range updatedServers {
		if title == t {
			return true
		}
	}
	updatedServers = append(updatedServers, title)
	return false
}

func (i *Instance) getOldServers() []*RawServerData {
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

func (i *Instance) updateServer(tx *DB, s *RawServerData) {
	if isUpdated(s.Title) {
		return
	}

	// get server's db id (or create)
	id := tx.InsertOrSelect(s)
	// create new player history point
	tx.AddServerPopulation(id, s)
	// update server (urls and player stats)
	tx.UpdateServerStats(id, s)
}
