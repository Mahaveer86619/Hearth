package ingestion_pipeline

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/Mahaveer86619/Hearth/pkg/models"
)

var (
	// Matches standard Go log: 2026/02/10 04:57:00 INFO Message...
	// Group 1: Timestamp
	// Group 2: Level (INFO, WARN, ERROR)
	// Group 3: The rest of the message (Service + Body)
	reGoLog = regexp.MustCompile(`^(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2})\s+(\w+)?\s*(.*)`)

	// Matches SQL Latency: [422.554ms] [rows:1] ...
	reSqlTrace = regexp.MustCompile(`^\[(\d+\.\d+ms)\] \[rows:(\d+)\] (.*)`)

	// Matches OCR: OCR response for <ID>: { ... }
	reOcr = regexp.MustCompile(`OCR response for ([A-Z0-9]+): \{(.+)\}`)
)

func Normalize(raw []byte) models.LogEntry {
	rawStr := string(raw)
	entry := models.LogEntry{
		Timestamp: time.Now(),
		Service:   "SYSTEM",
		Severity:  "INFO",
		Type:      "raw",
		Raw:       rawStr,
		Metadata:  make(map[string]interface{}),
	}

	// --- Strategy 1: JSON Logs (HTTP) ---
	if strings.HasPrefix(strings.TrimSpace(rawStr), "{") {
		var jsonMap map[string]interface{}
		if err := json.Unmarshal(raw, &jsonMap); err == nil {
			entry.Type = "json"
			entry.Service = "API"

			if val, ok := jsonMap["time"].(string); ok {
				if t, err := time.Parse(time.RFC3339Nano, val); err == nil {
					entry.Timestamp = t
				}
			}
			if msg, ok := jsonMap["uri"].(string); ok {
				entry.Message = "HTTP " + msg
			}
			if status, ok := jsonMap["status"].(float64); ok {
				entry.Metadata["status"] = int(status)
				if status >= 400 {
					entry.Severity = "WARN"
				}
				if status >= 500 {
					entry.Severity = "ERROR"
				}
			}
			if lat, ok := jsonMap["latency_human"].(string); ok {
				entry.Metadata["latency"] = lat
			}
			return entry
		}
	}

	// --- Strategy 2: Standard Go Logs ---
	// Format: "2026/02/10 10:00:00 INFO Message..."
	if matches := reGoLog.FindStringSubmatch(rawStr); len(matches) > 3 {
		// 1. Parse Timestamp
		if t, err := time.Parse("2006/01/02 15:04:05", matches[1]); err == nil {
			entry.Timestamp = t
		}

		// 2. Parse Level
		if matches[2] != "" {
			entry.Severity = matches[2]
		}

		// 3. Dynamic Service Extraction
		// Logic: Take the whole message string (Group 3), split by space, take the first word.
		fullMsg := strings.TrimSpace(matches[3])
		entry.Message = fullMsg

		// Split "OCR response for..." -> ["OCR", "response", "for"...]
		// Split "FSN: TMRGT..." -> ["FSN:", "TMRGT..."]
		parts := strings.SplitN(fullMsg, " ", 2)
		if len(parts) > 0 {
			// Clean the service name (remove trailing colon if present)
			serviceName := strings.TrimSuffix(parts[0], ":")

			// Save to Metadata as requested
			entry.Metadata["service_name"] = serviceName

			// Also update the main Service field for indexing
			entry.Service = serviceName
		}

		entry.Type = "std"

		// 4. Run Specific Metadata Parsing (Deep extraction)
		extractSpecificMetadata(&entry)

		return entry
	}

	// --- Strategy 3: SQL Traces (No Timestamp header) ---
	// Format: "[400ms] [rows:1] SELECT..."
	if matches := reSqlTrace.FindStringSubmatch(rawStr); len(matches) > 3 {
		entry.Type = "sql"
		entry.Service = "DB"
		entry.Severity = "WARN"
		entry.Metadata["latency"] = matches[1]
		entry.Metadata["rows"] = matches[2]
		entry.Message = "SQL: " + matches[3]
		return entry
	}

	// Fallback for unknown lines
	entry.Message = rawStr
	return entry
}

func extractSpecificMetadata(entry *models.LogEntry) {
	// OCR Deep Parsing
	// Even though we already extracted Service="OCR", we want the Shipment ID too.
	if entry.Service == "OCR" || strings.Contains(entry.Message, "OCR response") {
		matches := reOcr.FindStringSubmatch(entry.Message)
		if len(matches) > 2 {
			entry.Metadata["shipment_id"] = matches[1]
			if strings.Contains(matches[2], "Status:success") {
				entry.Metadata["ocr_status"] = "success"
			} else {
				entry.Metadata["ocr_status"] = "failed"
				entry.Severity = "ERROR"
			}
		}
	}

	// Slow SQL Header (Standard Log format)
	if strings.Contains(entry.Message, "SLOW SQL") {
		entry.Severity = "WARN"
		entry.Metadata["alert"] = "slow_sql"
	}
}
