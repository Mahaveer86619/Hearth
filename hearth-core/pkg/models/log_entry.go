package models

import "time"

type LogEntry struct {
	Timestamp time.Time      `json:"timestamp"`
	Service   string         `json:"service"`
	Severity  string         `json:"severity"` // INFO, ERROR, WARN, SQL, HTTP
	Message   string         `json:"message"`
	Type      string         `json:"type"`               // "json", "ocr", "sql", "std"
	Metadata  map[string]any `json:"metadata,omitempty"` // Flexible bag of attributes
	Raw       string         `json:"raw,omitempty"`      // Keep original for debugging
}
