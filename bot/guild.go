package bot

import (
	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/bwmarrin/discordgo"
)

var (
	adminRoles = []string{}
	modRoles   = []string{}
	owners     = []string{}
)

// LoadRoles saves to memory who are the owners, admins and mods from the configuration file.
func LoadRoles() {
	owners = config.Config.Owners
	adminRoles = config.Config.Admins
	modRoles = config.Config.Mods
}

// IsOwner returns true if userID matches an ID of the owners group.
func IsOwner(userID string) bool {
	return findIfExists(userID, owners)
}

// IsAdmin returns true if any of the roles is part of the admins group or if the userID matches an ID of the owners group.
func IsAdmin(roles []string, userID string) bool {
	if IsOwner(userID) {
		return true
	}
	return arrayFindIfExists(roles, adminRoles)
}

// IsMod returns true if any of the roles is part of the admins or mods group or if the userID matches an ID of the owners group.
func IsMod(roles []string, userID string) bool {
	if IsAdmin(roles, userID) {
		return true
	}
	return arrayFindIfExists(roles, modRoles)
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

func Is_Q_A_B(roles []string) bool {
	return arrayFindIfExists(roles, []string{config.Config.Roles.Q, config.Config.Roles.A, config.Config.Roles.B})
}

func GetMemberInfo(s *discordgo.Session, userID string) (*discordgo.Member, error) {
	return s.GuildMember(config.Config.Server, userID)
}

func DowngradeToC(roles []string) []string {
	newRoles := []string{}
	for _, role := range roles {
		if role != config.Config.Roles.Q && role != config.Config.Roles.A && role != config.Config.Roles.B && role != config.Config.Roles.C {
			newRoles = append(newRoles, role)
		}
	}
	newRoles = append(newRoles, config.Config.Roles.C)
	return newRoles
}
