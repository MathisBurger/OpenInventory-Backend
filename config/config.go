package config

import (
	"encoding/json"
	"os"
)

// Database parameter of the configuration
// Contains all information to connect to database
type database struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Host     string `json:"host"`
}

// Configuration settings for the webserver
// Contains also web endpoints and logging settings
type serverconfig struct {
	Prefork            bool              `json:"prefork"`
	Port               string            `json:"port"`
	LogRequestErrors   bool              `json:"logRequestErrors"`
	WebEndpoints       map[string]string `json:"web-endpoints"`
	WebIndexEntrypoint string            `json:"web-index-entrypoint"`
}

// Main configuration struct
// Contains database and serverconfig
type Config struct {
	name      string       `json:"name"`
	Db        database     `json:"database"`
	ServerCFG serverconfig `json:"server-config"`
}

//////////////////////////////////////////////////
//                  ParseConfig                 //
//   Function to parse the configuration file   //
//   returns an Pointer of the Config struct    //
//////////////////////////////////////////////////
func ParseConfig() (c *Config, err error) {
	f, err := os.Open("./config/config.json")
	if err != nil {

		return
	}
	c = new(Config)
	err = json.NewDecoder(f).Decode(c)
	return
}
