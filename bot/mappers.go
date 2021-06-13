package bot

import (
	"errors"

	"github.com/aiuzu42/AiuzuBotDiscord/models"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func memberToLocalUser(m *discordgo.Member) (models.User, error) {
	user := models.User{}
	if m == nil || m.User == nil {
		log.Error("[memberToLocalUser]Mandatory data missing")
		return user, errors.New("Mandatory data missing")
	}
	if m.User.Discriminator == "0000" {
		return user, errors.New("webhook")
	}
	user.UserID = m.User.ID
	fullName := m.User.Username + "#" + m.User.Discriminator
	user.FullName = fullName
	user.Name = m.User.Username
	user.Nickname = m.Nick
	t, err := m.JoinedAt.Parse()
	if err != nil {
		log.Error("[memberToLocalUser]Unable to parse date info for " + m.User.ID)
	} else {
		user.Server.AppendJoinDate(t)
	}
	user.Vxp = 0
	return user, nil
}

func userAndMemberToLocalUser(u *discordgo.User, m *discordgo.Member) (models.User, error) {
	user := models.User{}
	if u == nil {
		log.Error("[userAndMemberToLocalUser]Mandatory data missing")
		return user, errors.New("Mandatory data missing")
	}
	if u.Discriminator == "0000" {
		return user, errors.New("webhook")
	}
	user.UserID = u.ID
	fullName := u.Username + "#" + u.Discriminator
	user.FullName = fullName
	user.Name = u.Username
	if m != nil {
		user.Nickname = m.Nick
		t, err := m.JoinedAt.Parse()
		if err != nil {
			log.Error("[userAndMemberToLocalUser]Unable to parse date info for " + m.User.ID)
		} else {
			user.Server.AppendJoinDate(t)
		}
	}
	user.Vxp = 0
	return user, nil
}
