package database

import (
	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/models"
)

type Databse interface {
	InitDB(c config.DBConnection) *models.AppError
	GetUser(userID string, username string) (models.User, *models.AppError)
	AddUser(user models.User) *models.AppError
	IncreaseMessageCount(userID string) *models.AppError
	AddJoinDate(userID string, date string) *models.AppError
	AddLeaveDate(userID string, date string) *models.AppError
}
