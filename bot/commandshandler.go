package bot

import (
	"strings"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/models"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func CommandsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.GuildID != config.Config.Server {
		gName, _ := s.Guild(m.GuildID)
		log.Warn("Returned because wrong server " + m.GuildID + " " + gName.Name)
		return
	}

	updateUserData(m)

	if strings.HasPrefix(m.Content, prefix) == true {
		r := []rune(m.Content)
		st := string(r[pLen:])
		args := strings.Split(st, " ")
		switch args[0] {
		case "reloadConfig":
			reloadRolesCommand(s, m.ChannelID)
		case "say":
			sayCommand(s, m, r)
		case "detallesFull":
			fullDetailsCommand(s, m, args)
		case "detalles":
			detailsCommand(s, m, args)
		case "detalleSanciones":
			sanctionsCommand(s, m, args)
		case "setStatus":
			setStatus(s, m, r)
		case "syncTodos":
			syncDatabase(s, m)
		default:
			if IsMod(m.Member.Roles, m.Author.ID) {
				sendErrorResponse(s, m.ChannelID, "El comando que intentas usar no existe.")
			}
		}
	}

}

func sayCommand(s *discordgo.Session, m *discordgo.MessageCreate, r []rune) {
	if !IsOwner(m.Author.ID) {
		log.Warn("User: " + m.Author.ID + " tried to use command details without permission.")
		return
	}
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Error("say command cant delete message" + err.Error())
	}
	msg := string(r[pLen+4:])
	_, err = s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		log.Error("say command cant send message" + err.Error())
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

func fullDetailsCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("User: " + m.Author.ID + " tried to use command fullDetails without permission.")
		sendErrorResponse(s, m.ChannelID, "No tienes permiso de usar ese comando")
		return
	}
	if len(args) != 2 {
		sendErrorResponse(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!details {user}")
		return
	}
	user, appErr := repo.GetUser(args[1], args[1])
	if appErr != nil && appErr.Code == models.UserNotFoundCode {
		sendErrorResponse(s, m.ChannelID, "El usuario no se encontro en base de datos")
	} else if appErr != nil {
		sendErrorResponse(s, m.ChannelID, "Hubo un error con la base de datos")
	} else {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, createMessageEmbedUserFull(user))
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func detailsCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("User: " + m.Author.ID + " tried to use command details without permission.")
		sendErrorResponse(s, m.ChannelID, "No tienes permiso de usar ese comando")
		return
	}
	if len(args) != 2 {
		sendErrorResponse(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!details {user}")
		return
	}
	user, appErr := repo.GetUser(args[1], args[1])
	if appErr != nil && appErr.Code == models.UserNotFoundCode {
		sendErrorResponse(s, m.ChannelID, "El usuario no se encontro en base de datos")
	} else if appErr != nil {
		sendErrorResponse(s, m.ChannelID, "Hubo un error con la base de datos")
	} else {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, createMessageEmbedUser(user))
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func sanctionsCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("User: " + m.Author.ID + " tried to use command sanctions without permission.")
		sendErrorResponse(s, m.ChannelID, "No tienes permiso de usar ese comando")
		return
	}
	if len(args) != 2 {
		sendErrorResponse(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!detalleSanciones {user}")
		return
	}
	user, appErr := repo.GetUser(args[1], args[1])
	if appErr != nil && appErr.Code == models.UserNotFoundCode {
		sendErrorResponse(s, m.ChannelID, "El usuario no se encontro en base de datos")
	} else if appErr != nil {
		sendErrorResponse(s, m.ChannelID, "Hubo un error con la base de datos")
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

func setStatus(s *discordgo.Session, m *discordgo.MessageCreate, r []rune) {
	if !IsOwner(m.Author.ID) {
		log.Warn("User: " + m.Author.ID + " tried to use command setStatus without permission.")
		return
	}
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Error("setStatus command cant delete message" + err.Error())
	}
	msg := string(r[pLen+10:])
	err = s.UpdateStatus(0, msg)
	if err != nil {
		log.Error("setStatus command cant update status" + err.Error())
	}
}

func syncDatabase(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !IsOwner(m.Author.ID) {
		log.Warn("User: " + m.Author.ID + " tried to use command syncDatabase without permission.")
		return
	}
	members, err := s.GuildMembers(m.GuildID, "", 500)
	if err != nil {
		log.Error("SyncDatabase 1: " + err.Error())
		sendErrorResponse(s, m.ChannelID, "Hubo un error sincronizando")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Iniciando sincronizacion...")
	log.Warn(len(members))
	for _, member := range members {
		if member.User == nil {
			log.Warn("Member without user " + m.ID)
			continue
		}
		_, appErr := repo.GetUser(member.User.ID, "")
		full := member.User.Username + "#" + member.User.Discriminator
		if appErr != nil && appErr.Code == models.UserNotFoundCode {
			user := models.User{UserID: member.User.ID, Name: member.User.Username, Nickname: member.Nick, FullName: full}
			user.Server.JoinDates = append(user.Server.JoinDates, string(member.JoinedAt))
			user.Server.MessageCount = 0
			user.Server.WasModerator = false
			if findIfExists(config.Config.RolUltimatum, member.Roles) {
				user.Server.Ultimatum = true
			} else {
				user.Server.Ultimatum = false
			}
			log.Warn("Se agregaria a " + full + " a la db")
			appErr := repo.AddUser(user)
			if appErr != nil {
				log.Error("error with new member at sync " + appErr.Message)
			}
		} else {
			log.Warn("Ya estaba en db: " + full)
		}
	}
	_, err = s.ChannelMessageSend(m.ChannelID, "Sincronizacion terminada")
	if err != nil {
		log.Error(err.Error())
	}
}
