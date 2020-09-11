package database

import "github.com/aiuzu42/AiuzuBotDiscord/models"

type memory struct {
	db []models.User
}

func (m *memory) GetUser(userID string) (models.User, *models.AppError) {
	for _, u := range m.db {
		if userID == u.ID {
			return u, nil
		}
	}
	return models.User{}, &models.AppError{Code: userNotFoundCode, Message: "User not found."}
}

func (m *memory) AddUser(user models.User) *models.AppError {
	for _, u := range m.db {
		if user.ID == u.ID {
			return &models.AppError{Code: userAlredyExists, Message: "User Alredy Exists"}
		}
	}
	m.db = append(m.db, user)
	return nil
}
