package main

import (
	"flag"
	"time"

	"github.com/lmas/ss13_se"
)

var (
	flagAddr = flag.String("addr", ":8000", "Adress and port to run the web server on")
	flagPath = flag.String("path", "servers.db", "File path to database")
)

func main() {
	flag.Parse()

	// TODO: load config from a toml file
	conf := ss13_se.Conf{
		WebAddr:       *flagAddr,
		ReadTimeout:   30 * time.Second,
		WriteTimeout:  30 * time.Second,
		ScrapeTimeout: 10 * time.Minute,
		Storage: &ss13_se.StorageSqlite{
			Path: *flagPath,
		},
	}
	app, err := ss13_se.New(conf)
	if err != nil {
		panic(err)
	}

	err = app.Run()
	if err != nil {
		panic(err)
	}
}
