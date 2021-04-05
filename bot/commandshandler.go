package bot

import (
	"strconv"
	"strings"
	"time"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/database"
	db "github.com/aiuzu42/AiuzuBotDiscord/database"
	"github.com/aiuzu42/AiuzuBotDiscord/version"
	"github.com/aiuzu42/AiuzuBotDiscord/youtube"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	NO_AUTH        = "No tienes permiso de usar ese comando"
	GENERIC_ERROR  = "Java lío :v tuve un error al intentar eso."
	USER_NOT_FOUND = "No se encontro a ese usuario"
	START_YT_BOT   = "Se esta intentando iniciar el bot de Youtube, favor de esperar..."
	STOP_YT_BOT    = "Se esta intentando detener el bot de Youtube, favor de esperar..."
)

// Commands handler job is to pasrse new messages to update the user data and execute bot commands if appropiate.
func CommandsHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.GuildID == "" {
		handleDM(s, m)
		return
	}
	if m.GuildID != config.Config.Server {
		gData, _ := s.Guild(m.GuildID)
		gName := ""
		if gData != nil {
			gName = gData.Name
		}
		log.Warn("[CommandsHandler]Returned because wrong server " + m.GuildID + " " + gName)
		return
	}
	updateUserData(m)

	if findIfExists(m.ChannelID, config.Config.Channels.Suggestions) {
		convertToSuggestion(s, m)
		return
	}

	if strings.HasPrefix(m.Content, prefix) == true {
		r := []rune(m.Content)
		st := string(r[pLen:])
		args := strings.Split(st, " ")
		wasCustom := false
		for _, custom := range config.Config.CustomSays {
			if custom.CommandName == args[0] {
				wasCustom = true
				customSayCommand(s, m, custom.Channel, st)
				break
			}
		}
		if !wasCustom {
			switch args[0] {
			case "reloadConfig":
				reloadRolesCommand(s, m)
			case "say":
				sayCommand(s, m.ChannelID, m.ChannelID, m.ID, st)
			case "detallesFull":
				fullDetailsCommand(s, m, args)
			case "detalles":
				detailsCommand(s, m, args)
			case "detalleSanciones":
				sanctionsCommand(s, m, args)
			case "setStatus":
				setGameStatus(s, m, args)
			case "setListenStatus":
				setListenStatus(s, m, args)
			case "setStreamStatus":
				setStreamStatus(s, m, args)
			case "syncTodos":
				syncDatabase(s, m)
			case "ultimatum":
				ultimatumCommand(s, m, st)
			case "primerAviso":
				primerAvisoCommand(s, m, st)
			case "sancionFuerte":
				sancionFuerteCommand(s, m, st)
			case "sancion":
				sancionCommand(s, m, st)
			case "ayuda", "help":
				ayudaCommand(s, m, st)
			case "actualizar":
				actualizarCommand(s, m, st)
			case "createdDate":
				createdDateCommand(s, m, args)
			case "version":
				versionCommand(s, m, args)
			case "startYt":
				go startYtCommand(s, m, args)
			case "stopYt":
				go stopYtCommand(s, m)
			case "dm", "md":
				dmCommand(s, m, st)
			default:
				if IsMod(m.Member.Roles, m.Author.ID) {
					sendMessage(s, m.ChannelID, "El comando que intentas usar no existe.", "[CommandsHandler][0]")
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

func customSayCommand(s *discordgo.Session, m *discordgo.MessageCreate, tarCh string, st string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[customSayCommand]User: " + m.Author.ID + " tried to use command customSayCommand without permission.")
		sendMessage(s, m.ChannelID, NO_AUTH, "[customSayCommand][0]")
		return
	}
	sayCommand(s, m.ChannelID, tarCh, m.ID, st)
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
		sendMessage(s, m.ChannelID, NO_AUTH, "[fullDetailsCommand][0]")
		return
	}
	if len(args) != 2 {
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!detallesFull {user}", "[fullDetailsCommand][1]")
		return
	}
	user, dbErr := repo.GetUser(args[1], args[1])
	if dbErr != nil && dbErr.Code == db.UserNotFoundCode {
		sendMessage(s, m.ChannelID, "El usuario no se encontro en base de datos", "[fullDetailsCommand][2]")
	} else if dbErr != nil {
		log.Error("[fullDetailsCommand]Database error: " + dbErr.Message)
		sendMessage(s, m.ChannelID, "Hubo un error con la base de datos", "[fullDetailsCommand][2]")
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
		sendMessage(s, m.ChannelID, NO_AUTH, "[detailsCommand][0]")
		return
	}
	if len(args) != 2 {
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!detalles {user}", "[detailsCommand][1]")
		return
	}
	user, dbErr := repo.GetUser(args[1], args[1])
	if dbErr != nil && dbErr.Code == db.UserNotFoundCode {
		sendMessage(s, m.ChannelID, "El usuario no se encontro en base de datos", "[detailsCommand][2]")
	} else if dbErr != nil {
		log.Error("[detailsCommand]Database error: " + dbErr.Message)
		sendMessage(s, m.ChannelID, "Hubo un error con la base de datos", "[detailsCommand][3]")
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
		sendMessage(s, m.ChannelID, NO_AUTH, "[sanctionsCommand][0]")
		return
	}
	if len(args) != 2 {
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!detalleSanciones {user}", "[sanctionsCommand][1]")
		return
	}
	user, dbErr := repo.GetUser(args[1], args[1])
	if dbErr != nil && dbErr.Code == db.UserNotFoundCode {
		sendMessage(s, m.ChannelID, "El usuario no se encontro en base de datos", "[sanctionsCommand][2]")
	} else if dbErr != nil {
		log.Error("[sanctionsCommand]Database error: " + dbErr.Message)
		sendMessage(s, m.ChannelID, "Hubo un error con la base de datos", "[sanctionsCommand][3]")
	} else {
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, createMessageEmbedSanctions(user))
		if err != nil {
			log.Error("[sanctionsCommand]Error sending message: " + err.Error())
		}
	}
}

func reloadRolesCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !IsOwner(m.Author.ID) {
		log.Warn("[reloadRolesCommand]User: " + m.Author.ID + " tried to use command reloadRolesCommand without permission.")
		return
	}
	if err := config.ReloadConfig(); err != nil {
		log.Error("[reloadRolesCommand]Error reloading config: " + err.Error())
		s.ChannelMessageSend(m.ChannelID, "Hubo un error, no se pudo recargar la congfiguracion")
	}
	LoadRoles()
	err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		log.Error("[reloadRolesCommand]Error marking message: " + err.Error())
	}
}

func sendMessage(s *discordgo.Session, channelID string, msg string, logMsg string) {
	_, err := s.ChannelMessageSend(channelID, msg)
	if err != nil {
		log.Error(logMsg + "Error sending message [" + msg + "]: " + err.Error())
	}
}

func setGameStatus(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsAdmin(m.Member.Roles, m.Author.ID) {
		log.Warn("[setGameStatus]User: " + m.Author.ID + " tried to use command setGameStatus without permission.")
		sendMessage(s, m.ChannelID, NO_AUTH, "[setGameStatus][0]")
		return
	}
	if len(args) < 2 {
		log.Error("[setGameStatus][1]Invalid number of arguments")
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, favor de revisar el comando", "[setGameStatus][1]")
		return
	}
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Error("[setGameStatus]Unable to delete message: " + err.Error())
	}
	msg := strings.Join(args[1:], " ")
	err = s.UpdateGameStatus(0, msg)
	if err != nil {
		log.Error("[setGameStatus]Unable to update status: " + err.Error())
	}
}

func setListenStatus(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsAdmin(m.Member.Roles, m.Author.ID) {
		log.Warn("[setListenStatus]User: " + m.Author.ID + " tried to use command setListenStatus without permission.")
		return
	}
	if len(args) < 2 {
		log.Error("[setListenStatus][1]Invalid number of arguments")
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, favor de revisar el comando", "[setListenStatus][1]")
		return
	}
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Error("[setListenStatus]Unable to delete message: " + err.Error())
	}
	msg := strings.Join(args[1:], " ")
	err = s.UpdateListeningStatus(msg)
	if err != nil {
		log.Error("[setListenStatus]Unable to update status: " + err.Error())
	}
}

func setStreamStatus(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsAdmin(m.Member.Roles, m.Author.ID) {
		log.Warn("[setStreamStatus]User: " + m.Author.ID + " tried to use command setStreamStatus without permission.")
		return
	}
	if len(args) < 3 {
		log.Error("[setStreamStatus][1]Invalid number of arguments")
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, favor de revisar el comando", "[setStreamStatus][1]")
		return
	}
	url := args[1]
	msg := strings.Join(args[2:], " ")
	if url == "" && msg == "" {
		log.Error("[setStreamStatus][2]Invalid number of arguments")
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, favor de revisar el comando", "[setStreamStatus][2]")
		return
	}
	err := s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.Error("[setStreamStatus]Unable to delete message: " + err.Error())
	}
	err = s.UpdateStreamingStatus(0, msg, url)
	if err != nil {
		log.Error("[setStreamStatus]Unable to update status: " + err.Error())
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
		sendMessage(s, m.ChannelID, "Hubo un error sincronizando", "[syncDatabase][0]")
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
		sendMessage(s, m.ChannelID, NO_AUTH, "[ultimatumCommand][0]")
		return
	}
	arg, reason := argumentsHandler(st)
	if arg == "" && reason == "" {
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!ultimatum {userID} [razon]", "[ultimatumCommand][1]")
		return
	}
	dbErr := repo.SetUltimatum(arg)
	if dbErr != nil {
		log.Error("[ultimatumCommand]Error updating ultimatum in database: " + dbErr.Message)
		if dbErr.Code == database.UserNotFoundCode {
			sendMessage(s, m.ChannelID, USER_NOT_FOUND, "[ultimatumCommand][2]")
			return
		} else if dbErr.Code == database.UserAlredyInUltimatumCode {
			sendMessage(s, m.ChannelID, arg+" ya esta en Ultimatum o ha estado antes.", "[ultimatumCommand][3]")
		} else {
			sendMessage(s, m.ChannelID, GENERIC_ERROR, "[ultimatumCommand][4]")
			return
		}
	}
	ult := []string{config.Config.Roles.Ultimatum}
	err := s.GuildMemberEdit(m.GuildID, arg, ult)
	if err != nil {
		log.Error("[ultimatumCommand]Error setting roles: " + err.Error())
		sendMessage(s, m.ChannelID, GENERIC_ERROR, "[ultimatumCommand][5]")
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

func sancionFuerteCommand(s *discordgo.Session, m *discordgo.MessageCreate, st string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[sancionFuerteCommand]User: " + m.Author.ID + " tried to use command sancionFuerteCommand without permission.")
		sendMessage(s, m.ChannelID, NO_AUTH, "[sancionFuerteCommand][0]")
		return
	}
	arg, reason := argumentsHandler(st)
	if arg == "" && reason == "" {
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!sancionFuerte {userID} [razon]", "[sancionFuerteCommand][1]")
		return
	}
	dbErr := repo.SetUltimatum(arg)
	if dbErr != nil {
		log.Error("[sancionFuerteCommand]Error updating ultimatum in database: " + dbErr.Message)
		if dbErr.Code == database.UserNotFoundCode {
			sendMessage(s, m.ChannelID, USER_NOT_FOUND, "[sancionFuerteCommand][2]")
			return
		} else if dbErr.Code != database.UserAlredyInUltimatumCode {
			sendMessage(s, m.ChannelID, GENERIC_ERROR, "[sancionFuerteCommand][3]")
			return
		}
	}
	ult := []string{config.Config.Roles.Ultimatum}
	err := s.GuildMemberEdit(m.GuildID, arg, ult)
	if err != nil {
		log.Error("[sancionFuerteCommand]Error setting roles: " + err.Error())
		sendMessage(s, m.ChannelID, GENERIC_ERROR, "[sancionFuerteCommand][4]")
		return
	}
	dbErr = repo.IncreaseSanction(arg, reason, m.Author.ID, m.Author.Username, "sancionFuerte")
	if dbErr != nil {
		log.Error("[sancionFuerteCommand]Error increasing user sanctions: " + dbErr.Message)
		sendMessage(s, m.ChannelID, GENERIC_ERROR, "[sancionFuerteCommand][5]")
		return
	}
	dbErr = repo.SetPrimerAviso(arg)
	if dbErr != nil && dbErr.Code != database.PrimerAvisoUnavailable {
		sendMessage(s, m.ChannelID, GENERIC_ERROR, "[sancionFuerteCommand][6]")
		return
	}
	endMsg := "Se aplico una sancion fuerte a <@" + arg + ">, razon: " + reason
	_, err = s.ChannelMessageSend(m.ChannelID, endMsg)
	if err != nil {
		log.Error("[sancionFuerteCommand]Error sending success message: " + err.Error())
	}
	userData, dbErr := repo.GetUser(arg, "")
	if dbErr != nil {
		log.Error("[sancionFuerteCommand]Error getting user data: " + dbErr.Message)
	}
	_, err = s.ChannelMessageSendEmbed(config.Config.Channels.Sancionados, createMessageEmbedSancionFuerte(arg, userData.FullName, reason, userData.Sanctions.Count))
	if err != nil {
		log.Error("[sancionFuerteCommand]Error sending success message 2: " + err.Error())
	}
}

func sancionCommand(s *discordgo.Session, m *discordgo.MessageCreate, st string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[sancionCommand]User: " + m.Author.ID + " tried to use command sancionCommand without permission.")
		sendMessage(s, m.ChannelID, NO_AUTH, "[sancionCommand][0]")
		return
	}
	arg, reason := argumentsHandler(st)
	if arg == "" && reason == "" {
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!sancion {userID} [razon]", "[sancionCommand][1]")
		return
	}
	memberData, err := GetMemberInfo(s, arg)
	if err != nil {
		log.Error("[sancionCommand]Error getting user data from API: " + err.Error())
		sendMessage(s, m.ChannelID, GENERIC_ERROR, "[sancionCommand][2]")
		return
	}
	userData, dbErr := repo.GetUser(arg, "")
	if dbErr != nil {
		log.Error("[sancionCommand]Error getting user data: " + dbErr.Message)
		sendMessage(s, m.ChannelID, GENERIC_ERROR, "[sancionCommand][3]")
		return
	}
	if !userData.Sanctions.Aviso {
		_, err := s.ChannelMessageSend(m.ChannelID, "<@"+arg+"> Aun tiene posibilidad de primer aviso. Considera primero darle su primer aviso con ai!primerAviso antes de sancionar")
		if err != nil {
			log.Error("[sancionCommand]Error sending ban needed message: " + err.Error())
		}
		return
	}
	action := ""
	if userData.Sanctions.Count >= 2 {
		action = "Se mando a Ultimatum por tener mas de 3 sanciones"
		err = s.GuildMemberEdit(config.Config.Server, arg, []string{config.Config.Roles.Ultimatum})
		if err != nil {
			log.Error("[sancionCommand]Error setting user roles 1: " + err.Error())
			sendMessage(s, m.ChannelID, GENERIC_ERROR, "[sancionCommand][4]")
			return
		}
		dbErr = repo.SetUltimatum(arg)
		if dbErr != nil && dbErr.Code != database.UserAlredyInUltimatumCode {
			log.Error("[sancionCommand]Error setting user to ultimatum in db: " + err.Error())
			sendMessage(s, m.ChannelID, GENERIC_ERROR, "[sancionCommand][5]")
			return
		}
	} else if Is_Ultimatum_Silenciado(memberData.Roles) {
		action = "Se aplico sanción y se debe de banear"
		_, err := s.ChannelMessageSend(m.ChannelID, "El usuario ya esta en ultimatum o silenciado!! Se debe de banear.")
		if err != nil {
			log.Error("[sancionCommand]Error sending ban needed message: " + err.Error())
		}
	} else if Is_Q_A_B(memberData.Roles) {
		action = "Se bajo a C y aplico sanción"
		nRoles := DowngradeToC(memberData.Roles)
		err = s.GuildMemberEdit(config.Config.Server, arg, nRoles)
		if err != nil {
			log.Error("[sancionCommand]Error setting user roles 2: " + err.Error())
			sendMessage(s, m.ChannelID, GENERIC_ERROR, "[sancionCommand][6]")
			return
		}
	} else {
		action = "Se aplico sanción y se silenció"
		err = s.GuildMemberEdit(config.Config.Server, arg, []string{config.Config.Roles.Silenced})
		if err != nil {
			log.Error("[sancionCommand]Error silencing user: " + err.Error())
			sendMessage(s, m.ChannelID, GENERIC_ERROR, "[sancionCommand][7]")
			return
		}
	}
	dbErr = repo.IncreaseSanction(arg, reason, m.Author.ID, m.Author.Username, "sancion")
	if dbErr != nil {
		log.Error("[sancionCommand]Error increasing user sanctions: " + dbErr.Message)
		sendMessage(s, m.ChannelID, GENERIC_ERROR, "[sancionCommand][8]")
		return
	}
	_, err = s.ChannelMessageSend(m.ChannelID, "Se aplico una sancion a: <@"+arg+">")
	if err != nil {
		log.Error("[sancionCommand]Error sending success message: " + err.Error())
	}
	_, err = s.ChannelMessageSendEmbed(config.Config.Channels.Sancionados, createMessageEmbedSancion(arg, userData.FullName, reason, userData.Sanctions.Count+1, action))
	if err != nil {
		log.Error("[sancionCommand]Error sending success message 2: " + err.Error())
	}
}

func primerAvisoCommand(s *discordgo.Session, m *discordgo.MessageCreate, st string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[primerAvisoCommand]User: " + m.Author.ID + " tried to use command primerAvisoCommand without permission.")
		sendMessage(s, m.ChannelID, NO_AUTH, "[primerAvisoCommand][0]")
		return
	}
	arg, reason := argumentsHandler(st)
	if arg == "" || reason == "" {
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!primerAviso {userID} {razon}", "[primerAvisoCommand][1]")
		return
	}
	dbErr := repo.SetPrimerAviso(arg)
	if dbErr != nil {
		if dbErr.Code == database.PrimerAvisoUnavailable {
			_, err := s.ChannelMessageSend(m.ChannelID, "A ese usuario ya se le habia dado primer aviso, toca sanción.")
			if err != nil {
				log.Error("[primerAvisoCommand]Error sending primer aviso unavailable message: " + err.Error())
			}
		} else if dbErr.Code == database.UserNotFoundCode {
			sendMessage(s, m.ChannelID, USER_NOT_FOUND, "[primerAvisoCommand][2]")
		} else {
			sendMessage(s, m.ChannelID, GENERIC_ERROR, "[primerAvisoCommand][3]")
		}
		return
	}
	_, err := s.ChannelMessageSendComplex(config.Config.Channels.Primer, createMessageComplexFirstStrike(arg, reason, config.Config.Messages.Primer, YELLOW))
	if err != nil {
		log.Error("[primerAvisoCommand]Error sending success message: " + err.Error())
	}
}

func ayudaCommand(s *discordgo.Session, m *discordgo.MessageCreate, st string) {
	var level int
	if IsOwner(m.Author.ID) {
		level = 3
	} else if IsAdmin(m.Member.Roles, m.Author.ID) {
		level = 2
	} else if IsMod(m.Member.Roles, m.Author.ID) {
		level = 1
	} else {
		level = 0
	}
	if st == "ayuda" || st == "help" {
		response := "say\nreporte"
		if level >= 1 {
			response = response + "\ndetallesFull\ndetalles\ndetalleSanciones\nultimatum\nprimerAviso\nsancion\nsancionFuerte\nactualizar\ncreatedDate\nversion\nmd\ndm"
		}
		if level >= 2 {
			response = response + "\nstartYt\nstopYt\nsetStatus\nsetListenStatus\nsetStreamStatus"
		}
		if level >= 3 {
			response = response + "\nsyncTodos\nreloadConfig"
		}
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, createMessageGeneralHelp(response))
		if err != nil {
			log.Error("[ayudaCommand]Error sending general help message: " + err.Error())
		}
		return
	}
	command := saySplit(st)
	var desc string
	var synt string
	permError := false
	switch command {
	case "reloadConfig":
		if level < 3 {
			permError = true
		}
		desc = "Carga modificaciones realizadas al archivo de configuración"
		synt = "ai!reloadConfig"
	case "say":
		desc = "El bot dice lo que le pidas y borra el mensaje original"
		synt = "ai!say {mensaje}"
	case "detallesFull":
		if level < 1 {
			permError = true
		}
		desc = "Muestra todos los detalles del usuario, excepto el desglose de las sanciones"
		synt = "ai!detallesFull {nombre_con_identificador o id}"
	case "detalles":
		if level < 1 {
			permError = true
		}
		desc = "Muestra los detalles basicos de un usuario"
		synt = "ai!detalles {nombre_con_identificador o id}"
	case "detalleSanciones":
		if level < 1 {
			permError = true
		}
		desc = "Muestra el detalle de las sanciones del usuario"
		synt = "ai!detalleSanciones {nombre_con_identificador o id}"
	case "setStatus":
		if level < 2 {
			permError = true
		}
		desc = "Actualiza el estatus de \"Jugando a\" del bot y borra el mensaje original"
		synt = "ai!setStatus {status}"
	case "setListenStatus":
		if level < 2 {
			permError = true
		}
		desc = "Actualiza el estatus de \"Escuchando a\" del bot y borra el mensaje original"
		synt = "ai!setListenStatus {status}"
	case "setStreamStatus":
		if level < 2 {
			permError = true
		}
		desc = "Actualiza el estatus de \"Streaming\" del bot y borra el mensaje original"
		synt = "ai!setStreamStatus {url} {status}"
	case "syncTodos":
		if level < 3 {
			permError = true
		}
		desc = "Revisa todos los usuarios del servidor y agrega a base de datos a los que no esten registrados, operacion pesada"
		synt = "ai!syncTodos"
	case "ultimatum":
		if level < 1 {
			permError = true
		}
		desc = "Se pasa al usuario con ese ID a ultimatum, se registra en nuestra base de datos, se le quitan todos los roles y se le asigna solo el rol de Ultimatum"
		synt = "ai!ultimatum {userID} [razon]"
	case "primerAviso":
		if level < 1 {
			permError = true
		}
		desc = "Si tiene derecho a primer aviso se aplica y notifica, si no lo tiene se notifica que se debe sancionar"
		synt = "ai!primerAviso {userID} {razon}"
	case "sancionFuerte":
		if level < 1 {
			permError = true
		}
		desc = "Se aplica una sancion fuerte y se pasa a ultimatum, AiuzuBot notifica de esto en el canal apropiado. Se registra la sanción."
		synt = "ai!sancion {userID} [razon]"
	case "sancion":
		if level < 1 {
			permError = true
		}
		desc = "Se aplica una sancion fuerte, los detalles de esto dependen del caso especifico, AiuzuBot notifica de esto en el canal apropiado. Se registra la sanción."
		synt = "ai!sancion {id} [razon]"
	case "reporte":
		desc = "Te permite reportar a algun usuario mediante su ID ai!reporte {userID} {razón} o reportar cualquier cosa que notes (fallos, situaciones, etc) ai!reporte {mensaje}"
		synt = "ai!reporte {userID} {razón} ó ai!reporte {mensaje}"
	case "actualizar":
		if level < 1 {
			permError = true
		}
		desc = "Actualiza el nombre y el apodo del usuario en caso de que se encuentre desactualizado."
		synt = "ai!actualizar {id}"
	case "ayuda", "help":
		desc = "El comando de ayuda te explica como usar los comandos de AiuzuBot y que hace cada uno."
		synt = "ai!ayuda [comando]"
	case "createdDate":
		if level < 1 {
			permError = true
		}
		desc = "Te dice la fecha de creacion de la cuenta asociada con ese ID en formato dd-MM-yyyy mm:HH"
		synt = "ai!createdDate {userID}"
	case "version":
		if level < 1 {
			permError = true
		}
		desc = "Te dice el numero de version de Aiuzu Bot"
		synt = "ai!version"
	case "stopYt":
		if level < 2 {
			permError = true
		}
		desc = "Detiene el AiuzuBot de Youtube del canal que este en las configuraciones"
		synt = "ai!stopYt"
	case "startYt":
		if level < 2 {
			permError = true
		}
		desc = "Inicia el AiuzuBot de Youtube del canal que este en las configuraciones"
		synt = "ai!startYt"
	case "dm", "md":
		if level < 1 {
			permError = true
		}
		desc = "Manda un MD al usuario del ID indicado"
		synt = "ai!dm {id} {msg} ó ai!md {id} {msg}"
	default:
		permError = true
	}
	if permError {
		_, err := s.ChannelMessageSend(m.ChannelID, "No encuentro el comando "+command+" (o no tienes permisos :v)")
		if err != nil {
			log.Error("[ayudaCommand]Error sending failed message: " + err.Error())
		}
		return
	}
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, createMessageCommandHelp(command, desc, synt))
	if err != nil {
		log.Error("[ayudaCommand]Error sending help message: " + err.Error())
	}
}

func actualizarCommand(s *discordgo.Session, m *discordgo.MessageCreate, st string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[actualizarCommand]User: " + m.Author.ID + " tried to use command actualizarCommand without permission.")
		sendMessage(s, m.ChannelID, NO_AUTH, "[actualizarCommand][0]")
		return
	}
	arg, ext := argumentsHandler(st)
	if arg == "" || ext != "" {
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!actualizar {userID}", "[actualizarCommand][1]")
		return
	}
	mem, err := GetMemberInfo(s, arg)
	if err != nil {
		log.Error("[actualizarCommand]Error getting member info: " + err.Error())
		sendMessage(s, m.ChannelID, GENERIC_ERROR, "[actualizarCommand][2]")
		return
	}
	upd, err := updateUserNames(mem.User.ID, mem.Nick, mem.User.Username, mem.User.Discriminator)
	if err != nil {
		log.Error("[actualizarCommand]Error updating info: " + err.Error())
		sendMessage(s, m.ChannelID, GENERIC_ERROR, "[actualizarCommand][3]")
		return
	}
	success := "Ya estaba actualizado."
	if upd {
		success = "Se actualizo correctamente!"
	}
	_, err = s.ChannelMessageSend(m.ChannelID, success)
	if err != nil {
		log.Error("[actualizarCommand]Error sending success message: " + err.Error())
	}
}

func convertToSuggestion(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		log.Error("[convertToSuggestion]Error adding check mark: " + err.Error())
	}
	err = s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
	if err != nil {
		log.Error("[convertToSuggestion]Error adding x: " + err.Error())
	}
}

func createdDateCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[createdDateCommand]User: " + m.Author.ID + " tried to use command createdDate without permission.")
		sendMessage(s, m.ChannelID, NO_AUTH, "[createdDateCommand][0]")
		return
	}
	if len(args) != 2 {
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!createdDate {userID}", "[createdDateCommand][1]")
		return
	}
	idInt, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		sendMessage(s, m.ChannelID, "El ID ["+args[1]+"] no es valido.", "[createdDateCommand][2]")
		return
	}
	idInt = ((idInt >> 22) + DISCORD_EPOCH) / 1000
	createdTime := time.Unix(idInt, 0)
	createdDate := createdTime.Format("02-01-2006 15:04")
	_, err = s.ChannelMessageSend(m.ChannelID, "La cuenta se creo en: "+createdDate)
	if err != nil {
		log.Error("[createdDateCommand]Error sending success message: " + err.Error())
	}
}

func versionCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[versionCommand]User: " + m.Author.ID + " tried to use command version without permission.")
		sendMessage(s, m.ChannelID, NO_AUTH, "[versionCommand][0]")
		return
	}
	if len(args) > 1 {
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!version", "[versionCommand][1]")
		return
	}
	_, err := s.ChannelMessageSend(m.ChannelID, "Aiuzu Bot version: "+version.Version)
	if err != nil {
		log.Error("[versionCommand]Error sending success message: " + err.Error())
	}
}

func startYtCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !IsAdmin(m.Member.Roles, m.Author.ID) {
		log.Warn("[startYtCommand]User: " + m.Author.ID + " tried to use command startYt without permission.")
		sendMessage(s, m.ChannelID, NO_AUTH, "[startYtCommand][0]")
		return
	}
	liveId := ""
	if len(args) > 1 {
		liveId = args[1]
	}
	nm, err := s.ChannelMessageSend(m.ChannelID, START_YT_BOT)
	if err != nil {
		log.Error("[startYtCommand]Error sending wait message: " + err.Error())
		return
	}
	var result string
	if err := youtube.StartBot(config.Config.Youtube.BotName, liveId); err == nil {
		result = " se inicio con exito el bot! ✅"
	} else {
		result = " no se pudo iniciar el bot! ❌"
	}
	_, err = s.ChannelMessageEdit(nm.ChannelID, nm.ID, START_YT_BOT+result)
	if err != nil {
		log.Error("[startYtCommand]Error sending result message: " + err.Error())
	}
}

func stopYtCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !IsAdmin(m.Member.Roles, m.Author.ID) {
		log.Warn("[stopYtCommand]User: " + m.Author.ID + " tried to use command stopYt without permission.")
		sendMessage(s, m.ChannelID, NO_AUTH, "[stopYtCommand][0]")
		return
	}
	nm, err := s.ChannelMessageSend(m.ChannelID, STOP_YT_BOT)
	if err != nil {
		log.Error("[stopYtCommand]Error sending wait message: " + err.Error())
		return
	}
	var result string
	if err := youtube.StopBot(config.Config.Youtube.BotName); err == nil {
		result = " se detuvo con exito el bot! ✅"
	} else {
		result = " no se pudo detener el bot! ❌"
	}
	_, err = s.ChannelMessageEdit(nm.ChannelID, nm.ID, STOP_YT_BOT+result)
	if err != nil {
		log.Error("[stopYtCommand]Error sending result message: " + err.Error())
	}
}

func dmCommand(s *discordgo.Session, m *discordgo.MessageCreate, st string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[dmCommand]User: " + m.Author.ID + " tried to use command dmCommand without permission.")
		return
	}
	arg, msg := argumentsHandler(st)
	if arg == "" || msg == "" {
		sendMessage(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!md {id} {msg}", "[dmCommand][1]")
		return
	}
	dmc, err := s.UserChannelCreate(arg)
	if err != nil {
		log.Error("[dmCommand]Error creating DM channel: " + err.Error())
		sendMessage(s, m.ChannelID, GENERIC_ERROR, "[sendMDCommand][2]")
		return
	}
	_, err = s.ChannelMessageSend(dmc.ID, msg)
	if err != nil {
		log.Error("[dmCommand]Error sending DM: " + err.Error())
		sendMessage(s, m.ChannelID, GENERIC_ERROR, "[sendMDCommand][3]")
	}
}
