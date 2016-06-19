package ss13

import "github.com/BurntSushi/toml"

type ServerConfig struct {
	Title   string
	GameUrl string
	SiteUrl string
}

type Config struct {
	// Path to sqlite database file.
	DatabasePath string

	// Serve web pages on this address.
	ListenAddr string

	// List of "private" servers to manually poll for updates (private as
	// in they do not show up on the byond hub page).
	Servers []ServerConfig

	// Update all servers every x minutes.
	UpdateEvery int

	// Timeout after x seconds, when trying to update a server.
	UpdateTimeout int
}

func LoadConfig(path string) (*Config, error) {
	c := Config{}
	_, e := toml.DecodeFile(path, &c)
	if e != nil {
		return nil, e
	}
	return &c, nil
}
