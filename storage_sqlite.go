package ss13_se

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const sqliteScheme string = `
CREATE TABLE IF NOT EXISTS server_entry(
	id TEXT UNIQUE,
	title STRING,
	site_url STRING,
	game_url STRING,
	time DATETIME,
	players INTEGER
);

CREATE INDEX IF NOT EXISTS idx_server_entry ON server_entry(time, players, title);

CREATE TABLE IF NOT EXISTS server_history (
	id INTEGER PRIMARY KEY,
	time DATETIME,
	server_id TEXT,
	players INTEGER
);

CREATE INDEX IF NOT EXISTS idx_server_history ON server_history(time, server_id);
`

type StorageSqlite struct {
	*sqlx.DB
	Path string
}

func (store *StorageSqlite) Open() error {
	db, err := sqlx.Connect("sqlite3", store.Path)
	if err != nil {
		return err
	}

	_, err = db.Exec(sqliteScheme)
	if err != nil {
		return err
	}

	store.DB = db
	return nil
}

func (store *StorageSqlite) SaveServers(servers []ServerEntry) error {
	tx, err := store.Begin()
	if err != nil {
		return err
	}

	q := `INSERT OR REPLACE INTO server_entry (id, title, site_url, game_url, time, players) VALUES(?, ?, ?, ?, ?, ?);`
	for _, s := range servers {
		_, err := tx.Exec(q, s.ID, s.Title, s.SiteURL, s.GameURL, s.Time, s.Players)
		if err != nil {
			tx.Rollback() // TODO: handle error?
			return err
		}
	}

	return tx.Commit()
}

func (store *StorageSqlite) GetServer(id string) (ServerEntry, error) {
	var server ServerEntry
	q := `SELECT * FROM server_entry WHERE id = ? LIMIT 1;`
	err := store.Get(&server, q, id)
	if err != nil {
		return ServerEntry{}, err
	}
	return server, nil
}

func (store *StorageSqlite) GetServers() ([]ServerEntry, error) {
	var servers []ServerEntry
	q := `SELECT * FROM server_entry ORDER BY players DESC, id ASC;`
	err := store.Select(&servers, q)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (store *StorageSqlite) RemoveServers(servers []ServerEntry) error {
	tx, err := store.Begin()
	if err != nil {
		return err
	}

	qHistory := `DELETE FROM server_history WHERE server_id = ?;`
	qEntry := `DELETE FROM server_entry WHERE id = ?;`
	for _, s := range servers {
		_, err := tx.Exec(qHistory, s.ID)
		if err != nil {
			tx.Rollback() // TODO: handle error?
			return err
		}

		_, err = tx.Exec(qEntry, s.ID)
		if err != nil {
			tx.Rollback() // TODO: handle error?
			return err
		}
	}

	return tx.Commit()
}

func (store *StorageSqlite) SaveServerHistory(points []ServerPoint) error {
	tx, err := store.Begin()
	if err != nil {
		return err
	}

	q := `INSERT INTO server_history (time, server_id, players) VALUES(?, ?, ?);`
	for _, p := range points {
		_, err := tx.Exec(q, p.Time, p.ServerID, p.Players)
		if err != nil {
			tx.Rollback() // TODO: handle error?
			return err
		}
	}

	return tx.Commit()
}

func (store *StorageSqlite) GetServerHistory(days int) ([]ServerPoint, error) {
	var points []ServerPoint
	delta := time.Now().AddDate(0, 0, -days)
	q := `SELECT time,server_id,players FROM server_history WHERE time > ? ORDER BY time DESC, server_id ASC;`
	err := store.Select(&points, q, delta)
	if err != nil {
		return nil, err
	}
	return points, nil
}

func (store *StorageSqlite) GetSingleServerHistory(id string, days int) ([]ServerPoint, error) {
	var points []ServerPoint
	delta := time.Now().AddDate(0, 0, -days)
	q := `SELECT time,server_id,players FROM server_history WHERE server_id = ? AND time > ? ORDER BY time DESC;`
	err := store.Select(&points, q, id, delta)
	if err != nil {
		return nil, err
	}
	return points, nil
}
