package bot

import (
	"errors"
	"strings"

	"github.com/aiuzu42/AiuzuBotDiscord/database"
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
