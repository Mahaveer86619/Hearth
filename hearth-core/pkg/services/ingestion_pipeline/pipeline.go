package ingestion_pipeline

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Mahaveer86619/Hearth/pkg/db"
	"github.com/Mahaveer86619/Hearth/pkg/logger"
	"github.com/Mahaveer86619/Hearth/pkg/models"
)

func (s *IngestionPipelineService) Start() {
	logger.Info("Ingestion Pipeline", "Pipeline Service Started: processors ready")
	// Start multiple workers for parallel processing if needed
	go s.processLoop()
}

func (s *IngestionPipelineService) Stop() {
	close(s.quit)
}

// Ingest is called by the TCP server
func (s *IngestionPipelineService) Ingest(data []byte) {
	// Non-blocking ingest. If buffer is full, we drop logs to save the system.
	select {
	case s.ingestChan <- data:
	default:
		// TODO: Increment "dropped_logs" metric
		logger.Warn("Ingestion Pipeline", "Pipeline Buffer Full! Dropping log.")
	}
}

func (s *IngestionPipelineService) processLoop() {
	for {
		select {
		case <-s.quit:
			return
		case rawLine := <-s.ingestChan:
			s.processLog(rawLine)
		}
	}
}

func (s *IngestionPipelineService) processLog(raw []byte) {
	// 1. NORMALIZE (Phase 2 placeholder)
	// For now, we just wrap it simply
	entry := models.LogEntry{
		Timestamp:  time.Now(),
		Service:    "unknown",
		RawMessage: string(raw),
		Type:       "raw",
	}

	// 2. SERIALIZE
	data, err := json.Marshal(entry)
	if err != nil {
		logger.Error("Ingestion Pipeline", "Failed to marshal log: %v", err)
		return
	}

	// 3. BROADCAST (Hot Path)
	// Fire and forget to Redis
	if db.RDB != nil {
		db.PublishLog(context.Background(), data)
	}

	// 4. ARCHIVE (Cold Path - Phase 3)
	// s.archiver.Write(data)
}
