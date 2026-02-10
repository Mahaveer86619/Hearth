package logger

import (
	"fmt"
	"log"
)

// Info logs an informational message with the service name prefix
// Format: 2026/02/10 06:43:14 >[Service] -> [Message]
func Info(service string, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Printf("-> %s :> %s", service, msg)
}

// Warn logs a warning message with the service name prefix
func Warn(service string, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Printf("-> %s :> %s", service, msg)
}

// Error logs an error message with the service name prefix
func Error(service string, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Printf("-> %s :> %s", service, msg)
}
