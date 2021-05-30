package bot

import (
	"strconv"

	"github.com/aiuzu42/AiuzuBotDiscord/models"
	"github.com/bwmarrin/discordgo"
)

const (
	AQUA        = 1752220
	GREEN       = 3066993
	BLUE        = 3447003
	PURPLE      = 10181046
	GOLD        = 15844367
	ORANGE      = 15105570
	RED         = 15158332
	GREY        = 9807270
	DARKER_GREY = 8359053
	NAVY        = 3426654
)

func createMessageEmbedUserFull(user models.User) *discordgo.MessageEmbed {
	me := discordgo.MessageEmbed{}
	me.Title = "Información del usuario " + user.Name
	fields := []*discordgo.MessageEmbedField{}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "ID", Value: user.UserID})
	if user.Nickname != "" {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Apodo", Value: user.Nickname})
	}
	if len(user.OldNicknames) > 0 {
		chain := ""
		for _, o := range user.OldNicknames {
			chain = chain + o + " | "
		}
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Otros apodos", Value: chain})
	}
	if user.FullName != "" {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Identificador", Value: user.FullName})
	}
	if len(user.OldNames) > 0 {
		chain := ""
		for _, o := range user.OldNames {
			chain = chain + o + " | "
		}
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Otros nombres", Value: chain})
	}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Numero de sanciones", Value: strconv.Itoa(user.Sanctions.Count)})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "N. de mensajes", Value: strconv.Itoa(user.Server.MessageCount)})
	if user.Server.LastMessage != "" {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ultimo mensaje", Value: user.Server.LastMessage})
	}
	if len(user.Server.JoinDates) > 0 {
		chain := ""
		for _, d := range user.Server.JoinDates {
			chain = chain + d + ", "
		}
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Entro en las fechas", Value: chain})
	}
	if len(user.Server.LeftDates) > 0 {
		chain := ""
		for _, d := range user.Server.LeftDates {
			chain = chain + d + ", "
		}
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Se ha ido en las fechas", Value: chain})
	}
	me.Fields = fields
	return &me
}

func createMessageEmbedSancion(id string, fullName string, reason string, n int, action string) *discordgo.MessageEmbed {
	me := discordgo.MessageEmbed{}
	me.Title = "Sancion normal aplicada"
	fields := []*discordgo.MessageEmbedField{}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Sancionado:", Value: fullName})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "ID:", Value: id})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Razón", Value: reason})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Acción", Value: action})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Sanciones totales", Value: strconv.Itoa(n)})
	me.Fields = fields
	return &me
}

func createMessageGeneralHelp(commands string) *discordgo.MessageEmbed {
	me := discordgo.MessageEmbed{}
	me.Title = "Ayuda de AiuzuBot"
	fields := []*discordgo.MessageEmbedField{}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Comandos disponibles:", Value: commands})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Mas ayuda:", Value: "Para obtener mas informacion de cada comando usa:\nai!ayuda nombre_del_comando"})
	me.Fields = fields
	return &me
}

func createMessageCommandHelp(command string, desc string, synt string) *discordgo.MessageEmbed {
	me := discordgo.MessageEmbed{}
	me.Title = "Ayuda de " + command
	fields := []*discordgo.MessageEmbedField{}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Descripcion:", Value: desc})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Syntaxis:", Value: synt})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Notas:", Value: "Si algo esta entre {} significa que es obligatorio, si algo esta entre [] significa que es opcional."})
	me.Fields = fields
	return &me
}

func createMessageReport(reporterID string, reporterName string, reportedID string, reportedName string, reason string) *discordgo.MessageEmbed {
	me := discordgo.MessageEmbed{}
	me.Title = "Reporte"
	fields := []*discordgo.MessageEmbedField{}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Quien reporta:", Value: reporterID + "\n" + reporterName})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Reportado:", Value: reportedID + "\n" + reportedName})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Razón:", Value: reason})
	me.Fields = fields
	me.Color = RED
	return &me
}

func createMessageReportBasic(reporterID string, reporterName string, reason string) *discordgo.MessageEmbed {
	me := discordgo.MessageEmbed{}
	me.Title = "Reporte"
	fields := []*discordgo.MessageEmbedField{}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Quien reporta:", Value: reporterID + "\n" + reporterName})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Que reporta:", Value: reason})
	me.Fields = fields
	me.Color = RED
	return &me
}

func createMessageEmbedDMLog(user string, msg string) *discordgo.MessageEmbed {
	me := discordgo.MessageEmbed{}
	me.Title = "Aiuzu Bot DM"
	fields := []*discordgo.MessageEmbedField{}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Usuario:", Value: "<@" + user + ">"})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Msg:", Value: msg})
	me.Fields = fields
	me.Color = LIGHT_BLUE
	return &me
}
