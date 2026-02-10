package models

import "time"

type LogEntry struct {
	Timestamp time.Time      `json:"timestamp"`
	Service   string         `json:"service"`  // e.g., "OCR", "BILLING", "API"
	Severity  string         `json:"severity"` // INFO, ERROR, WARN, FATAL
	Message   string         `json:"message"`
	Type      string         `json:"type"`               // "json", "ocr", "sql", "std"
	Metadata  map[string]any `json:"metadata,omitempty"` // Flexible bag for specific data (latency, rows, etc.)
	Raw       string         `json:"raw,omitempty"`      // Original string for debugging
}
