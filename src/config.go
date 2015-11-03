package ss13

import (
	"encoding/json"
	"io/ioutil"
)

type ServerConfig struct {
	Title   string
	GameUrl string
	SiteUrl string
}

type Config struct {
	PollServers []ServerConfig
	Timeout     int
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tmp := &Config{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, err
	}

	return tmp, nil
}
