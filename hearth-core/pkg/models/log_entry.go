package models

import "time"

type LogEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	Service    string    `json:"service"`
	RawMessage string    `json:"raw_message"`
	Type       string    `json:"type"`
}
