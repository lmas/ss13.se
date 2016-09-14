package ss13

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const SCHEMA string = `
CREATE TABLE IF NOT EXISTS servers (
	id INTEGER PRIMARY KEY,
	last_updated DATETIME,
	title TEXT UNIQUE,
	game_url TEXT,
	site_url TEXT,
	players_current INTEGER,
	players_avg INTEGER,
	players_min INTEGER,
	players_max INTEGER,
	players_mon INTEGER,
	players_tue INTEGER,
	players_wed INTEGER,
	players_thu INTEGER,
	players_fri INTEGER,
	players_sat INTEGER,
	players_sun INTEGER
);

CREATE TABLE IF NOT EXISTS server_populations (
	id INTEGER PRIMARY KEY,
	timestamp DATETIME,
	players INTEGER,
	server_id INTEGER REFERENCES servers(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS server_pop_index ON server_populations(server_id);
`

type Database struct {
	*sqlx.DB
}

type TX struct {
	*sqlx.Tx
}

func OpenDatabase(path string) (*Database, error) {
	db, e := sqlx.Connect("sqlite3", path)
	if e != nil {
		return nil, e
	}
	_, e = db.Exec(SCHEMA)
	if e != nil {
		return nil, e
	}
	return &Database{db}, nil
}

func (db *Database) AllServers() ([]*Server, error) {
	var tmp []*Server
	e := db.Select(&tmp, `SELECT * FROM servers ORDER BY
		last_updated DESC,
		players_current DESC,
		players_avg DESC,
		title;`)
	if e != nil {
		return nil, e
	}
	return tmp, nil
}

func (db *Database) GetServer(id int) (*Server, error) {
	var tmp Server
	e := db.Get(&tmp, "SELECT * FROM servers WHERE id = ? LIMIT 1;", id)
	if e != nil {
		return nil, e
	}
	return &tmp, nil
}

func (db *Database) GetOldServers(ts time.Time) ([]*Server, error) {
	var tmp []*Server
	e := db.Select(&tmp, "SELECT * FROM servers WHERE last_updated < ?;", ts)
	if e != nil {
		return nil, e
	}
	return tmp, nil
}

func (db *Database) RemoveOldServers(ts time.Time) error {
	_, e := db.Exec("DELETE FROM servers WHERE last_updated < datetime(?, '-7 days');", ts)
	return e
}

func (db *Database) GetServerPopulation(id int, d time.Duration) ([]*ServerPopulation, error) {
	var tmp []*ServerPopulation
	t := time.Now().Add(-d)
	e := db.Select(&tmp, `SELECT * FROM server_populations WHERE 
		server_id = ? AND timestamp > ?
		ORDER BY timestamp DESC, server_id;`, id, t)
	if e != nil {
		return nil, e
	}
	return tmp, nil
}

func (db *TX) InsertOrSelect(s *RawServerData) (int, error) {
	var tmp Server
	e := db.Get(&tmp, "SELECT * FROM servers WHERE title = ? LIMIT 1;", s.Title)
	if e == nil {
		return tmp.ID, nil
	}

	r, e := db.Exec(`INSERT INTO servers (
		last_updated, title, game_url, site_url, players_current,
		players_avg, players_min, players_max
		) VALUES(?, ?, ?, ?, ?, ?, ?, ?);`,
		s.Timestamp,
		s.Title,
		s.Game_url,
		s.Site_url,
		s.Players, s.Players, s.Players, s.Players)
	if e != nil {
		return -1, e
	}
	id, e := r.LastInsertId()
	if e != nil {
		return -1, e
	}
	return int(id), nil
}

func (db *TX) AddServerPopulation(id int, s *RawServerData, now time.Time) error {
	_, e := db.Exec(`INSERT INTO server_populations (
		timestamp, players, server_id
		) VALUES (?, ?, ?);`, now, s.Players, id)
	return e
}

func (db *TX) UpdateServerStats(id int, s *RawServerData, now time.Time) error {
	period := now.Add(-time.Duration(30*24) * time.Hour)
	rows, e := db.Queryx(`SELECT timestamp, players FROM server_populations
	WHERE server_id = ? AND timestamp > ?
	ORDER BY timestamp DESC`, id, period)
	if e != nil {
		return e
	}
	defer rows.Close()

	var timestamp time.Time
	var players, sum, count int
	var day_sums [7]int
	day_counts := [7]int{1, 1, 1, 1, 1, 1, 1} // Using 1's to prevent ZeroDiv error
	min := s.Players
	max := s.Players
	for rows.Next() {
		rows.Scan(&timestamp, &players)
		count++
		sum += players
		if players < min {
			min = players
		}
		if players > max {
			max = players
		}

		day := timestamp.Weekday()
		day_sums[day] += players
		day_counts[day]++
	}

	_, e = db.Exec(`UPDATE servers SET last_updated = ?, players_current = ?,
	players_avg = ?, players_min = ?, players_max = ?, players_mon = ?,
	players_tue = ?, players_wed = ?, players_thu = ?, players_fri = ?,
	players_sat = ?, players_sun = ?
	WHERE id = ?;`, s.Timestamp, s.Players, sum/count, min, max,
		day_sums[1]/day_counts[1],
		day_sums[2]/day_counts[2],
		day_sums[3]/day_counts[3],
		day_sums[4]/day_counts[4],
		day_sums[5]/day_counts[5],
		day_sums[6]/day_counts[6],
		day_sums[0]/day_counts[0],
		id)
	return e
}
