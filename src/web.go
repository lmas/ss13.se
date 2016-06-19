package ss13

import (
	"html/template"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func New(debug bool, path string) (*Instance, error) {
	// WHen in debug mode we load the assets from disk instead of the
	// embedded ones.
	SetRawAssets(debug)

	tmpl := template.New("AllTemplates").Funcs(funcmap)
	tmplfiles, err := AssetDir("templates/")
	if err != nil {
		panic(err)
	}
	for p, b := range tmplfiles {
		name := filepath.Base(p)
		template.Must(tmpl.New(name).Parse(string(b)))
	}

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
		Config: c,
		Debug:  debug,

		db:     db,
		router: mux.NewRouter().StrictSlash(true),
		tmpls:  tmpl,
	}

	i.router.HandleFunc("/", i.page_index)
	i.router.HandleFunc("/about", i.page_about)
	i.router.HandleFunc("/r/ver", i.page_apollo)
	i.router.HandleFunc("/server/{id}", i.page_server)
	i.router.HandleFunc("/server/{id}/{slug}", i.page_server)
	i.router.HandleFunc("/static/{file:.*}", i.page_static)
	i.router.NotFoundHandler = http.HandlerFunc(i.page_404)

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
				Log("Update completed in %s", dur)
			}
			time.Sleep(td)
		}
	}()
	return http.ListenAndServe(i.Config.ListenAddr, i.router)
}

func (i *Instance) page_404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	i.tmpls.ExecuteTemplate(w, "page_404.html", nil)
}

func (i *Instance) page_static(w http.ResponseWriter, r *http.Request) {
	p := path.Join("static", mux.Vars(r)["file"])
	b, e := Asset(p)
	if e != nil {
		i.page_404(w, r)
		return
	}

	ctype := mime.TypeByExtension(filepath.Ext(p))
	w.Header().Add("Content-Type", ctype)
	_, e = w.Write(b)
	if e != nil {
		Log("Error sending static file %s: %s", p, e)
	}
}

func (i *Instance) page_index(w http.ResponseWriter, r *http.Request) {
	servers := i.db.AllServers()
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

	s, err := i.db.GetServer(int(id))
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
		"weekhistory":  i.db.GetServerPopulation(int(id), time.Duration(7*24+12)*time.Hour),
		"monthhistory": i.db.GetServerPopulation(int(id), time.Duration(31*24)*time.Hour),
		"weekdayavg":   weekdayavg,
	})
}

func (i *Instance) page_apollo(w http.ResponseWriter, r *http.Request) {
	// Go away, this it not an easter egg.
	http.Redirect(w, r, "byond://192.95.55.67:3333", http.StatusFound)
}
