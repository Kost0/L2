package models

type Event struct {
	UserID string `json:"userID"`
	Date   string `json:"date"`
	Event  string `json:"event"`
}

var Events = make([]*Event, 0)
