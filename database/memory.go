package database

import (
	"time"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/models"
)

type Memory struct {
	db          []models.User
	updateQuery map[string]string
	queryStatus bool
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

func (m *Memory) SetUltimatum(userID string) *dBError {
	for i := range m.db {
		if userID == m.db[i].UserID {
			m.db[i].Server.Ultimatum = true
			m.db[i].Sanctions.Aviso = true
			return nil
		}
	}
	return &dBError{Code: UserNotFoundCode, Message: UserNotFoundMessage}
}

func (m *Memory) IncreaseSanction(userID string, reason string, mod string, modName string, command string) *dBError {
	for i := range m.db {
		if userID == m.db[i].UserID {
			m.db[i].Sanctions.Count = m.db[i].Sanctions.Count + 1
			details := models.Details{AdminID: mod, AdminName: modName, Command: command, Date: time.Now().Format(time.RFC822), Notes: reason}
			m.db[i].Sanctions.SanctionDetails = append(m.db[i].Sanctions.SanctionDetails, details)
			return nil
		}
	}
	return &dBError{Code: UserNotFoundCode, Message: UserNotFoundMessage}
}

func (m *Memory) SetPrimerAviso(userID string) *dBError {
	for i := range m.db {
		if userID == m.db[i].UserID {
			m.db[i].Sanctions.Aviso = true
			return nil
		}
	}
	return &dBError{Code: UserNotFoundCode, Message: UserNotFoundMessage}
}

func (m *Memory) UpdateUser(userID string) *dBError {
	if !m.queryStatus {
		return &dBError{Code: WrongParametersCode, Message: WrongParametersMessage}
	}
	for i := range m.db {
		if userID == m.db[i].UserID {
			if val, ok := m.updateQuery["nickname"]; ok {
				m.db[i].Nickname = val
			}
			if val, ok := m.updateQuery["oldNicknames"]; ok {
				m.db[i].OldNicknames = append(m.db[i].OldNicknames, val)
			}
			if val, ok := m.updateQuery["fullName"]; ok {
				m.db[i].FullName = val
			}
			if val, ok := m.updateQuery["oldNames"]; ok {
				m.db[i].OldNames = append(m.db[i].OldNames, val)
			}
			return nil
		}
	}
	return &dBError{Code: UserNotFoundCode, Message: UserNotFoundMessage}
}

func (m *Memory) AddToUpdateQuery(t string, key string, value string) {
	m.updateQuery[key] = value
	m.queryStatus = true
}
func (m *Memory) ClearUpdateQuery() {
	m.updateQuery = make(map[string]string)
	m.queryStatus = false
}
