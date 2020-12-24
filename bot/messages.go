package bot

import (
	"strconv"

	"github.com/aiuzu42/AiuzuBotDiscord/models"
	"github.com/bwmarrin/discordgo"
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
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Numero de sanciones", Value: strconv.Itoa(user.Sanctions.Count)})
	if user.Sanctions.Aviso {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Primer aviso disponible?", Value: "NO"})
	} else {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Primer aviso disponible?", Value: "SI"})
	}
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
	if user.Server.Ultimatum {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ha estado en ultimatum?", Value: "SI"})
	} else {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ha estado en ultimatum?", Value: "NO"})
	}
	me.Fields = fields
	return &me
}

func createMessageEmbedUser(user models.User) *discordgo.MessageEmbed {
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
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Numero de sanciones", Value: strconv.Itoa(user.Sanctions.Count)})
	if user.Sanctions.Aviso {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Primer aviso disponible?", Value: "NO"})
	} else {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Primer aviso disponible?", Value: "SI"})
	}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "N. de mensajes", Value: strconv.Itoa(user.Server.MessageCount)})
	if user.Server.LastMessage != "" {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ultimo mensaje", Value: user.Server.LastMessage})
	}
	if user.Server.Ultimatum {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ha estado en ultimatum?", Value: "SI"})
	} else {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ha estado en ultimatum?", Value: "NO"})
	}
	me.Fields = fields
	return &me
}

func createMessageEmbedSanctions(user models.User) *discordgo.MessageEmbed {
	me := discordgo.MessageEmbed{}
	me.Title = "Sanciones del usuario " + user.Name
	fields := []*discordgo.MessageEmbedField{}
	if user.Sanctions.Aviso {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ya tiene un primer aviso?", Value: "SI"})
	} else {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ya tiene un primer aviso?", Value: "NO"})
	}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Numero de sanciones", Value: strconv.Itoa(user.Sanctions.Count)})
	if user.Server.Ultimatum {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ha estado en ultimatum?", Value: "SI"})
	} else {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ha estado en ultimatum?", Value: "NO"})
	}
	for _, s := range user.Sanctions.SanctionDetails {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Sancion", Value: s.String()})
	}
	me.Fields = fields
	return &me
}

func createMessageEmbedUltimatum(id string, reason string) *discordgo.MessageEmbed {
	me := discordgo.MessageEmbed{}
	me.Title = "Ultimatum"
	fields := []*discordgo.MessageEmbedField{}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Se movio a Ultimatum a:", Value: "<@" + id + ">"})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Razón", Value: reason})
	me.Fields = fields
	return &me
}

func createMessageEmbedSancionFuerte(id string, fullName string, reason string, n int) *discordgo.MessageEmbed {
	me := discordgo.MessageEmbed{}
	me.Title = "Sancion fuerte aplicada"
	fields := []*discordgo.MessageEmbedField{}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Sancionado:", Value: fullName})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "ID:", Value: id})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Razón", Value: reason})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Acción", Value: "Movido a Ultimatum, numero de sanciones incrementado"})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Sanciones totales", Value: strconv.Itoa(n)})
	me.Fields = fields
	return &me
}

func createMessageEmbedPrimerAviso(id string, fullName string, reason string) *discordgo.MessageEmbed {
	me := discordgo.MessageEmbed{}
	me.Title = "Primer aviso aplicado"
	fields := []*discordgo.MessageEmbedField{}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Sancionado:", Value: fullName})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "ID:", Value: id})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Razón", Value: reason})
	fields = append(fields, &discordgo.MessageEmbedField{Name: "Acción", Value: "Se aplico primer aviso"})
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
