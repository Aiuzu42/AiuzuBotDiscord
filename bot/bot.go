package bot

import (
	"errors"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	db "github.com/aiuzu42/AiuzuBotDiscord/database"
)

var (
	repo   db.Databse
	prefix = "ai!"
	pLen   = 3
)

// SelectRepository starts and instance of the repository of the selected type.
// Accepted types "memory", "mongoDB"
func SelectRepository(dbType string) error {
	switch dbType {
	case "memory":
		repo = &db.Memory{}
	case "mongoDB":
		repo = &db.MongoDB{}
	default:
		return errors.New("No database type selected or invalid type [" + dbType + "]")
	}
	appErr := repo.InitDB(config.Config.DBConn)
	if appErr != nil {
		return errors.New(appErr.Message)
	}
	return nil
}
