package main

import (
	"time"

	"github.com/lmas/ss13_se"
)

func main() {
	// TODO: load config from a toml file
	conf := ss13_se.Conf{
		WebAddr:       ":8000",
		ReadTimeout:   30 * time.Second,
		WriteTimeout:  30 * time.Second,
		ScrapeTimeout: 10 * time.Minute,
		Storage: &ss13_se.StorageSqlite{
			Path: "./tmp/servers.db",
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
