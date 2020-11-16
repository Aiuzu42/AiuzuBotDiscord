package bot

import (
	"errors"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/database"
)

var (
	repo   database.Databse
	prefix = "ai!"
	pLen   = 3
)

const (
	defaultPrefix = "ai!"
)

// SelectRepository starts and instance of the repository of the selected type.
// Accepted types "memory", "mongoDB"
func SelectRepository(dbType string) error {
	switch dbType {
	case "memory":
		repo = &database.Memory{}
	case "mongoDB":
		repo = &database.MongoDB{}
	default:
		return errors.New("No database type selected or invalid type [" + dbType + "]")
	}
	appErr := repo.InitDB(config.Config.DBConn)
	if appErr != nil {
		return errors.New(appErr.Message)
	}
	return nil
}
