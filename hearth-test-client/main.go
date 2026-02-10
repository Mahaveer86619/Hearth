package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the client configuration
type Config struct {
	TargetHost     string
	TargetHTTPPort string
	TargetTCPPort  string
	ServerPort     string
}

func loadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return Config{
		TargetHost:     getEnv("TARGET_HOST", "localhost"),
		TargetHTTPPort: getEnv("TARGET_HTTP_PORT", "4040"),
		TargetTCPPort:  getEnv("TARGET_TCP_PORT", "4050"),
		ServerPort:     getEnv("PORT", "3000"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	cfg := loadConfig()

	fs := http.FileServer(http.Dir("./templates"))
	http.Handle("/", fs)

	http.HandleFunc("/api/test-http", func(w http.ResponseWriter, r *http.Request) {
		handleHTTPTest(w, cfg)
	})

	http.HandleFunc("/api/test-tcp", func(w http.ResponseWriter, r *http.Request) {
		handleTCPTest(w, cfg)
	})

	log.Printf("Test Client starting on port %s...", cfg.ServerPort)
	log.Printf("Targeting Hearth Core at %s (HTTP: %s, TCP: %s)", cfg.TargetHost, cfg.TargetHTTPPort, cfg.TargetTCPPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatal(err)
	}
}

type TestResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type HearthResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func handleHTTPTest(w http.ResponseWriter, cfg Config) {
	url := fmt.Sprintf("http://%s:%s/health", cfg.TargetHost, cfg.TargetHTTPPort)
	resp, err := http.Get(url)

	result := TestResult{}

	if err != nil {
		result.Success = false
		result.Message = "Failed to connect to HTTP endpoint"
		result.Details = err.Error()
	} else {
		defer resp.Body.Close()

		var hearthResp HearthResponse
		if err := json.NewDecoder(resp.Body).Decode(&hearthResp); err != nil {
			result.Success = false
			result.Message = "Connected but failed to parse JSON response"
			result.Details = fmt.Sprintf("Status: %d", resp.StatusCode)
		} else {
			result.Success = resp.StatusCode == 200
			result.Message = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, hearthResp.Message)
			result.Details = hearthResp.Data
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func handleTCPTest(w http.ResponseWriter, cfg Config) {
	address := fmt.Sprintf("%s:%s", cfg.TargetHost, cfg.TargetTCPPort)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)

	result := TestResult{}

	if err != nil {
		result.Success = false
		result.Message = "Failed to establish TCP connection"
		result.Details = err.Error()
	} else {
		defer conn.Close()

		// Send a test message
		message := "Test message from dummy client"
		_, writeErr := fmt.Fprintf(conn, message)

		if writeErr != nil {
			result.Success = false
			result.Message = "Connected but failed to send data"
			result.Details = writeErr.Error()
		} else {
			result.Success = true
			result.Message = "Successfully connected and sent data"
			result.Details = fmt.Sprintf("Sent: %s to %s", message, address)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
