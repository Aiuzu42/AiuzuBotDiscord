package bot

import (
	"strconv"

	"github.com/aiuzu42/AiuzuBotDiscord/models"
	"github.com/bwmarrin/discordgo"
)

func createMessageEmbedUser(user models.User) *discordgo.MessageEmbed {
	me := discordgo.MessageEmbed{}
	me.Title = "Details of user " + user.Name
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
	if user.Sanctions.Count > 1 {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Nivel de sancion", Value: "ST-" + strconv.Itoa(user.Sanctions.Count)})
	}
	if user.Sanctions.Aviso {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ya tiene un primer aviso"})
	}
	fields = append(fields, &discordgo.MessageEmbedField{Name: "N. de mensajes", Value: strconv.Itoa(user.Server.MessageCount)})
	if len(user.Server.JoinDates) > 1 {
		chain := ""
		for _, d := range user.Server.JoinDates {
			chain = chain + d + ", "
		}
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Entro en las fechas", Value: chain})
	}
	if len(user.Server.LeftDates) > 1 {
		chain := ""
		for _, d := range user.Server.LeftDates {
			chain = chain + d + ", "
		}
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Se ha ido en las fechas", Value: chain})
	}
	if user.Server.Ultimatum {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ha estado en ultimatum"})
	}
	me.Fields = fields
	return &me
}

func createMessageEmbedSanctions(user models.User) *discordgo.MessageEmbed {
	me := discordgo.MessageEmbed{}
	me.Title = "Sanciones del usuario " + user.Name
	fields := []*discordgo.MessageEmbedField{}
	if user.Sanctions.Aviso {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ya tiene un primer aviso"})
	}
	if user.Sanctions.Count > 1 {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Nivel de sancion", Value: "ST-" + strconv.Itoa(user.Sanctions.Count)})
	}
	if user.Server.Ultimatum {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Ha estado en ultimatum"})
	}
	for _, s := range user.Sanctions.SanctionDetails {
		fields = append(fields, &discordgo.MessageEmbedField{Name: "Sancion", Value: s.String()})
	}
	me.Fields = fields
	return &me
}
