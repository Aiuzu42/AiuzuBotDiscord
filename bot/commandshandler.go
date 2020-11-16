package bot

import (
	"strings"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/models"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	NO_AUTH = "No tienes permiso de usar ese comando"
)

// Commands handler job is to pasrse new messages to update the user data and execute bot commands if appropiate.
func CommandsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.GuildID != config.Config.Server {
		gName, _ := s.Guild(m.GuildID)
		log.Warn("[CommandsHandler]Returned because wrong server " + m.GuildID + " " + gName.Name)
		return
	}

	updateUserData(m)

	if strings.HasPrefix(m.Content, prefix) == true {
		r := []rune(m.Content)
		st := string(r[pLen:])
		args := strings.Split(st, " ")
		switch args[0] {
		case "reloadConfig":
			reloadRolesCommand(s, m.ChannelID, m.Author.ID)
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
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Error("[sayCommand]Can't delete message: " + err.Error())
	}
	msg := string(r[pLen+4:])
	_, err = s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		log.Error("[sayCommand]Can't send message: " + err.Error())
	}
}

func updateUserData(m *discordgo.MessageCreate) {
	if m.Author == nil {
		log.Error("[updateUserData]Message with nil author")
		return
	}
	err := repo.IncreaseMessageCount(m.Author.ID)
	if err != nil && err.Code == models.UserNotFoundCode {
		user, errM := userAndMemberToLocalUser(m.Author, m.Member)
		if errM != nil {
			log.Error("[updateUserData]Unable to create user: " + errM.Error())
			return
		}
		appErr := repo.AddUser(user)
		if appErr != nil {
			log.Error("[updateUserData]Unable to store user in database: " + appErr.Message)
		}
	} else if err != nil {
		log.Error("[updateUserData]Error trying to increase message count: " + err.Message)
	}
}

func fullDetailsCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[fullDetailsCommand]User: " + m.Author.ID + " tried to use command fullDetails without permission.")
		sendErrorResponse(s, m.ChannelID, NO_AUTH)
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
		log.Error("[fullDetailsCommand]Database error: " + appErr.Message)
		sendErrorResponse(s, m.ChannelID, "Hubo un error con la base de datos")
	} else {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, createMessageEmbedUserFull(user))
		if err != nil {
			log.Error("[fullDetailsCommand]Error sending message: " + err.Error())
		}
	}
}

func detailsCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[detailsCommand]User: " + m.Author.ID + " tried to use command details without permission.")
		sendErrorResponse(s, m.ChannelID, NO_AUTH)
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
		log.Error("[detailsCommand]Database error: " + appErr.Message)
		sendErrorResponse(s, m.ChannelID, "Hubo un error con la base de datos")
	} else {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, createMessageEmbedUser(user))
		if err != nil {
			log.Error("[detailsCommand]Error sending message: " + err.Error())
		}
	}
}

func sanctionsCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[sanctionsCommand]User: " + m.Author.ID + " tried to use command sanctions without permission.")
		sendErrorResponse(s, m.ChannelID, NO_AUTH)
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
		log.Error("[sanctionsCommand]Database error: " + appErr.Message)
		sendErrorResponse(s, m.ChannelID, "Hubo un error con la base de datos")
	} else {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, createMessageEmbedSanctions(user))
		if err != nil {
			log.Error("[sanctionsCommand]Error sending message: " + err.Error())
		}
	}
}

func reloadRolesCommand(s *discordgo.Session, channelID string, id string) {
	if !IsOwner(id) {
		log.Warn("[reloadRolesCommand]User: " + id + " tried to use command reloadRolesCommand without permission.")
		return
	}
	if err := config.ReloadConfig(); err != nil {
		log.Error("[reloadRolesCommand]Error reloading config: " + err.Error())
		s.ChannelMessageSend(channelID, "Hubo un error, no se pudo recargar la congfiguracion")
	}
	LoadRoles()
}

func sendErrorResponse(s *discordgo.Session, channelID string, msg string) {
	_, err := s.ChannelMessageSend(channelID, msg)
	if err != nil {
		log.Error("[sendErrorResponse]Error sending message [" + msg + "]: " + err.Error())
	}
}

func setStatus(s *discordgo.Session, m *discordgo.MessageCreate, r []rune) {
	if !IsOwner(m.Author.ID) {
		log.Warn("[setStatus]User: " + m.Author.ID + " tried to use command setStatus without permission.")
		sendErrorResponse(s, m.ChannelID, NO_AUTH)
		return
	}
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Error("[setStatus]Unable to delete message: " + err.Error())
	}
	msg := string(r[pLen+10:])
	err = s.UpdateStatus(0, msg)
	if err != nil {
		log.Error("[setStatus]Unable to update status: " + err.Error())
	}
}

func syncDatabase(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !IsOwner(m.Author.ID) {
		log.Warn("[syncDatabase]User: " + m.Author.ID + " tried to use command syncDatabase without permission.")
		return
	}
	members, err := s.GuildMembers(m.GuildID, "", 500)
	if err != nil {
		log.Error("[syncDatabase]Unable to obtain guild members: " + err.Error())
		sendErrorResponse(s, m.ChannelID, "Hubo un error sincronizando")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Iniciando sincronizacion...")
	log.Warnf("[syncDatabase]Len: %v", len(members))
	for _, member := range members {
		if member.User == nil {
			log.Warn("[syncDatabase]Member without user " + m.ID)
			continue
		}
		_, appErr := repo.GetUser(member.User.ID, "")
		full := member.User.Username + "#" + member.User.Discriminator
		if appErr != nil && appErr.Code == models.UserNotFoundCode {
			user, err := memberToLocalUser(member)
			if err != nil {
				log.Error("[syncDatabase]Error parsing user info: " + err.Error())
				continue
			}
			log.Warn("[syncDatabase]Se agregaria a " + full + " a la db")
			appErr := repo.AddUser(user)
			if appErr != nil {
				log.Error("[syncDatabase]Error with new member at sync: " + appErr.Message)
			}
		} else {
			log.Warn("[syncDatabase]Ya estaba en db: " + full)
		}
	}
	_, err = s.ChannelMessageSend(m.ChannelID, "Sincronizacion terminada")
	if err != nil {
		log.Error("[syncDatabase]Error sending end message: " + err.Error())
	}
}
