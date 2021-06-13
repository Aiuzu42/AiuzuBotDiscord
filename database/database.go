package database

import (
	"strconv"
	"time"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
	"github.com/aiuzu42/AiuzuBotDiscord/models"
)

const (
	UserNotFoundCode          = 1
	UserAlredyExistsCode      = 2
	CantConnectToDatabaseCode = 3
	DatabaseErrorCode         = 4
	WrongParametersCode       = 5
	DecodingErrorCode         = 6
	UserNotFoundMessage       = "User not found"
	UserAlredyExistsMessage   = "User alredy Exists"
	WrongParametersMessage    = "Wrong parameters"
)

//Database defines the methods required for a struct to be a valid database connection.
type Databse interface {
	InitDB(c config.DBConnection) *dBError
	GetUser(userID string, username string) (models.User, *dBError)
	AddUser(user models.User) *dBError
	IncreaseMessageCount(userID string, xp int) (models.User, *dBError)
	AddJoinDate(userID string, date time.Time) *dBError
	AddLeaveDate(userID string, date time.Time) *dBError
	IncreaseSanction(userID string, reason string, mod string, modName string, command string) (models.User, *dBError)
	UpdateUser(userID string) *dBError
	AddToUpdateQuery(t string, key string, value string)
	ClearUpdateQuery()
	ModifyVxp(userID string, vxp int) (int, *dBError)
	SetVxp(userID string, vxp int) *dBError
}

type dBError struct {
	Code    int
	Message string
}

func (a dBError) String() string {
	return strconv.Itoa(a.Code) + ": " + a.Message
}
