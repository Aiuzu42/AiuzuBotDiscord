package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type configuration struct {
	Server   string       `json:"server"`
	Token    string       `json:"token"`
	DBConn   DBConnection `json:"dbConnection"`
	Owners   []string     `json:"owners"`
	Mods     []string     `json:"mods"`
	Admins   []string     `json:"admins"`
	FChannel string       `json:"fChannel"`
}

type DBConnection struct {
	DBType     string `json:"type"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	DB         string `json:"db"`
	Collection string `json:"collection"`
	User       string `json:"user"`
	Pass       string `json:"pass"`
	Options    string `json:"options"`
}

const (
	filePath = "config.json"
)

var Config configuration

//InitConfig should be only used to load config at the start of the program, it panics if the config cannot be loaded for any reason.
func InitConfig() {
	var err error
	Config, err = loadConfig()
	if err != nil {
		fmt.Println(err.Error())
		log.Panic(err.Error())
	}
}

//ReloadConfig can be used to reload config at any point, if it fails to reload it keeps the old config and returns an error.
func ReloadConfig() error {
	localConfig, err := loadConfig()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	Config = localConfig
	return nil
}

func loadConfig() (configuration, error) {
	var localConfig configuration
	file, err := os.Open(filePath)
	if err != nil {
		return configuration{}, errors.New("Unable to load configuration file [" + err.Error() + "]")
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&localConfig)
	if err != nil {
		return configuration{}, errors.New("Error parsing configuration file [" + err.Error() + "]")
	}
	return localConfig, nil
}
