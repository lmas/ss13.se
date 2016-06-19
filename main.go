package main

import (
	"flag"
	"log"
	"time"

	"github.com/lmas/ss13_se/src"
)

var (
	fAddr           = flag.String("addr", ":8000", "address to listen on, for the web server")
	fDatabase       = flag.String("database", "./ss13.db", "database file")
	fDebug          = flag.Bool("debug", false, "run in debug mode")
	fPrivateServers = flag.String("servers", "./servers.json", "file with a list of private servers to poll")
	fTimeout        = flag.Int("timeout", 15, "time (in minutes) between each update")
	fVerbose        = flag.Bool("verbose", false, "show verbose messages")
)

func main() {
	flag.Parse()

	db, e := ss13.OpenSqliteDB(*fDatabase)
	CheckError(e)

	ins := &ss13.Instance{
		Debug:           *fDebug,
		DB:              db,
		PrivServersFile: *fPrivateServers,
	}
	ins.Init()

	td := time.Duration(*fTimeout) * time.Minute
	go func() {
		if *fVerbose {
			Log("Updating servers every %s", td)
		}
		for {
			start := time.Now()
			ins.UpdateServers()
			dur := time.Since(start)
			if *fVerbose {
				Log("Update completed in %s", dur)
			}
			time.Sleep(td)
		}
	}()

	if *fVerbose {
		Log("Listening on %s", *fAddr)
	}
	e = ins.Serve(*fAddr)
	CheckError(e)
}

func Log(f string, args ...interface{}) {
	log.Printf(f+"\n", args...)
}

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
