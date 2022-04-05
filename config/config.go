package config

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/aiuzu42/AiuzuBotDiscord/youtube"
	log "github.com/sirupsen/logrus"
)

type configuration struct {
	Server            string                `json:"server"`
	Token             string                `json:"token"`
	LogLevel          string                `json:"logLevel"`
	DBConn            DBConnection          `json:"dbConnection"`
	Owners            []string              `json:"owners"`
	Mods              []string              `json:"mods"`
	Admins            []string              `json:"admins"`
	Channels          ChannelsInfo          `json:"channels"`
	Roles             RolesInfo             `json:"roles"`
	CustomSays        []CustomSay           `json:"customSays"`
	Youtube           YoutubeData           `json:"youtube"`
	Messages          BotMessages           `json:"messages"`
	LeaveNotification LeaveNotificationData `json:"leaveNotification"`
	Vxp               VxpConfig             `json:"vxp"`
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
	Sancionados string   `json:"sancionados"`
	Suggestions []string `json:"suggestions"`
	Reports     string   `json:"reports"`
	BotDM       string   `json:"botDM"`
	Youtube     string   `json:"youtube"`
	Upgrades    string   `json:"upgrades"`
	Logs        string   `json:"logs"`
}

type CustomSay struct {
	CommandName string `json:"commandName"`
	Channel     string `json:"channel"`
}

type RolesInfo struct {
	Q string `json:"q"`
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"c"`
}

type YoutubeData struct {
	Url         string `json:"url"`
	BotName     string `json:"botName"`
	SendMessage bool   `json:"sendMessage"`
	SetStatus   bool   `json:"setStatus"`
	Message     string `json:"message"`
	StatusMsg   string `json:"statusMsg"`
}

type BotMessages struct {
	Primer string `json:"primerAviso"`
}

type LeaveNotificationData struct {
	Active  bool   `json:"active"`
	Channel string `json:"channel"`
}

type VxpConfig struct {
	VxpMultipliers  []VxpMultiplier `json:"vxpMultipliers"`
	Active          bool            `json:"active"`
	IgnoredChannels []string        `json:"ignoredChannels"`
	RolUpgrades     []RolUpgrade    `json:"rolUpgrades"`
	MaxPerDay       int             `json:"maxPerDay"`
}

type VxpMultiplier struct {
	Rol  string `json:"rol"`
	Mult int    `json:"mult"`
}

type RolUpgrade struct {
	Rol      string   `json:"rol"`
	Value    int      `json:"value"`
	ToDelete []string `json:"toDelete"`
}

const (
	filePath = "config.json"
	apiUsrKey = "API_USR"
	apiPassKey = "API_PASS"
)

// Config contains the bot configuration.
var (
	Config configuration
	Loc *time.Location
	ApiUsr string
	ApiPass string
)

//InitConfig should be only used to load config at the start of the program, it panics if the config cannot be loaded for any reason.
func InitConfig() error {
	var err error
	Config, err = loadConfig()
	if err != nil {
		return err
	}
	err = loadEnvVariables()
	if err != nil {
		return err
	}
	switch Config.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	setAppVariables()
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
	err = loadEnvVariables()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	Config = localConfig
	setAppVariables()
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

func setAppVariables() {
	youtube.ApiUrl = Config.Youtube.Url
}

func loadEnvVariables() error {
	ApiPass = os.Getenv(apiPassKey)
	ApiUsr = os.Getenv(apiUsrKey)
	if ApiUsr == "" || ApiPass == "" {
		return errors.New("env variables not set")
	}
	return nil
}
