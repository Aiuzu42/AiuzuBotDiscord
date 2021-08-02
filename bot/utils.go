package bot

import (
	"errors"
	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"strings"
	"time"

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

func calculateVxp(roles []string, ch string, dayVxp int64, vxpToday int) (int, bool) {
	if !config.Config.Vxp.Active || findIfExists(ch, config.Config.Vxp.IgnoredChannels) {
		return 0, false
	}
	max := 1
	for _, m := range config.Config.Vxp.VxpMultipliers {
		if findIfExists(m.Rol, roles) {
			if m.Mult == 0 {
				return 0, false
			}
			if m.Mult > max {
				max = m.Mult
			}
		}
	}
	now := newToday()
	if dayVxp < now {
		return max, true
	}
	if vxpToday >= config.Config.Vxp.MaxPerDay {
		return 0, false
	}
	if vxpToday+max > config.Config.Vxp.MaxPerDay {
		return config.Config.Vxp.MaxPerDay - vxpToday, false
	}
	return max, false
}

func caluclateRolUpgrade(vxp int, mult int) (string, []string, bool) {
	prev := vxp - mult
	for _, up := range config.Config.Vxp.RolUpgrades {
		if prev < up.Value && vxp >= up.Value {
			return up.Rol, up.ToDelete, true
		}
	}
	return "", nil, false
}

func setupRolUpgrade(rol string, toDelete []string, roles []string) []string {
	var toReturn []string
	found := false
	for _, r := range roles {
		if r == rol {
			found = true
		}
		if !findIfExists(r, toDelete) {
			toReturn = append(toReturn, r)
		}
	}
	if !found {
		toReturn = append(toReturn, rol)
	}
	return toReturn
}

func newToday() int64 {
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, config.Loc)
	return today.Unix()
}
