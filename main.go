package ss13_se

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Used internally for logging a global # of players
const internalServerTitle string = "_ss13.se"

type Conf struct {
	// Web stuff
	WebAddr      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// Scraper stuff
	ScrapeTimeout time.Duration

	// Misc.
	Storage Storage
}

type App struct {
	conf      Conf
	web       *http.Server
	store     Storage
	templates map[string]*template.Template
}

func New(c Conf) (*App, error) {
	templates, err := loadTemplates()
	if err != nil {
		return nil, err
	}

	w := &http.Server{
		Addr:         c.WebAddr,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
	}

	a := &App{
		conf:      c,
		web:       w,
		store:     c.Storage,
		templates: templates,
	}

	r := mux.NewRouter()
	r.Handle("/", handler(a.pageIndex))
	r.Handle("/server/{id}", handler(a.pageServer))
	r.Handle("/server/{id}/daily", handler(a.pageDailyChart))
	r.Handle("/server/{id}/weekly", handler(a.pageWeeklyChart))
	r.Handle("/server/{id}/average", handler(a.pageAverageChart))
	a.web.Handler = r

	return a, nil
}

func (a *App) Log(msg string, args ...interface{}) {
	log.Printf(msg+"\n", args...)
}

func (a *App) Run() error {
	a.Log("Opening storage...")
	err := a.store.Open()
	if err != nil {
		return err
	}

	a.Log("Running updater")
	go a.runUpdater()

	a.Log("Running server on %s", a.conf.WebAddr)
	return a.web.ListenAndServe()
}

func (a *App) runUpdater() {
	webClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	for {
		now := time.Now()
		servers, err := scrapeByond(webClient, now)
		dur := time.Since(now)
		a.Log("Scrape done in %s, errors: %v", dur, err)

		if err == nil {
			servers = append(servers, a.makeHubEntry(now, servers))

			if err := a.store.SaveServers(servers); err != nil {
				a.Log("Error saving servers: %s", err)
			}

			if err := a.updateHistory(now, servers); err != nil {
				a.Log("Error saving server history: %s", err)
			}

			if err := a.updateOldServers(now, servers); err != nil {
				a.Log("Error updating old servers: %s", err)
			}
		}

		time.Sleep(a.conf.ScrapeTimeout)
	}
}

func (a *App) updateHistory(t time.Time, servers []ServerEntry) error {
	var history []ServerPoint
	for _, s := range servers {
		history = append(history, ServerPoint{
			Time:     t,
			ServerID: s.ID,
			Players:  s.Players,
		})
	}
	return a.store.SaveServerHistory(history)
}

func (a *App) updateOldServers(t time.Time, servers []ServerEntry) error {
	var old []ServerEntry
	for _, s := range servers {
		if !s.Time.Equal(t) {
			s.Players = 0
			old = append(old, s)
		}
	}

	var remove []ServerEntry
	for index, s := range old {
		delta := t.Sub(s.Time)
		if delta.Hours() > 24*1 { // TODO: CHANGE TO A WEEK AFTER TESTING
			remove = append(remove, s)
			old = append(old[:index], old[index+1:]...)
		}
	}

	if len(remove) > 0 {
		a.Log("Removing servers: %s", serverNameList(remove)) // TODO: remove after testing?
		if err := a.store.RemoveServers(remove); err != nil {
			return err
		}
	}

	if len(old) > 0 {
		a.Log("Old servers: %s", serverNameList(old)) // TODO: remove after testing
		if err := a.updateHistory(t, old); err != nil {
			return err
		}
	}
	return nil
}

// TODO: can probably remove this func after we're done testing
func serverNameList(servers []ServerEntry) string {
	var names []string
	for _, s := range servers {
		names = append(names, s.Title)
	}
	return strings.Join(names, ", ")
}

func (a *App) makeHubEntry(t time.Time, servers []ServerEntry) ServerEntry {
	var totalPlayers int
	for _, s := range servers {
		totalPlayers += s.Players
	}

	return ServerEntry{
		ID:      makeID(internalServerTitle),
		Title:   internalServerTitle,
		SiteURL: "",
		GameURL: "",
		Time:    t,
		Players: totalPlayers,
	}
}
