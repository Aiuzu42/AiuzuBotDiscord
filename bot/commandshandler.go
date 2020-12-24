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
				customSayCommand(s, m, custom.Channel, st)
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
			case "primerAviso":
				primerAvisoCommand(s, m, st)
			case "sancionFuerte":
				sancionFuerteCommand(s, m, st)
			case "sancion":
				sancionCommand(s, m, st)
			case "ayuda":
				ayudaCommand(s, m, st)
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

func customSayCommand(s *discordgo.Session, m *discordgo.MessageCreate, tarCh string, st string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[customSayCommand]User: " + m.Author.ID + " tried to use command customSayCommand without permission.")
		sendErrorResponse(s, m.ChannelID, NO_AUTH)
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
	ult := []string{config.Config.Roles.Ultimatum}
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

func sancionFuerteCommand(s *discordgo.Session, m *discordgo.MessageCreate, st string) {
	if !IsMod(m.Member.Roles, m.Author.ID) {
		log.Warn("[sancionFuerteCommand]User: " + m.Author.ID + " tried to use command sancionFuerteCommand without permission.")
		sendErrorResponse(s, m.ChannelID, NO_AUTH)
		return
	}
	arg, reason := argumentsHandler(st)
	if arg == "" && reason == "" {
		sendErrorResponse(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!sancionFuerte {userID} [razon]")
		return
	}
	dbErr := repo.SetUltimatum(arg)
	if dbErr != nil {
		log.Error("[sancionFuerteCommand]Error updating ultimatum in database: " + dbErr.Message)
		if dbErr.Code == database.UserNotFoundCode {
			sendErrorResponse(s, m.ChannelID, USER_NOT_FOUND)
			return
		} else if dbErr.Code != database.UserAlredyInUltimatumCode {
			sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
			return
		}
	}
	ult := []string{config.Config.Roles.Ultimatum}
	err := s.GuildMemberEdit(m.GuildID, arg, ult)
	if err != nil {
		log.Error("[sancionFuerteCommand]Error setting roles: " + err.Error())
		sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
		return
	}
	dbErr = repo.IncreaseSanction(arg, reason, m.Author.ID, m.Author.Username, "sancionFuerte")
	if dbErr != nil {
		log.Error("[sancionFuerteCommand]Error increasing user sanctions: " + dbErr.Message)
		sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
		return
	}
	dbErr = repo.SetPrimerAviso(arg)
	if dbErr != nil && dbErr.Code != database.PrimerAvisoUnavailable {
		sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
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
		sendErrorResponse(s, m.ChannelID, NO_AUTH)
		return
	}
	arg, reason := argumentsHandler(st)
	if arg == "" && reason == "" {
		sendErrorResponse(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!sancion {userID} [razon]")
		return
	}
	memberData, err := GetMemberInfo(s, arg)
	if err != nil {
		log.Error("[sancionCommand]Error getting user data from API: " + err.Error())
		sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
		return
	}
	userData, dbErr := repo.GetUser(arg, "")
	if dbErr != nil {
		log.Error("[sancionCommand]Error getting user data: " + dbErr.Message)
		sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
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
			sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
			return
		}
		dbErr = repo.SetUltimatum(arg)
		if dbErr != nil && dbErr.Code != database.UserAlredyInUltimatumCode {
			log.Error("[sancionCommand]Error setting user to ultimatum in db: " + err.Error())
			sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
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
			sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
			return
		}
	} else {
		action = "Se aplico sanción y se silenció"
		err = s.GuildMemberEdit(config.Config.Server, arg, []string{config.Config.Roles.Silenced})
		if err != nil {
			log.Error("[sancionCommand]Error silencing user: " + err.Error())
			sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
			return
		}
	}
	dbErr = repo.IncreaseSanction(arg, reason, m.Author.ID, m.Author.Username, "sancion")
	if dbErr != nil {
		log.Error("[sancionCommand]Error increasing user sanctions: " + dbErr.Message)
		sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
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
		sendErrorResponse(s, m.ChannelID, NO_AUTH)
		return
	}
	arg, reason := argumentsHandler(st)
	if arg == "" && reason == "" {
		sendErrorResponse(s, m.ChannelID, "Numero de argumentos incorrecto, el comando es: ai!primerAviso {userID} [razon]")
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
			sendErrorResponse(s, m.ChannelID, USER_NOT_FOUND)
		} else {
			sendErrorResponse(s, m.ChannelID, GENERIC_ERROR)
		}
		return
	}
	_, err := s.ChannelMessageSend(m.ChannelID, "Se le aplico correctamente primer aviso a <@"+arg+">")
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
	if st == "ayuda" {
		response := "say"
		if level >= 1 {
			response = response + "\ndetallesFull\ndetalles\ndetalleSanciones\nultimatum\nprimerAviso\nsancion\nsancionFuerte"
		}
		if level >= 3 {
			response = response + "\nsyncTodos\nsetStatus\nreloadConfig"
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
		if level < 3 {
			permError = true
		}
		desc = "Actualiza el estatus del bot y borra el mensaje original"
		synt = "ai!setStatus {status}"
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
		synt = "ai!primerAviso {userID} [razon]"
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
	case "ayuda":
		desc = "El comando de ayuda te explica como usar los comandos de AiuzuBot y que hace cada uno."
		synt = "ai!ayuda [comando]"
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
