package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// TODO: change these to be cmd args instead!
const db_file = "/home/lmas/projects/ss13_se/src/db.sqlite3"

// Dir to save new graphs in
const save_dir = "/home/lmas/projects/ss13_se/src/static/graphs"

// How far back in time the graphs will go
var last_week = time.Now().AddDate(0, 0, -7)
var last_month = time.Now().AddDate(0, -1, 0)

var week_days = [7]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

func checkerror(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// open a db connection
	db, err := sql.Open("sqlite3", db_file)
	checkerror(err)
	defer db.Close()

	// loop over each server in db
	rows, err := db.Query("select id, title from gameservers_server")
	checkerror(err)
	defer rows.Close()

	var (
		id    int
		title string
	)
	for rows.Next() {
		err := rows.Scan(&id, &title)
		checkerror(err)
		creategraphs(db, "week-", id, title)
		createweekdaygraph(db, "avg_days-", id, title)
	}
	err = rows.Err()
	checkerror(err)
}

func createtempfile(prefix string) (f *os.File, name string) {
	// create a tmp file
	file, err := ioutil.TempFile("", prefix)
	checkerror(err)
	return file, file.Name()
}

func getgraphpath(title string, prefix string) (path string) {
	err := os.MkdirAll(save_dir, 0777)
	checkerror(err)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(title)))
	return filepath.Join(save_dir, fmt.Sprintf("%s%s", prefix, hash))
}

func creategraphs(db *sql.DB, prefix string, id int, title string) {
	// create a tmp file
	ifile, ifilename := createtempfile(prefix)
	defer ifile.Close()

	// Make sure we have somewhere to save the stored graphs in
	ofilename := getgraphpath(title, prefix)

	// get the server's data and write it to the file
	rows, err := db.Query("select created,players from gameservers_serverhistory where server_id = ? and created >= ? order by created asc", id, last_week)
	checkerror(err)
	defer rows.Close()

	var (
		created time.Time
		players int
	)
	for rows.Next() {
		err := rows.Scan(&created, &players)
		checkerror(err)
		_, err = ifile.WriteString(fmt.Sprintf("%d, %d\n", created.Unix(), players))
		checkerror(err)
	}
	err = rows.Err()
	checkerror(err)

	// run the plotter against the data file
	err = exec.Command("./plot_time.sh", ifilename, ofilename).Run()
	checkerror(err)

	// close and remove the tmp file
	ifile.Close()
	os.Remove(ifilename)
}

func createweekdaygraph(db *sql.DB, prefix string, id int, title string) {
	// create a tmp file
	ifile, ifilename := createtempfile(prefix)
	defer ifile.Close()

	// Make sure we have somewhere to save the stored graphs in
	ofilename := getgraphpath(title, prefix)

	// get the server's data and write it to the file
	// TODO: Move sunday (first day in list at 0) to the end...
	rows, err := db.Query("select strftime('%w', created) as weekday, avg(players) from gameservers_serverhistory where server_id = ? and created >= ? group by weekday;", id, last_week)
	checkerror(err)
	defer rows.Close()

	var (
		day     int
		players float64
	)
	for rows.Next() {
		err := rows.Scan(&day, &players)
		checkerror(err)
		_, err = ifile.WriteString(fmt.Sprintf("%s, %f\n", week_days[day], players))
		checkerror(err)
	}
	err = rows.Err()
	checkerror(err)

	// run the plotter against the data file
	err = exec.Command("./plot_bar.sh", ifilename, ofilename).Run()
	checkerror(err)

	// close and remove the tmp file
	ifile.Close()
	os.Remove(ifilename)
}
