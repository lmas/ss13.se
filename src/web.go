package ss13

import (
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmas/ss13_se/src/assetstatic"
	"github.com/lmas/ss13_se/src/assettemplates"
)

func (i *Instance) Init() {
	i.DB.InitSchema()
}

func (i *Instance) Serve(addr string) error {
	i.addr = addr
	if i.Debug == false {
		gin.SetMode(gin.ReleaseMode)
	}

	// TODO: replace Default with New and use custom logger and stuff?
	i.router = gin.Default()
	i.router.NoRoute(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.HTML(http.StatusNotFound, "page_404.html", nil)
		}
	}())

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

	// Load templates
	tmpl := template.New("AllTemplates").Funcs(funcmap)
	tmplfiles, err := assettemplates.AssetDir("templates/")
	if err != nil {
		panic(err)
	}
	for p, b := range tmplfiles {
		name := filepath.Base(p)
		template.Must(tmpl.New(name).Parse(string(b)))
	}
	i.router.SetHTMLTemplate(tmpl)

	// Load static files
	staticfiles, e := assetstatic.AssetDir("static/")
	if e != nil {
		panic(e)
	}
	for p, _ := range staticfiles {
		ctype := mime.TypeByExtension(filepath.Ext(p))
		// Need to make a local copy of the var or else all files will
		// return the content of a single file (quirk with range).
		b := staticfiles[p]
		i.router.GET(fmt.Sprintf("/%s", p), func(c *gin.Context) {
			c.Data(http.StatusOK, ctype, b)
		})
	}

	// Setup all URLS
	i.router.GET("/", i.page_index)

	i.router.GET("/server/:server_id/*slug", i.page_server)
	i.router.GET("/server/:server_id", i.page_server)

	//i.router.GET("/stats", page_stats)
	i.router.GET("/about", i.page_about)

	i.router.GET("/r/ver", i.page_apollo)

	return i.router.Run(i.addr)
}

func (i *Instance) page_index(c *gin.Context) {
	servers := i.DB.AllServers()
	c.HTML(http.StatusOK, "page_index.html", gin.H{
		"pagetitle": "Index",
		"servers":   servers,
	})
}

func (i *Instance) page_about(c *gin.Context) {
	c.HTML(http.StatusOK, "page_about.html", nil)
}

func (i *Instance) page_server(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("server_id"), 10, 0)
	if err != nil {
		c.HTML(http.StatusNotFound, "page_404.html", nil)
		return
	}

	s, err := i.DB.GetServer(int(id))
	if err != nil {
		c.HTML(http.StatusNotFound, "page_404.html", nil)
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
	c.HTML(http.StatusOK, "page_server.html", gin.H{
		"pagetitle":    s.Title,
		"server":       s,
		"weekhistory":  i.DB.GetServerPopulation(int(id), time.Duration(7*24+12)*time.Hour),
		"monthhistory": i.DB.GetServerPopulation(int(id), time.Duration(31*24)*time.Hour),
		"weekdayavg":   weekdayavg,
	})
}

func (i *Instance) page_apollo(c *gin.Context) {
	// Go away, this it not an easter egg.
	c.Redirect(http.StatusFound, "byond://192.95.55.67:3333")
}
