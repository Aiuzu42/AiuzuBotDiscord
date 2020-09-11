package database

import "github.com/aiuzu42/AiuzuBotDiscord/models"

const (
	userNotFoundCode = 1
	userAlredyExists = 2
)

type databse interface {
	GetUser(userID string) (models.User, *models.AppError)
	AddUser(user models.User) *models.AppError
}
