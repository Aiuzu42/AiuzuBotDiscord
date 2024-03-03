package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/aiuzu42/AiuzuBotDiscord/config"
)

type User struct {
	UserID       string        `bson:"userID"`
	Name         string        `bson:"name"`
	FullName     string        `bson:"fullName"`
	OldNames     []string      `bson:"oldNames,omitempty"`
	Nickname     string        `bson:"nickname"`
	OldNicknames []string      `bson:"oldNicknames,omitempty"`
	Sanctions    Sanction      `bson:"sanctions"`
	Server       ServerDetails `bson:"server"`
	Vxp          int           `bson:"vxp"`
	DayVxp       int64         `bson:"dayVxp"`
	VxpToday     int           `bson:"vxpToday"`
}

type Sanction struct {
	Count           int       `bson:"count"`
	SanctionDetails []Details `bson:"sanctionDetails,omitempty"`
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
	JoinDates    []string `bson:"joinDates,omitempty"`
	LeftDates    []string `bson:"leftDates,omitempty"`
	LastMessage  string   `bson:"lastMessage"`
}

func (d Details) String() string {
	var strB strings.Builder
	strB.WriteString(d.AdminName + " con ID: " + d.AdminID)
	if d.Command != "" {
		strB.WriteString(" sanciono con el comando " + d.Command)
	}
	if d.Date != "" {
		strB.WriteString(" el dia " + d.Date)
	}
	if d.Notes != "" {
		strB.WriteString(". Notas: " + d.Notes)
	}
	return strB.String()
}

func (u User) String() string {
	j, _ := json.Marshal(u)
	return string(j)
}

func (s *ServerDetails) AppendJoinDate(date time.Time) {
	date = date.In(config.Loc)
	s.JoinDates = append(s.JoinDates, date.Format(time.RFC822))
}

func (s *ServerDetails) AppendLeftDates(date time.Time) {
	date = date.In(config.Loc)
	s.LeftDates = append(s.LeftDates, date.Format(time.RFC822))
}
