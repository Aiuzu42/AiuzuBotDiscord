package models

import (
	"encoding/json"
	"strconv"
)

const (
	UserNotFoundCode      = 1
	UserAlredyExists      = 2
	CantConnectToDatabase = 3
	DatabaseError         = 4
	DiscordError          = 5
)

type User struct {
	UserID       string        `bson:"userID"`
	Name         string        `bson:"name"`
	FullName     string        `bson:"fullName"`
	Nickname     string        `bson:"nickname"`
	OldNicknames []string      `bson:"oldNicknames"`
	Sanctions    Sanction      `bson:"sanctions"`
	Server       ServerDetails `bson:"server"`
}

type Sanction struct {
	Count           int       `bson:"count"`
	Aviso           bool      `bson:"aviso"`
	SanctionDetails []Details `bson:"sanctionDetails"`
}

type Details struct {
	AdminID   string `bson:"adminID"`
	AdminName string `bson:"adminName"`
	Command   string `bson:"command"`
	Date      string `bson:"date"`
	Notes     string `bson:"notes"`
}

type ServerDetails struct {
	MessageCount int      `bson:"messageCount"`
	JoinDates    []string `bson:"joinDates"`
	LeftDates    []string `bson:"leftDates"`
	Ultimatum    bool     `bson:"ultimatum"`
	WasModerator bool     `bson:"wasModerator"`
}

type AppError struct {
	Code    int
	Message string
}

func (d Details) String() string {
	return d.AdminName + " con ID: " + d.AdminID + " sanciono con el comando " + d.Command + " el dia " + d.Date
}

func (a AppError) String() string {
	return strconv.Itoa(a.Code) + ": " + a.Message
}

func (u User) String() string {
	j, _ := json.Marshal(u)
	return string(j)
}
