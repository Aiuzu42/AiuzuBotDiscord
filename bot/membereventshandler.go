package bot

import (
	"time"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	db "github.com/aiuzu42/AiuzuBotDiscord/database"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// NewMemberHandler is the handler that performs tasks when a new member is added to the guild.
func NewMemberHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if m.GuildID != config.Config.Server {
		return
	}
	_, dbErr := repo.GetUser(m.Member.User.ID, "")
	if dbErr != nil && dbErr.Code == db.UserAlredyExistsCode {
		ja, err := m.Member.JoinedAt.Parse()
		if err != nil {
			log.Error("[NewMemberHandler]Error parsing join date for user [" + m.Member.User.ID + "]: " + err.Error())
			return
		}
		_, dbErr = repo.AddJoinDate(m.Member.User.ID, ja)
		if dbErr != nil {
			log.Error("[NewMemberHandler]Error adding join date: " + dbErr.Message)
		}
	} else if dbErr != nil {
		log.Error("[NewMemberHandler]Database error: " + dbErr.Message)
	} else {
		user, err := memberToLocalUser(m.Member)
		if err != nil {
			log.Error("[NewMemberHandler]Unable to create user: " + err.Error())
			return
		}
		dbErr := repo.AddUser(user)
		if dbErr != nil {
			log.Error("[NewMemberHandler]Error adding user to database " + dbErr.Message)
		}
	}
}

// MemberLeaveHandler is the handler that performs tasks when a member is removed from the guild.
func MemberLeaveHandler(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	if m.GuildID != config.Config.Server {
		return
	}
	msg := m.User.Username + "#" + m.User.Discriminator + " abandonó el servidor. ID: " + m.User.ID
	ult, dbErr := repo.AddLeaveDate(m.Member.User.ID, time.Now())
	if dbErr != nil && dbErr.Code == db.UserNotFoundCode {
		log.Info("[MemberLeaveHandler]Adding member that wasnt in DB and leave the server.")
		user, err := memberToLocalUser(m.Member)
		if err != nil {
			log.Error("[MemberLeaveHandler]Unable to create user: " + err.Error())
			return
		}
		user.Server.AppendLeftDates(time.Now())
		dbErr := repo.AddUser(user)
		if dbErr != nil {
			log.Error("[MemberLeaveHandler]Database error adding user: " + dbErr.Message + msg)
			return
		}
	} else if dbErr != nil {
		log.Warn("[MemberLeaveHandler]" + msg)
		log.Error("[MemberLeaveHandler]Database error updating user: " + dbErr.Message)
		return
	}
	if ult {
		msg = msg + " y era ULTIMATUM."
	}
	_, err := s.ChannelMessageSend(config.Config.FChannel, msg)
	if err != nil {
		log.Warn("[MemberLeaveHandler]" + msg)
		log.Error("[MemberLeaveHandler]Error sending message: " + err.Error() + msg)
	}
}
