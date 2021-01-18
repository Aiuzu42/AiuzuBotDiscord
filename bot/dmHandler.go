package bot

import (
	"strings"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/database"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func handleDM(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Info("DM from: " + m.Author.ID + ": " + m.Content)
	if strings.HasPrefix(m.Content, prefix) == true {
		r := []rune(m.Content)
		st := string(r[pLen:])
		args := strings.Split(st, " ")
		switch args[0] {
		case "reporte":
			dmReportCommand(s, m, st)
		}
	}
}

func dmReportCommand(s *discordgo.Session, m *discordgo.MessageCreate, st string) {
	arg, reason := argumentsHandler(st)
	if arg == "" && reason == "" {
		sendErrorResponse(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!reporte {mensaje} ó ai!reporte {userID} {razon}")
		return
	}
	userData, dbErr := repo.GetUser(arg, "")
	if dbErr != nil {
		if dbErr.Code != database.UserNotFoundCode {
			log.Error("[dmReportCommand]Error finding user: " + dbErr.Message)
		}
		_, err := s.ChannelMessageSendEmbed(config.Config.Channels.Reports, createMessageReportBasic(m.Author.ID, m.Author.Username+"#"+m.Author.Discriminator, arg+" "+reason))
		if err != nil {
			log.Error("[dmReportCommand]Error sending message basic: " + err.Error())
		}
	} else {
		if reason == "" {
			reason = "NA"
		}
		_, err := s.ChannelMessageSendEmbed(config.Config.Channels.Reports, createMessageReport(m.Author.ID, m.Author.Username+"#"+m.Author.Discriminator, arg, userData.FullName, reason))
		if err != nil {
			log.Error("[dmReportCommand]Error sending message: " + err.Error())
		}
	}
}
