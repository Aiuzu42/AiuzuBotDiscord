package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aiuzu42/AiuzuBotDiscord/bot"
	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/version"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	AppSession App
)

// StartApp Initialize the bot.
// It load roles from config, start discordgo session, selects repository type, adds handlers, prints bot version and handles stop condition.
func StartApp() {
	AppSession = App{}
	var err error
	bot.LoadRoles()
	AppSession.b, err = discordgo.New("Bot " + config.Config.Token)
	if err != nil {
		log.Fatal("[StartApp]Error starting up discordgo: " + err.Error())
	}
	err = bot.SelectRepository(config.Config.DBConn.DBType)
	if err != nil {
		log.Fatal("[StartApp]Error selecting repository: " + err.Error())
	}
	AppSession.b.Identify.Intents = discordgo.IntentsAll
	AppSession.b.AddHandler(bot.CommandsHandler)
	AppSession.b.AddHandler(bot.NewMemberHandler)
	AppSession.b.AddHandler(bot.MemberLeaveHandler)
	AppSession.b.AddHandler(bot.MemberUpdateHandler)
	err = AppSession.b.Open()
	if err != nil {
		log.Fatal("[StartApp]Error opening Discord websocket connection: " + err.Error())
	}
	http.HandleFunc("/", indexController)
	http.HandleFunc("/msg", msgController)
	http.HandleFunc("/edit", editController)
	http.HandleFunc("/editmsg", editMsgController)
	http.HandleFunc("/msgembed", msgEmbedController)
	http.HandleFunc("/editmsgembed", editMsgEmbedController)
	http.HandleFunc("/sendMsgDirect", sendMsgDirectController)
	go func() {
		err := http.ListenAndServe(":9090", nil)
		if err != nil {
			log.Fatal("[StartApp]Error initiating web server: " + err.Error())
		}
	}()
	fmt.Println("AiuzuBot Discord is now running. Version: " + version.Version)
	log.Info("AiuzuBot Discord is now running. Version: " + version.Version)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	AppSession.b.Close()
}
