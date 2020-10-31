package bot

import (
	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/models"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	adminRoles = []string{}
	modRoles   = []string{}
	owners     = []string{}
)

func LoadRoles() {
	owners = config.Config.Owners
	adminRoles = config.Config.Admins
	modRoles = config.Config.Mods
}

func IsOwner(userID string) bool {
	return findIfExists(userID, owners)
}

func IsAdmin(roles []string, userID string) bool {
	if IsOwner(userID) {
		return true
	}
	return arrayFindIfExists(roles, adminRoles)
}

func IsMod(roles []string, userID string) bool {
	if IsAdmin(roles, userID) {
		return true
	}
	return arrayFindIfExists(roles, modRoles)
}

func AddToAdmins(roles []string) {
	for _, r := range roles {
		if !findIfExists(r, adminRoles) {
			adminRoles = append(adminRoles, r)
		}
	}
}

func AddToMods(roles []string) {
	for _, r := range roles {
		if !findIfExists(r, modRoles) {
			modRoles = append(modRoles, r)
		}
	}
}

func RemoveFromMods(role string) bool {
	for i, e := range modRoles {
		if role == e {
			copy(modRoles[i:], modRoles[i+1:])
			modRoles[len(modRoles)-1] = ""
			modRoles = modRoles[:len(modRoles)-1]
			return true
		}
	}
	return false
}

func RemoveFromAdmins(role string) bool {
	for i, e := range adminRoles {
		if role == e {
			copy(adminRoles[i:], adminRoles[i+1:])
			adminRoles[len(adminRoles)-1] = ""
			adminRoles = adminRoles[:len(adminRoles)-1]
			return true
		}
	}
	return false
}

func ListAdminRoles(s *discordgo.Session, id string) ([]string, *models.AppError) {
	guildRoles, err := s.GuildRoles(id)
	if err != nil {
		log.Warn("ListAdminRoles " + err.Error())
		return adminRoles, &models.AppError{Code: models.DiscordError, Message: "Discord error"}
	}
	res := make([]string, len(adminRoles))
	for i, ar := range adminRoles {
		for _, gr := range guildRoles {
			if ar == gr.ID {
				res[i] = gr.Name
				break
			}
		}
		if res[i] == "" {
			res[i] = ar
		}
	}
	return res, nil
}

func ListModRoles(s *discordgo.Session, id string) ([]string, *models.AppError) {
	guildRoles, err := s.GuildRoles(id)
	if err != nil {
		log.Warn("ListModRoles " + err.Error())
		return modRoles, &models.AppError{Code: models.DiscordError, Message: "Discord error"}
	}
	res := make([]string, len(modRoles))
	for i, ar := range modRoles {
		for _, gr := range guildRoles {
			if ar == gr.ID {
				res[i] = gr.Name
				break
			}
		}
		if res[i] == "" {
			res[i] = ar
		}
	}
	return res, nil
}

func arrayFindIfExists(a []string, b []string) bool {
	for _, ea := range a {
		for _, eb := range b {
			if ea == eb {
				return true
			}
		}
	}
	return false
}

func findIfExists(a string, b []string) bool {
	for _, eb := range b {
		if a == eb {
			return true
		}
	}
	return false
}
