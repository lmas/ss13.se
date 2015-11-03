package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/lmas/ss13_se/src"
)

var (
	// General flags
	f_debug   = flag.Bool("debug", false, "Run in debug mode")
	f_verbose = flag.Bool("verbose", false, "Show verbose messages")

	// Flags for the run (webserver) command
	f_addr = flag.String("addr", "127.0.0.1:8000", "Server's listening address")

	// Flags for the update command
	f_daemon  = flag.Bool("daemon", false, "Continuously update")
	f_timeout = flag.Int64("timeout", 15, "Time (in minutes) between each update in daemon mode")
)

type Command struct {
	Name   string
	Desc   string
	DoFunc func()
}

var commands = []Command{
	// Main commands
	Command{"serve", "Run a local web server.", doserve},
	Command{"update", "Update the population stats.", doupdate},

	// Misc. commands
	Command{"version", "Show current version and exit.", doversion},
}

func main() {
	flag.Parse()
	ss13.SetDebug(*f_debug)

	command := flag.Arg(0)
	for _, a := range commands {
		if command == a.Name {
			a.DoFunc()
			return
		}
	}

	dohelp()
}

func dohelp() {
	fmt.Println(`SS13 population stats server.

Usage:
	TODO_NAME [flags] command [arguments]

Available commands:`)
	for _, a := range commands {
		fmt.Printf("\t%s\t%s\n", a.Name, a.Desc)
	}
	fmt.Println(`
Run TODO_NAME -h to get a list of command line flags.
`)
}

func doversion() {
	fmt.Println("TODO: Use a fancy name, add a banner and set a version.")
}

func doserve() {
	if *f_verbose {
		fmt.Printf("Listening on %s.\n\n", *f_addr)
	}

	instance := &ss13.Instance{
		Debug: *f_debug,
		DB:    ss13.OpenSqliteDB("new.db"), // TODO
	}
	instance.Init()
	instance.Serve(*f_addr)
}

func doupdate() {
	td := time.Duration(*f_timeout) * time.Minute
	if *f_verbose {
		if *f_daemon {
			fmt.Printf("Running updates every %v minutes\n", *f_timeout)
		} else {
			fmt.Println("Updating...")
		}
	}

	instance := &ss13.Instance{
		Debug: *f_debug,
		DB:    ss13.OpenSqliteDB("new.db"), // TODO
	}
	instance.Init()

	for {
		start := time.Now()
		instance.UpdateServers()
		stop := time.Now()
		if *f_verbose {
			fmt.Printf("Update completed in %v\n", stop.Sub(start))
		}

		if !*f_daemon {
			return
		}
		fmt.Println("Next at ", time.Now().Add(td))
		time.Sleep(td)
	}
}
