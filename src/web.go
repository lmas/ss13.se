package ss13

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

	// Load templates
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
	tmpl := template.Must(template.New("ServerTemplates").Funcs(funcmap).ParseGlob("templates/*"))
	i.router.SetHTMLTemplate(tmpl)

	// Setup all URLS
	i.router.Static("/static", "./static")

	i.router.GET("/", i.page_index)

	i.router.GET("/server/:server_id/*slug", i.page_server)
	i.router.GET("/server/:server_id", i.page_server)

	//i.router.GET("/stats", page_stats)
	i.router.GET("/about", i.page_about)

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
