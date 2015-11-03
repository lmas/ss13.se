package ss13

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (i *Instance) Init() {
	InitSchema(i.DB)

	if i.Debug == false {
		gin.SetMode(gin.ReleaseMode)
	}

	// TODO: replace Default with New and use custom logger and stuff
	i.router = gin.Default()

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

	i.router.GET("/", i.server_index)

	i.router.GET("/server/:server_id/*slug", i.server_detail)
	i.router.GET("/server/:server_id", i.server_detail)

	//i.router.GET("/stats", page_stats)
	//i.router.GET("/about", page_about)
}

func (i *Instance) Serve(addr string) error {
	i.addr = addr
	return i.router.Run(i.addr)
}

func (i *Instance) server_index(c *gin.Context) {
	servers := AllServers(i.DB)
	c.HTML(http.StatusOK, "server_index.html", gin.H{
		"pagetitle": "Index",
		"servers":   servers,
	})
}

func (i *Instance) server_detail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("server_id"), 10, 0)
	check_error(err)
	s, err := GetServer(i.DB, int(id))
	if err != nil {
		// TODO
		//c.HTML(http.StatusNotFound, "error_404.html", nil)
		c.String(http.StatusNotFound, "Server not found")
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
	c.HTML(http.StatusOK, "server_detail.html", gin.H{
		"pagetitle":    s.Title,
		"server":       s,
		"weekhistory":  GetServerPopulation(i.DB, int(id), time.Duration(7*24+12)*time.Hour),
		"monthhistory": GetServerPopulation(i.DB, int(id), time.Duration(31*24)*time.Hour),
		"weekdayavg":   weekdayavg,
	})
}
