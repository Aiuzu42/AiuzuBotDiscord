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
