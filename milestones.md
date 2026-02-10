# Hearth Development Milestones

## Phase 1: The Pulse (Hot Path & Plumbing)
**Goal:** Get logs from TCP -> Redis -> WebSocket immediately.
- [ ] **1.1 Redis Infrastructure:** Implement `pkg/db/redis.go` to handle connection and Pub/Sub.
- [ ] **1.2 The Pipeline Orchestrator:** Create `pkg/pipeline/service.go`. This is the central nervous system that accepts raw bytes, parses them, and routes them.
- [ ] **1.3 TCP Wiring:** Update `pkg/web/tcp_server.go` to feed data into the Pipeline instead of just printing to stdout.
- [ ] **1.4 WebSocket Subscription:** Update `pkg/handlers/ws_handler.go` to subscribe to Redis and stream real logs to the browser.

## Phase 2: The Logic (Normalization & Parsing)
**Goal:** Turn raw text into structured JSON that we can query.
- [ ] **2.1 Log Normalizer:** Implement logic to detect if a log is JSON or Raw Text.
- [ ] **2.2 Special Parsers:** Add specific regex support for the `SLOW SQL` and `OCR response` logs seen in your project.
- [ ] **2.3 Metadata Tagging:** tag logs with `service_name`, `timestamp`, and `severity`.

## Phase 3: The Ash (Cold Path & Archival)
**Goal:** Save logs to disk and upload to MinIO without blocking.
- [ ] **3.1 WAL (Write-Ahead Log):** Implement a buffered file writer that appends logs to a local temp file.
- [ ] **3.2 The Rotator:** Create a background ticker that checks file size/age.
- [ ] **3.3 The Compressor:** Implement Zstd compression logic.
- [ ] **3.4 MinIO Uploader:** Implement `pkg/storage/minio.go` to upload compressed chunks.

## Phase 4: Client Integration
**Goal:** Connect your actual Go application (`ims-wvs-server`).
- [ ] **4.1 Go SDK:** Create a non-blocking `io.Writer` package for your app to import.