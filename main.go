package main

import (
	"os"

	"github.com/aiuzu42/AiuzuBotDiscord/app"
	"github.com/aiuzu42/AiuzuBotDiscord/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	setupLog()
	err := config.InitConfig()
	if err != nil {
		log.Fatal("[main]Error loading configuration: " + err.Error())
	}
	app.StartApp()
}

func setupLog() {
	file, err := os.OpenFile("discordLog.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
}
