package bot

import (
	"time"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/models"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func NewMemberHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if m.GuildID != config.Config.Server {
		return
	}
	_, appErr := repo.GetUser(m.Member.User.ID, "")
	if appErr != nil && appErr.Code == models.UserAlredyExists {
		appErr = repo.AddJoinDate(m.Member.User.ID, string(m.Member.JoinedAt))
		if appErr != nil {
			log.Error("error in new member event1 " + appErr.Message)
		}
	} else if appErr != nil {
		log.Error("error in new member event2 " + appErr.Message)
	} else {
		full := m.Member.User.Username + "#" + m.Member.User.Discriminator
		user := models.User{UserID: m.Member.User.ID, Name: m.Member.User.Username, Nickname: m.Member.Nick, FullName: full}
		user.Server.JoinDates = append(user.Server.JoinDates, string(m.Member.JoinedAt))
		user.Server.MessageCount = 0
		user.Server.WasModerator = false
		user.Server.Ultimatum = false
		appErr := repo.AddUser(user)
		if appErr != nil {
			log.Error("error in new member event3 " + appErr.Message)
		}
	}
}

func MemberLeaveHandler(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	if m.GuildID != config.Config.Server {
		return
	}
	msg := m.User.Username + "#" + m.User.Discriminator + " abandonó el servidor. ID: " + m.User.ID
	_, err := s.ChannelMessageSend(config.Config.FChannel, msg)
	if err != nil {
		log.Error("error in member leave event1" + err.Error())
	}
	_, appErr := repo.GetUser(m.Member.User.ID, "")
	if appErr != nil && appErr.Code == models.UserAlredyExists {
		appErr = repo.AddLeaveDate(m.Member.User.ID, time.Now().Format(time.RFC3339))
		if appErr != nil {
			log.Error("error in member leave event2 " + appErr.Message)
		}
	} else if appErr != nil {
		log.Error("error in member leave event3 " + appErr.Message)
	} else {
		full := m.Member.User.Username + "#" + m.Member.User.Discriminator
		user := models.User{UserID: m.Member.User.ID, Name: m.Member.User.Username, Nickname: m.Member.Nick, FullName: full}
		user.Server.JoinDates = append(user.Server.JoinDates, string(m.Member.JoinedAt))
		user.Server.MessageCount = 0
		user.Server.WasModerator = false
		user.Server.Ultimatum = false
		user.Server.LeftDates = append(user.Server.LeftDates, time.Now().Format(time.RFC3339))
		appErr := repo.AddUser(user)
		if appErr != nil {
			log.Error("error in member leave event4 " + appErr.Message)
		}
	}
}
