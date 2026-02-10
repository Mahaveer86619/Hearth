# Hearth 🌋

**Real-time Log Aggregation and Analysis Pipeline**

Hearth is a specialized, self-hosted log aggregator designed for high-throughput distributed systems. It offers an efficient alternative to complex solutions like ELK or costly managed services by employing a lean, "pipe-centric" architecture.

## 🚀 Executive Summary

Hearth functions as a robust "black box" flight recorder for your application stack—designed for continuous, non-intrusive operation, becoming visible only when its data is required.

It categorizes logs into two states:
- **Hot Data (Real-time Stream):** Ephemeral, real-time data for immediate debugging and monitoring.
- **Cold Data (Archived Storage):** Historical, compressed data for audit trails, compliance, and deep forensic analysis.

## 🧠 Core Philosophy

- **Application Resilience First:** Prioritizes host application stability. In extreme load, Hearth sheds log data rather than introducing latency or crashes.
- **Efficiency Over Complexity:** Focused on fast processing and retrieval. High-speed `grep` over heavy SQL-like query engines.
- **Optimized I/O:** accepts that deep historical searches may involve decompressing specific data chunks, optimizing instead for rapid writes and high compression.

## 🏗️ System Architecture

Hearth operates on a **"Push & Fan-out"** logic. Applications push data to Hearth, which fans it out to storage (Cold Path) and live viewers (Hot Path).

### A. Ingestion Layer (The "Furnace")
- **TCP Port 4050:** Raw TCP entry point for efficiency.
- **Normalization:** Automatically parses JSON or wraps raw text in a JSON envelope with timestamps and metadata.

### B. The Hot Path (Live Streaming)
- **The Broker (Redis):** Immediate fan-out via Redis Pub/Sub.
- **WebSocket Gateway (HTTP Port 4040):** CLI and GUI clients connect here for zero-latency visibility.
- **Server-Side Filtering:** In-memory regex filtering on the Redis stream.

### C. The Cold Path (Archival)
- **The Buffer (WAL):** Local binary Write-Ahead Log for durability.
- **Storage (MinIO):** Rotated files are compressed with **Zstandard (Zstd)** and uploaded to S3-compatible storage.
- **Deterministic Indexing:** `bucket/service/YYYY/MM/DD/HH/chunk_id.zst`

## 🛠️ Tech Stack

- **Core:** Go (Golang)
- **Stream Broker:** Redis
- **Cold Storage:** MinIO (S3-Compatible)
- **Compression:** Zstandard (Zstd)
- **Frontend/CLI:** Flutter (GUI) / Go (TUI) - *Planned*

## 🚦 Getting Started

### Prerequisites
- Docker & Docker Compose

### Running the Stack
```bash
docker-compose up --build
```

### Components & Ports
- **Hearth Core Ingestion (TCP):** `4050`
- **Hearth Core API/WebSocket (HTTP):** `4040`
- **Hearth Test Client (Web UI):** `3000`
- **MinIO Console:** `9001`
- **Redis:** `6379`

## 🧪 Testing Ingestion

You can use the built-in **Hearth Test Client** by navigating to `http://localhost:3000` in your browser. It provides a simple UI to verify both HTTP health endpoints and TCP ingestion connectivity.

Alternatively, test raw TCP ingestion via `telnet` or `nc`:
```bash
echo '{"level":"info","msg":"Hello from the furnace"}' | nc localhost 4050
```

## 📝 Integration Guide: The "Dual Writer" Pattern

Applications should inject a custom logger writer (e.g., Go `zap` or `logrus` hook):
1. **Primary (Non-Blocking TCP):** Attempts to write to `localhost:4050`. Must be non-blocking to protect application performance.
2. **Secondary (Fallback):** `os.Stdout` to ensure `docker logs` and standard orchestration logging still function.
