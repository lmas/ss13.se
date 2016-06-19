package main

import (
	"flag"
	"log"

	"github.com/lmas/ss13_se/src"
)

var (
	fConfig = flag.String("config", "config.toml", "path to config file")
	fDebug  = flag.Bool("debug", false, "run in debug mode")
)

func main() {
	flag.Parse()

	ins, e := ss13.New(*fDebug, *fConfig)
	CheckError(e)

	if *fDebug {
		Log("Updating servers every %d minutes", ins.Config.UpdateEvery)
		Log("Listening on %s", ins.Config.ListenAddr)
	}
	e = ins.Run()
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
