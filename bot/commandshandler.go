package bot

import (
	"strings"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/database"
	db "github.com/aiuzu42/AiuzuBotDiscord/database"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	NO_AUTH        = "No tienes permiso de usar ese comando"
	GENERIC_ERROR  = "Hubo un error al procesar el comando"
	USER_NOT_FOUND = "No se encontro a ese usuario"
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
		wasCustom := false
		for _, custom := range config.Config.CustomSays {
			if custom.CommandName == args[0] {
				wasCustom = true
				sayCommand(s, m.ChannelID, custom.Channel, m.ID, st)
				break
			}
		}
		if !wasCustom {
			switch args[0] {
			case "reloadConfig":
				reloadRolesCommand(s, m.ChannelID, m.Author.ID)
			case "say":
				sayCommand(s, m.ChannelID, m.ChannelID, m.ID, st)
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
			case "ultimatum":
				ultimatumCommand(s, m, st)
			default:
				if IsMod(m.Member.Roles, m.Author.ID) {
					sendErrorResponse(s, m.ChannelID, "El comando que intentas usar no existe.")
				}
			}
		}
	}

}

func sayCommand(s *discordgo.Session, originCh string, tarCh string, id string, st string) {
	err := s.ChannelMessageDelete(originCh, id)
	if err != nil {
		log.Error("[sayCommand]Can't delete message: " + err.Error())
	}
	msg := saySplit(st)
	_, err = s.ChannelMessageSend(tarCh, msg)
	if err != nil {
		log.Error("[sayCommand]Can't send message: " + err.Error())
	}
}

func updateUserData(m *discordgo.MessageCreate) {
	if m.Author == nil {
		log.Error("[updateUserData]Message with nil author")
		return
	}
	dbErr := repo.IncreaseMessageCount(m.Author.ID)
	if dbErr != nil && dbErr.Code == db.UserNotFoundCode {
		user, errM := userAndMemberToLocalUser(m.Author, m.Member)
		if errM != nil {
			log.Error("[updateUserData]Unable to create user: " + errM.Error())
			return
		}
		dbErr := repo.AddUser(user)
		if dbErr != nil {
			log.Error("[updateUserData]Unable to store user in database: " + dbErr.Message)
		}
	} else if dbErr != nil {
		log.Error("[updateUserData]Error trying to increase message count: " + dbErr.Message)
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
	user, dbErr := repo.GetUser(args[1], args[1])
	if dbErr != nil && dbErr.Code == db.UserNotFoundCode {
		sendErrorResponse(s, m.ChannelID, "El usuario no se encontro en base de datos")
	} else if dbErr != nil {
		log.Error("[fullDetailsCommand]Database error: " + dbErr.Message)
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
	user, dbErr := repo.GetUser(args[1], args[1])
	if dbErr != nil && dbErr.Code == db.UserNotFoundCode {
		sendErrorResponse(s, m.ChannelID, "El usuario no se encontro en base de datos")
	} else if dbErr != nil {
		log.Error("[detailsCommand]Database error: " + dbErr.Message)
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
	user, dbErr := repo.GetUser(args[1], args[1])
	if dbErr != nil && dbErr.Code == db.UserNotFoundCode {
		sendErrorResponse(s, m.ChannelID, "El usuario no se encontro en base de datos")
	} else if dbErr != nil {
		log.Error("[sanctionsCommand]Database error: " + dbErr.Message)
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
		_, dbErr := repo.GetUser(member.User.ID, "")
		full := member.User.Username + "#" + member.User.Discriminator
		if dbErr != nil && dbErr.Code == db.UserNotFoundCode {
			user, err := memberToLocalUser(member)
			if err != nil {
				log.Error("[syncDatabase]Error parsing user info: " + err.Error())
				continue
			}
			log.Warn("[syncDatabase]Se agregaria a " + full + " a la db")
			dbErr := repo.AddUser(user)
			if dbErr != nil {
				log.Error("[syncDatabase]Error with new member at sync: " + dbErr.Message)
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

func ultimatumCommand(s *discordgo.Session, m *discordgo.MessageCreate, st string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[ultimatumCommand]User: " + m.Author.ID + " tried to use command sanctions without permission.")
		sendErrorResponse(s, m.ChannelID, NO_AUTH)
		return
	}
	arg, reason := argumentsHandler(st)
	if arg == "" && reason == "" {
		sendErrorResponse(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!ultimatum {userID} [razon]")
		return
	}
	dbErr := repo.SetUltimatum(arg)
	if dbErr != nil {
		log.Error("[ultimatumCommand]Error updating ultimatum in database: " + dbErr.Message)
		if dbErr.Code == database.UserNotFoundCode {
			sendErrorResponse(s, m.ChannelID, USER_NOT_FOUND)
			return
		} else if dbErr.Code == database.UserAlredyInUltimatumCode {
			sendErrorResponse(s, m.ChannelID, arg+" ya esta en Ultimatum o ha estado antes.")
		} else {
			sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
			return
		}
	}
	ult := []string{config.Config.RolUltimatum}
	err := s.GuildMemberEdit(m.GuildID, arg, ult)
	if err != nil {
		log.Error("[ultimatumCommand]Error setting roles: " + err.Error())
		sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
		return
	}
	endMsg := "Se movió a <@" + arg + "> a Ultimatum"
	log.Info(m.Author.ID + " movio a " + arg + " a Ultimatum.")
	if reason == "" {
		endMsg = endMsg + ". Recuerda poner la razon en el canal de <#" + config.Config.Channels.Ultimatum + ">"
	}
	_, err = s.ChannelMessageSend(m.ChannelID, endMsg)
	if err != nil {
		log.Error("[ultimatumCommand]Error sending success message: " + err.Error())
	}
	if reason != "" {
		_, err = s.ChannelMessageSendEmbed(config.Config.Channels.Ultimatum, createMessageEmbedUltimatum(arg, reason))
		if err != nil {
			log.Error("[ultimatumCommand]Error sending notice message: " + err.Error())
		}
	}
}

func argumentsHandler(st string) (string, string) {
	arg := ""
	msg := ""
	args := strings.Split(st, " ")
	n := len(args)
	if n == 2 {
		arg = args[1]
	} else if n > 2 {
		arg = args[1]
		msg = strings.Join(args[2:], " ")
	}
	return arg, msg
}

func saySplit(st string) string {
	msg := ""
	args := strings.Split(st, " ")
	msg = strings.Join(args[1:], " ")
	return msg
}
