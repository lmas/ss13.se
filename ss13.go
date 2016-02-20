package main

import (
	"fmt"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/lmas/ss13_se/src"
)

func main() {
	app := cli.NewApp()
	app.Version = ss13.VERSION
	app.Usage = "" // TODO
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Run in debug mode",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Show verbose messages",
		},
		cli.StringFlag{
			Name:  "addr",
			Usage: "Set the listen address for the web server",
			Value: ":8000",
		},
		cli.BoolFlag{
			Name:  "daemon",
			Usage: "Continuously run when using the update command",
		},
		cli.IntFlag{
			Name:  "timeout",
			Usage: "Time (in minutes) between each update in daemon mode",
			Value: 15,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "Run the web server",
			Action: run_server,
		},
		{
			Name:   "update",
			Usage:  "Update the server stats",
			Action: update_stats,
		},
	}

	app.Run(os.Args)
}

func run_server(c *cli.Context) {
	if c.GlobalBool("verbose") {
		fmt.Printf("Listening on %s.\n\n", c.GlobalString("addr"))
	}

	instance := &ss13.Instance{
		Debug: c.GlobalBool("debug"),
		DB:    ss13.OpenSqliteDB("new.db"), // TODO
	}
	instance.Init()
	instance.Serve(c.GlobalString("addr"))
}

func update_stats(c *cli.Context) {
	td := time.Duration(c.GlobalInt("timeout")) * time.Minute
	if c.GlobalBool("verbose") {
		if c.GlobalBool("daemon") {
			fmt.Printf("Running updates every %v minutes\n", c.GlobalInt("timeout"))
		} else {
			fmt.Println("Updating...")
		}
	}

	instance := &ss13.Instance{
		Debug: c.GlobalBool("debug"),
		DB:    ss13.OpenSqliteDB("new.db"), // TODO
	}
	instance.Init()

	for {
		start := time.Now()
		instance.UpdateServers()
		stop := time.Now()
		if c.GlobalBool("verbose") {
			fmt.Printf("Update completed in %v\n", stop.Sub(start))
		}

		if !c.GlobalBool("daemon") {
			return
		}
		fmt.Println("Next at ", time.Now().Add(td))
		time.Sleep(td)
	}
}
