# Docwise

A RAG (Retrieval-Augmented Generation) app that lets you chat with PDF documents using a local LLM.

## Architecture

```
React (Vite) → Go API (Chi) → Ollama (LLM + Embeddings)
                    ↕
              PostgreSQL + pgvector
```

## Tech Stack

- **Backend**: Go 1.22+ with Chi router (`backend/`)
- **Frontend**: React 18 + TypeScript + Vite (`frontend/`)
- **LLM**: Ollama running Llama 3.1 8B
- **Embeddings**: nomic-embed-text (768 dimensions) via Ollama
- **Vector Store**: PostgreSQL with pgvector extension
- **Infrastructure**: Docker Compose (Postgres only for local dev, full stack available)

## RAG Pipeline

1. **Upload**: PDF → pdftotext extracts text → chunker splits into overlapping chunks (1000 chars, 200 overlap) → nomic-embed-text generates embeddings → stored in pgvector
2. **Query**: Question → embedded → cosine similarity search (`<=>` operator) finds top 5 chunks → chunks injected into prompt → Llama 3.1 generates answer → streamed via SSE

## Running Locally

Prerequisites: Go, Node.js 20+, Docker, Ollama, poppler (`brew install ollama poppler`)

```bash
docker compose up postgres -d     # Start PostgreSQL
ollama serve                       # Start Ollama (if not running as service)
cd backend && go run ./cmd/server  # Backend on :8080 (auto-runs migrations)
cd frontend && npm run dev         # Frontend on :5173 (proxies /api to :8080)
```

Required Ollama models: `llama3.1:8b`, `nomic-embed-text`

## API Endpoints

- `GET /api/health` — health check
- `POST /api/documents` — upload PDF (multipart/form-data, field: `file`)
- `GET /api/documents` — list all documents
- `GET /api/documents/{id}` — get single document
- `DELETE /api/documents/{id}` — delete document and its chunks
- `POST /api/chat` — ask a question (JSON: `{document_id, question}`), returns SSE stream

## Key Design Decisions

- PDF text extraction uses `pdftotext` (poppler) via `os/exec` rather than a Go PDF library, for reliability
- Document processing runs in a background goroutine with `context.Background()` (not the request context)
- Chat streaming: Ollama returns NDJSON → backend re-packages as SSE → frontend reads with ReadableStream API
- Frontend polls `GET /api/documents` every 3s to detect status changes (pending → processing → ready)
- Migrations run automatically on server startup via `store.RunMigrations()`
- No authentication — single-user local app

## Commands

- `make up` / `make down` — Docker Compose full stack
- `make dev-backend` / `make dev-frontend` — local dev
- `make ollama-pull` — pull required models
- `make test` — run Go tests
