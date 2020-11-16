package bot

import (
	"time"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/models"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// NewMemberHandler is the handler that performs tasks when a new member is added to the guild.
func NewMemberHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if m.GuildID != config.Config.Server {
		return
	}
	_, appErr := repo.GetUser(m.Member.User.ID, "")
	if appErr != nil && appErr.Code == models.UserAlredyExists {
		appErr = repo.AddJoinDate(m.Member.User.ID, string(m.Member.JoinedAt))
		if appErr != nil {
			log.Error("[NewMemberHandler]Error adding join date: " + appErr.Message)
		}
	} else if appErr != nil {
		log.Error("[NewMemberHandler]Database error: " + appErr.Message)
	} else {
		user, err := memberToLocalUser(m.Member)
		if err != nil {
			log.Error("[NewMemberHandler]Unable to create user: " + err.Error())
			return
		}
		appErr := repo.AddUser(user)
		if appErr != nil {
			log.Error("[NewMemberHandler]Error adding user to database " + appErr.Message)
		}
	}
}

// MemberLeaveHandler is the handler that performs tasks when a member is removed from the guild.
func MemberLeaveHandler(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	if m.GuildID != config.Config.Server {
		return
	}
	msg := m.User.Username + "#" + m.User.Discriminator + " abandonó el servidor. ID: " + m.User.ID
	user, appErr := repo.GetUser(m.Member.User.ID, "")
	if appErr != nil && appErr.Code == models.UserAlredyExists {
		appErr = repo.AddLeaveDate(m.Member.User.ID, time.Now().Format(time.RFC822))
		if appErr != nil {
			log.Warn("[MemberLeaveHandler]" + msg)
			log.Error("[MemberLeaveHandler]Error adding leave date: " + appErr.Message)
		}
	} else if appErr != nil {
		log.Warn("[MemberLeaveHandler]" + msg)
		log.Error("[MemberLeaveHandler]Database error getting user: " + appErr.Message)
	} else {
		log.Info("[MemberLeaveHandler]Adding member that wasnt in DB and leave the server.")
		user, err := memberToLocalUser(m.Member)
		if err != nil {
			log.Error("[MemberLeaveHandler]Unable to create user: " + err.Error())
			return
		}
		user.Server.LeftDates = append(user.Server.LeftDates, time.Now().Format(time.RFC822))
		appErr := repo.AddUser(user)
		if appErr != nil {
			log.Error("[MemberLeaveHandler]Database error adding user: " + appErr.Message + msg)
		}
	}
	if user.Server.Ultimatum {
		msg = msg + " y era ULTIMATUM."
	}
	_, err := s.ChannelMessageSend(config.Config.FChannel, msg)
	if err != nil {
		log.Warn("[MemberLeaveHandler]" + msg)
		log.Error("[MemberLeaveHandler]Error sending message: " + err.Error() + msg)
	}
}
