package bot

import (
	"errors"
	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"strings"

	"github.com/aiuzu42/AiuzuBotDiscord/database"
)

const (
	YELLOW     = 16766208
	LIGHT_BLUE = 6410746

	DISCORD_EPOCH = 1420070400000
)

func updateUserNames(id string, n string, un string, d string) (bool, error) {
	updated := false
	u, dbErr := repo.GetUser(id, "")
	if dbErr != nil && dbErr.Code == database.UserNotFoundCode {
		return false, errors.New("[updateUserNames]User updated but not found [" + id + "]")
	} else if dbErr != nil {
		return false, errors.New("[updateUserNames]Database error getting info: " + dbErr.Message)
	}
	repo.ClearUpdateQuery()
	if n != u.Nickname {
		updated = true
		repo.AddToUpdateQuery("$set", "nickname", n)
		repo.AddToUpdateQuery("$push", "oldNicknames", u.Nickname)
	}
	full := un + "#" + d
	if full != u.FullName {
		updated = true
		repo.AddToUpdateQuery("$set", "name", un)
		repo.AddToUpdateQuery("$set", "fullName", full)
		repo.AddToUpdateQuery("$push", "oldNames", u.FullName)
	}
	if updated {
		dbErr = repo.UpdateUser(id)
		if dbErr != nil {
			return false, errors.New("[updateUserNames]Database error updating: " + dbErr.Message)
		}
	}
	return updated, nil
}

func argumentsHandler(st string) (string, string) {
	arg := ""
	msg := ""
	args := strings.Split(st, " ")
	n := len(args)
	if n == 2 {
		arg = args[1]
	} else if n > 2 {
		arg = args[1]
		msg = strings.Join(args[2:], " ")
	}
	return arg, msg
}

func saySplit(st string) string {
	msg := ""
	args := strings.Split(st, " ")
	msg = strings.Join(args[1:], " ")
	return msg
}

func calculateVxp(roles []string, ch string) int {
	if !config.Config.Vxp.Active || findIfExists(ch, config.Config.Vxp.IgnoredChannels) {
		return 0
	}
	max := 1
	for _, m := range config.Config.Vxp.VxpMultipliers {
		if findIfExists(m.Rol, roles) {
			if m.Mult == 0 {
				return 0
			}
			if m.Mult > max {
				max = m.Mult
			}
		}
	}
	return max
}
