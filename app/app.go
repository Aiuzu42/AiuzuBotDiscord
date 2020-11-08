package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aiuzu42/AiuzuBotDiscord/bot"
	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/version"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func StartApp() {
	bot.LoadRoles()
	b, err := discordgo.New("Bot " + config.Config.Token)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = bot.SelectRepository(config.Config.DBConn.DBType)
	if err != nil {
		os.Exit(1)
	}
	b.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	b.AddHandler(bot.CommandsHandler)
	b.AddHandler(bot.NewMemberHandler)
	b.AddHandler(bot.MemberLeaveHandler)
	err = b.Open()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("AiuzuBot Discord is now running. Version: " + version.Version)
	log.Info("AiuzuBot Discord is now running. Version: " + version.Version)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	b.Close()
}
