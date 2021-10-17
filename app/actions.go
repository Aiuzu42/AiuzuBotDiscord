package app

import (
	"fmt"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (a App) SendMessage(m Message) error {
	_, err := a.b.ChannelMessageSend(m.ChannelID, m.Content)
	if err != nil {
		log.Error("[SendMessage]Error sending message: " + err.Error())
	}
	return err
}

func (a App) SendMessageEmbed(m EmbedMessage) error {
	embed := discordgo.MessageEmbed{}
	if m.Title != "" {
		embed.Title = m.Title
	}
	embed.Color = m.Color
	embed.Description = m.Content
	if m.Image != "" {
		embed.Image = &discordgo.MessageEmbedImage{URL: m.Image}
	}
	if m.Thumbnail != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: m.Thumbnail}
	}
	fields := []*discordgo.MessageEmbedField{}
	for _, f := range m.Fields {
		fields = append(fields, &discordgo.MessageEmbedField{Name: f.Name, Value: f.Value, Inline: f.Inline})
	}
	embed.Fields = fields
	_, err := a.b.ChannelMessageSendEmbed(m.ChannelID, &embed)
	return err
}

func (a App) GetChannels() (ContentWrapper, error) {
	st, err := a.b.GuildChannels(config.Config.Server)
	channels := ContentWrapper{}
	if err != nil {
		log.Error("[GetChannels]Error getting channels: " + err.Error())
		return channels, err
	}
	channels = ContentWrapper{}
	channels.C = make(map[string]string)
	for _, c := range st {
		channels.C[c.ID] = c.Name
	}
	return channels, nil
}

func (a App) GetMessage(channelID string, messageID string) (ContentWrapper, error) {
	m, err := a.b.ChannelMessage(channelID, messageID)
	if err != nil {
		log.Error("[GetMessage]Error getting message: " + err.Error())
		return ContentWrapper{}, err
	}
	messageToEdit, err := a.GetChannels()
	if err != nil {
		log.Error("[GetMessage]Error getting channels: " + err.Error())
		return messageToEdit, err
	}
	messageToEdit.ChannelID = channelID
	messageToEdit.MessageID = m.ID
	if m.Embeds != nil && len(m.Embeds) > 0 {
		messageToEdit.IsEmbed = true
		embedMessage := EmbedMessage{}
		embedMessage.ChannelID = messageToEdit.ChannelID
		embedMessage.Title = m.Embeds[0].Title
		messageToEdit.HexColor = ColorToHex(m.Embeds[0].Color)
		embedMessage.Content = m.Embeds[0].Description
		for i := 0; i < 25; i++ {
			embedMessage.Fields = append(embedMessage.Fields, Field{})
		}
		for i, field := range m.Embeds[0].Fields {
			embedMessage.Fields[i] = Field{Inline: field.Inline, Name: field.Name, Value: field.Value}
		}
		if m.Embeds[0].Image != nil {
			embedMessage.Image = m.Embeds[0].Image.URL
		}
		if m.Embeds[0].Thumbnail != nil {
			embedMessage.Thumbnail = m.Embeds[0].Thumbnail.URL
		}
		messageToEdit.Emb = embedMessage
	} else {
		messageToEdit.Content = m.Content
		messageToEdit.ContentL = len(m.Content)
		messageToEdit.IsEmbed = false
	}
	return messageToEdit, nil
}

func (a App) EditMessage(channelID string, messageID string, content string) error {
	_, err := a.b.ChannelMessageEdit(channelID, messageID, content)
	if err != nil {
		log.Error("[EditMessage]Error editing message: " + err.Error())
	}
	return err
}

func (a App) EditMessageEmbed(channelID string, messageID string, msg EmbedMessage) error {
	e := discordgo.MessageEmbed{Title: msg.Title, Color: msg.Color, Description: msg.Content}
	for _, field := range msg.Fields {
		e.Fields = append(e.Fields, &discordgo.MessageEmbedField{Name: field.Name, Value: field.Value, Inline: field.Inline})
	}
	if msg.Image != "" {
		e.Image = &discordgo.MessageEmbedImage{URL: msg.Image}
	}
	if msg.Thumbnail != "" {
		e.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: msg.Thumbnail}
	}
	_, err := a.b.ChannelMessageEditEmbed(channelID, messageID, &e)
	if err != nil {
		log.Error("[EditMessageEmbed]Error editing message: " + err.Error())
	}
	return err
}

func (m EmbedMessage) Validate() error {
	titleL := len(m.Title)
	if titleL > 256 {
		return fmt.Errorf("El titulo no puede tener mas de 256 caracteres")
	}
	contentL := len(m.Content)
	if contentL > 2048 {
		return fmt.Errorf("El contenido no puede tener mas de 2048 caracteres")
	}
	fieldsL := len(m.Fields)
	if fieldsL > 25 {
		return fmt.Errorf("No puede tener mas de 25 campos")
	}
	totalL := titleL + contentL + fieldsL
	for i, field := range m.Fields {
		nameL := len(field.Name)
		valueL := len(field.Value)
		if nameL > 256 {
			return fmt.Errorf("El nombre del campo %d no puede tener mas de 256 caracteres", i)
		}
		if valueL > 1024 {
			return fmt.Errorf("El valor del campo %d no puede tener mas de 1024 caracteres", i)
		}
		totalL = totalL + nameL + valueL
		if totalL > 6000 {
			return fmt.Errorf("El total de caracteres no debe de pasar de los 6000 caracteres")
		}
	}
	return nil
}
