package ss13

import (
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func New(debug bool, path string) (*Instance, error) {
	c, e := LoadConfig(path)
	if e != nil {
		return nil, e
	}

	db, e := OpenSqliteDB(c.DatabasePath)
	if e != nil {
		return nil, e
	}
	db.InitSchema()

	i := Instance{
		Debug:  debug,
		DB:     db,
		Config: c,
	}

	return &i, nil
}

func (i *Instance) Run() error {
	go func() {
		td := time.Duration(i.Config.UpdateEvery) * time.Minute
		for {
			start := time.Now()
			i.UpdateServers()
			dur := time.Since(start)
			if i.Debug {
				fmt.Printf("Update completed in %s\n", dur)
			}
			time.Sleep(td)
		}
	}()

	i.router = mux.NewRouter().StrictSlash(true)
	i.router.NotFoundHandler = http.HandlerFunc(i.page_404)

	// Custom template functions
	funcmap := template.FuncMap{
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

	// WHen in debug mode we load the assets from disk instead of the
	// embedded ones.
	SetRawAssets(i.Debug)

	// Load templates
	tmpl := template.New("AllTemplates").Funcs(funcmap)
	tmplfiles, err := AssetDir("templates/")
	if err != nil {
		panic(err)
	}
	for p, b := range tmplfiles {
		name := filepath.Base(p)
		template.Must(tmpl.New(name).Parse(string(b)))
	}
	i.tmpls = tmpl

	// Load static files
	staticfiles, e := AssetDir("static/")
	if e != nil {
		panic(e)
	}
	for p, _ := range staticfiles {
		// Need to make a local copy of the var or else all files will
		// return the content of a single file (quirk with range).
		b := staticfiles[p]
		ctype := mime.TypeByExtension(filepath.Ext(p))
		i.router.HandleFunc(fmt.Sprintf("/%s", p),
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", ctype)
				_, e := w.Write(b)
				LogError(e)
			})
	}

	// Setup all URLS
	i.router.HandleFunc("/", i.page_index)
	i.router.HandleFunc("/about", i.page_about)
	i.router.HandleFunc("/r/ver", i.page_apollo)
	i.router.HandleFunc("/server/{id}", i.page_server)
	i.router.HandleFunc("/server/{id}/{slug}", i.page_server)

	return http.ListenAndServe(i.Config.ListenAddr, i.router)
}

func (i *Instance) page_404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	i.tmpls.ExecuteTemplate(w, "page_404.html", nil)
}

func (i *Instance) page_index(w http.ResponseWriter, r *http.Request) {
	servers := i.DB.AllServers()
	i.tmpls.ExecuteTemplate(w, "page_index.html", D{
		"pagetitle": "Index",
		"servers":   servers,
	})
}

func (i *Instance) page_about(w http.ResponseWriter, r *http.Request) {
	i.tmpls.ExecuteTemplate(w, "page_about.html", nil)
}

func (i *Instance) page_server(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 0)
	if err != nil {
		i.page_404(w, r)
		return
	}

	s, err := i.DB.GetServer(int(id))
	if err != nil {
		i.page_404(w, r)
		return
	}
	type weekday struct {
		Day     string
		Players int
	}
	weekdayavg := [7]weekday{
		weekday{"Monday", s.PlayersMon},
		weekday{"Tuesday", s.PlayersTue},
		weekday{"Wednessday", s.PlayersWed},
		weekday{"Thursday", s.PlayersThu},
		weekday{"Friday", s.PlayersFri},
		weekday{"Saturday", s.PlayersSat},
		weekday{"Sunday", s.PlayersSun},
	}
	i.tmpls.ExecuteTemplate(w, "page_server.html", D{
		"pagetitle":    s.Title,
		"server":       s,
		"weekhistory":  i.DB.GetServerPopulation(int(id), time.Duration(7*24+12)*time.Hour),
		"monthhistory": i.DB.GetServerPopulation(int(id), time.Duration(31*24)*time.Hour),
		"weekdayavg":   weekdayavg,
	})
}

func (i *Instance) page_apollo(w http.ResponseWriter, r *http.Request) {
	// Go away, this it not an easter egg.
	http.Redirect(w, r, "byond://192.95.55.67:3333", http.StatusFound)
}
