package bot

import (
	"strings"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/models"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func CommandsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID != config.Config.Server || m.Author.ID == s.State.User.ID {
		return
	}

	updateUserData(m)

	if strings.HasPrefix(m.Content, prefix) == true {
		r := []rune(m.Content)
		st := string(r[pLen:])
		words := strings.Split(st, " ")
		switch words[0] {
		case "reloadRoles":
			reloadRolesCommand(s, m.ChannelID)
		default:
			if IsMod(m.Member.Roles, m.Author.ID) {
				sendErrorResponse(s, m.ChannelID, "El comando que intentas usar no existe.")
			}
		}
	}

}

func updateUserData(m *discordgo.MessageCreate) {
	if err := repo.IncreaseMessageCount(m.Author.ID); err != nil && err.Code == models.UserNotFoundCode {
		full := m.Author.Username + "#" + m.Author.Discriminator
		user := models.User{UserID: m.Author.ID, Name: m.Author.Username, FullName: full}
		user.Server.MessageCount = user.Server.MessageCount + 1
		user.Server.WasModerator = false
		user.Server.Ultimatum = false
		if m.Member != nil {
			user.Nickname = m.Member.Nick
			user.Server.JoinDates = append(user.Server.JoinDates, string(m.Member.JoinedAt))
		}
		repo.AddUser(user)
	}
}

func detailsCommand(s *discordgo.Session, args []string, m *discordgo.MessageCreate) {
	if !IsAdmin(m.Member.Roles, m.Author.ID) {
		log.Warn("User: " + m.Author.ID + " tried to use command details without permission.")
		return
	}
	if len(args) != 2 {
		sendErrorResponse(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!details {user}")
		return
	}
	user, appErr := repo.GetUser("", args[1])
	if appErr != nil && appErr.Code == models.UserNotFoundCode {
		sendErrorResponse(s, m.ChannelID, "El usuario no existe en base de datos")
	} else if appErr != nil {
		sendErrorResponse(s, m.ChannelID, "Hubo un error con la base de datos")
	} else {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, createMessageEmbedUser(user))
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func sancionesCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	if len(args) != 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!sanciones {user}")
		if err != nil {
			log.Error(err.Error())
		}
		return
	}
	user, err := repo.GetUser("", args[1])
	if err != nil {
		log.Error(err.String())
	} else {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, createMessageEmbedSanctions(user))
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func reloadRolesCommand(s *discordgo.Session, channelID string) {
	if err := config.ReloadConfig(); err != nil {
		s.ChannelMessageSend(channelID, "Hubo un error, no se pudo cargar la congfiguracion")
	}
	LoadRoles()
}

func sendErrorResponse(s *discordgo.Session, channelID string, msg string) {
	_, err := s.ChannelMessageSend(channelID, msg)
	if err != nil {
		log.Error("Error sending message: " + msg + " because " + err.Error())
	}
}
