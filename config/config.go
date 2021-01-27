package config

import (
	"encoding/json"
	"os"
)

type database struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Host     string `json:"host"`
}

type serverconfig struct {
	Prefork            bool              `json:"prefork"`
	Port               string            `json:"port"`
	WebEndpoints       map[string]string `json:"web-endpoints"`
	WebIndexEntrypoint string            `json:"web-index-entrypoint"`
}

type Config struct {
	name      string       `json:"name"`
	Db        database     `json:"database"`
	ServerCFG serverconfig `json:"server-config"`
}

func ParseConfig() (c *Config, err error) {
	f, err := os.Open("./config/config.json")
	if err != nil {

		return
	}
	c = new(Config)
	err = json.NewDecoder(f).Decode(c)
	return
}
