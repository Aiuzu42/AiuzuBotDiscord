package database

import (
	"time"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/models"
)

type Memory struct {
	db []models.User
}

func (m *Memory) InitDB(c config.DBConnection) *dBError {
	return nil
}

func (m *Memory) GetUser(userID string, username string) (models.User, *dBError) {
	for _, u := range m.db {
		if userID == u.UserID || username == u.FullName {
			return u, nil
		}
	}
	return models.User{}, &dBError{Code: UserNotFoundCode, Message: UserNotFoundMessage}
}

func (m *Memory) AddUser(user models.User) *dBError {
	for _, u := range m.db {
		if user.UserID == u.UserID {
			return &dBError{Code: UserAlredyExistsCode, Message: UserAlredyExistsMessage}
		}
	}
	m.db = append(m.db, user)
	return nil
}

func (m *Memory) IncreaseMessageCount(userID string) *dBError {
	for i := range m.db {
		if userID == m.db[i].UserID {
			m.db[i].Server.MessageCount = m.db[i].Server.MessageCount + 1
			return nil
		}
	}
	return &dBError{Code: UserNotFoundCode, Message: UserNotFoundMessage}
}

func (m *Memory) AddJoinDate(userID string, date time.Time) (bool, *dBError) {
	for i := range m.db {
		if userID == m.db[i].UserID {
			m.db[i].Server.AppendJoinDate(date)
			return m.db[i].Server.Ultimatum, nil
		}
	}
	return false, &dBError{Code: UserNotFoundCode, Message: UserNotFoundMessage}
}

func (m *Memory) AddLeaveDate(userID string, date time.Time) (bool, *dBError) {
	for i := range m.db {
		if userID == m.db[i].UserID {
			m.db[i].Server.AppendLeftDates(date)
			return m.db[i].Server.Ultimatum, nil
		}
	}
	return false, &dBError{Code: UserNotFoundCode, Message: UserNotFoundMessage}
}
