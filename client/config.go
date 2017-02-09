package client

import (
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
	//private
	filePath string
	//public
	Username         string
	Password         string
	RememberPassword bool
}

//Save the config to file
func (s Config) Save() {
	f, err := os.Open(s.filePath)
	if err != nil {
		log.Println("Couldn't save to the file.")
	}

	enc := toml.NewEncoder(f)
	err = enc.Encode(s)
	if err != nil {
		log.Println("Couldn't save to the file.")
	}
}

//Load the config from file
func (s *Config) Load() {
	_, err := toml.DecodeFile(s.filePath, s)
	if err != nil {
		log.Println("Couldn't read from the file, using defaults.")
		Conf = defaultConfig()
	} //else {} //uses the config from the file
}

func defaultConfig() Config {
	return Config{
		filePath:         "",
		Username:         "",
		Password:         "",
		RememberPassword: false,
	}
}
