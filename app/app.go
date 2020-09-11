package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aiuzu42/AiuzuBotDiscord/bot"
	"github.com/bwmarrin/discordgo"
)

func StartApp(token string) {
	b, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err.Error())
	}
	b.AddHandler(bot.CommandsHandler)
	err = b.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Discord bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	b.Close()
}
