package main

import (
	"time"

	"github.com/lmas/ss13hub"
)

func main() {
	// TODO: load config from a toml file
	conf := ss13hub.Conf{
		WebAddr:       ":8000",
		ReadTimeout:   30 * time.Second,
		WriteTimeout:  30 * time.Second,
		ScrapeTimeout: 10 * time.Minute,
		Storage: &ss13hub.StorageSqlite{
			Path: "./tmp/servers.db",
		},
	}
	app, err := ss13hub.New(conf)
	if err != nil {
		panic(err)
	}

	err = app.Run()
	if err != nil {
		panic(err)
	}
}
