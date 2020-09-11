package models

type User struct {
	ID        string
	Name      string
	Nickname  string
	Sanctions Sanction
	Server    ServerDetails
}

type Sanction struct {
	Count           int
	Aviso           bool
	SanctionDetails Details
}

type Details struct {
	AdminID   string
	AdminName string
	Command   string
	Date      string
}

type ServerDetails struct {
	MessageCount int
	JoinDates    []string
	LeftDates    []string
	Ultimatum    bool
	HasBeenAdmin bool
}

type AppError struct {
	Code    int
	Message string
}
