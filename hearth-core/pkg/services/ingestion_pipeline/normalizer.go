package ingestion_pipeline

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/Mahaveer86619/Hearth/pkg/models"
)

var (
	// Matches standard log formats like: 2026/02/10 19:13:06 INFO [SERVICE] Message
	// Or: 2026/02/10 19:13:06 INFO Message
	reStdLog = regexp.MustCompile(`^(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2})\s+([A-Z]+)\s+(.*)`)
)

// Normalize attempts to convert a raw byte slice into a structured LogEntry.
func Normalize(raw []byte) models.LogEntry {
	rawStr := strings.TrimSpace(string(raw))
	entry := models.LogEntry{
		Timestamp: time.Now(),
		Service:   "unknown",
		Severity:  "INFO",
		Type:      "raw",
		Raw:       rawStr,
		Metadata:  make(map[string]interface{}),
	}

	if rawStr == "" {
		return entry
	}

	// Strategy 1: JSON
	if strings.HasPrefix(rawStr, "{") && strings.HasSuffix(rawStr, "}") {
		var data map[string]any
		if err := json.Unmarshal(raw, &data); err == nil {
			entry.Type = "json"
			entry.Metadata = data
			extractGenericFields(&entry, data)
			return entry
		}
	}

	// Strategy 2: Standard Log Format (Regex)
	if matches := reStdLog.FindStringSubmatch(rawStr); len(matches) > 3 {
		entry.Type = "std"
		
		// Parse Time
		if t, err := time.Parse("2006/01/02 15:04:05", matches[1]); err == nil {
			entry.Timestamp = t
		}

		// Severity
		entry.Severity = strings.ToUpper(matches[2])
		
		// Message & Service Extraction
		fullMsg := matches[3]
		entry.Message = fullMsg
		
		// Try to see if there's a [SERVICE] or SERVICE: prefix
		if strings.HasPrefix(fullMsg, "[") {
			endIdx := strings.Index(fullMsg, "]")
			if endIdx > 0 {
				entry.Service = fullMsg[1:endIdx]
				entry.Message = strings.TrimSpace(fullMsg[endIdx+1:])
			}
		} else if parts := strings.SplitN(fullMsg, " ", 2); len(parts) > 1 && strings.HasSuffix(parts[0], ":") {
			entry.Service = strings.TrimSuffix(parts[0], ":")
			entry.Message = parts[1]
		}

		return entry
	}

	// Fallback: Generic Raw
	entry.Message = rawStr
	extractFromRaw(&entry, rawStr)

	return entry
}

func extractFromRaw(entry *models.LogEntry, raw string) {
	upperRaw := strings.ToUpper(raw)
	if strings.Contains(upperRaw, "ERROR") || strings.Contains(upperRaw, "ERRO") {
		entry.Severity = "ERROR"
	} else if strings.Contains(upperRaw, "WARN") {
		entry.Severity = "WARN"
	} else if strings.Contains(upperRaw, "DEBUG") {
		entry.Severity = "DEBUG"
	}
}

func extractGenericFields(entry *models.LogEntry, data map[string]interface{}) {
	msgFields := []string{"message", "msg", "content", "text", "body"}
	for _, f := range msgFields {
		if val, ok := data[f].(string); ok {
			entry.Message = val
			break
		}
	}
	serviceFields := []string{"service", "app", "application", "name", "service_name"}
	for _, f := range serviceFields {
		if val, ok := data[f].(string); ok {
			entry.Service = val
			break
		}
	}
	levelFields := []string{"level", "severity", "log_level", "type"}
	for _, f := range levelFields {
		if val, ok := data[f].(string); ok {
			entry.Severity = strings.ToUpper(val)
			break
		}
	}
	timeFields := []string{"timestamp", "time", "@timestamp", "ts", "created_at"}
	for _, f := range timeFields {
		if val := data[f]; val != nil {
			if s, ok := val.(string); ok {
				formats := []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05", "2006/01/02 15:04:05"}
				for _, fmt := range formats {
					if t, err := time.Parse(fmt, s); err == nil {
						entry.Timestamp = t
						break
					}
				}
			}
			if n, ok := val.(float64); ok {
				if n > 1e12 { entry.Timestamp = time.UnixMilli(int64(n)) } else { entry.Timestamp = time.Unix(int64(n), 0) }
			}
			break
		}
	}
	if entry.Message == "" { entry.Message = entry.Raw }
}
