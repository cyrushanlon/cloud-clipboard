package client

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

//username and password
//remember password
//other keyboard shortcut
//cloud clipboard stuff
//remote config location?

//Conf global config variable
var Conf Config

//Config holds application wide config variables
type Config struct {
	//public
	Username         string
	Password         string
	RememberPassword bool
	MulticastIP      string
}

//Save the config to file
func (s Config) Save() {
	f, err := os.Open("config.conf")
	if err != nil {
		//File doesnt exist
		log.Println("Creating config file.")
		f, err = os.Create("config.conf")
		if err != nil {
			log.Println("Couldn't create the file.")
			return
		}
	}

	enc := toml.NewEncoder(f)
	err = enc.Encode(s)
	if err != nil {
		log.Println("Couldn't save to the file.")
		return
	}
	fmt.Println("Saved config to config file.")
}

//Load the config from file
func (s *Config) Load() {
	_, err := toml.DecodeFile("config.conf", s)
	if err != nil {
		log.Println("Couldn't read from the file, using defaults.")
		Conf = defaultConfig()
		Conf.Save()
	} //else {} //uses the config from the file
}

//Delete removes the current config file from disk
func (s *Config) Delete() error {
	return os.Remove("config.conf")
}

func defaultConfig() Config {
	return Config{
		Username:         "cyrushanlon",
		Password:         "",
		RememberPassword: false,
		MulticastIP:      "232.49.101.200:6964",
	}
}