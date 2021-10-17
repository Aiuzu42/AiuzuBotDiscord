package app

import "github.com/bwmarrin/discordgo"

type App struct {
	b *discordgo.Session
}

type ServerError struct {
	Message string
	Code    int
}

type ContentWrapper struct {
	C         map[string]string
	ChannelID string
	MessageID string
	Default   string
	ErrorMsg  string
	HexColor  string
	Content   string
	ContentL  int
	IsEmbed   bool
	Emb       EmbedMessage
}

type Response struct {
	Message string
}

type Message struct {
	ChannelID string
	Content   string
}

type EmbedMessage struct {
	ChannelID string
	Content   string
	Color     int
	Title     string
	Fields    []Field
	Image     string
	Thumbnail string
}

type Field struct {
	Inline bool
	Name   string
	Value  string
}
