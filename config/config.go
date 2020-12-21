package config

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type configuration struct {
	Server     string       `json:"server"`
	Token      string       `json:"token"`
	LogLevel   string       `json:"logLevel"`
	DBConn     DBConnection `json:"dbConnection"`
	Owners     []string     `json:"owners"`
	Mods       []string     `json:"mods"`
	Admins     []string     `json:"admins"`
	Channels   ChannelsInfo `json:"channels"`
	Roles      RolesInfo    `json:"roles"`
	CustomSays []CustomSay  `json:"customSays"`
}

// DBConnections contains all the needed information to connect to a database.
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

type ChannelsInfo struct {
	F           string `json:"f"`
	Ultimatum   string `json:"ultimatum"`
	Sancionados string `json:"sancionados"`
}

type CustomSay struct {
	CommandName string `json:"commandName"`
	Channel     string `json:"channel"`
}

type RolesInfo struct {
	Ultimatum string `json:"ultimatum"`
	Q         string `json:"q"`
	A         string `json:"a"`
	B         string `json:"b"`
	C         string `json:"c"`
	Silenced  string `json:"silenced"`
}

const (
	filePath = "config.json"
)

// Config contains the bot configuration.
var Config configuration
var Loc *time.Location

//InitConfig should be only used to load config at the start of the program, it panics if the config cannot be loaded for any reason.
func InitConfig() error {
	var err error
	Config, err = loadConfig()
	if err != nil {
		return err
	}
	switch Config.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	Loc = time.FixedZone("UTC-6", -6*60*60)
	return nil
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
