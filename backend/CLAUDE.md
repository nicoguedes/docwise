# Backend

Go backend following standard project layout.

## Structure

- `cmd/server/main.go` — entry point, wires dependencies, sets up Chi router, graceful shutdown
- `internal/config/` — env-based configuration with defaults
- `internal/model/` — data types (Document, Chunk, ChatRequest, TextChunk)
- `internal/store/` — PostgreSQL repository (pgx connection pool, pgvector queries)
- `internal/service/` — business logic:
  - `pdf.go` — extracts text via `pdftotext` CLI (poppler)
  - `chunker.go` — splits text into overlapping chunks respecting sentence boundaries
  - `embedder.go` — calls Ollama `/api/embeddings` endpoint
  - `rag.go` — orchestrates the full pipeline (process document + ask question with streaming)
- `internal/handler/` — HTTP handlers (document CRUD, chat with SSE streaming, health)
- `internal/middleware/` — CORS configuration
- `migrations/` — SQL migration files (run automatically on startup)

## Dependencies

- `go-chi/chi/v5` — HTTP router
- `go-chi/cors` — CORS middleware
- `jackc/pgx/v5` — PostgreSQL driver
- `pgvector/pgvector-go` — pgvector type support for Go

## Database

PostgreSQL with pgvector. Two tables:
- `documents` — id (uuid), filename, file_size, page_count, status, timestamps
- `chunks` — id (uuid), document_id (FK), content, chunk_index, page_number, embedding (vector(768))

IVFFlat index on embedding column for cosine similarity search.

## Building

```bash
go build -o docwise ./cmd/server
```

Docker: multi-stage build (golang:1.22-alpine → alpine:3.20 with poppler-utils).
