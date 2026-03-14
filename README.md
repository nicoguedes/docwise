# Docwise

Chat with your PDF documents using a local LLM and RAG (Retrieval-Augmented Generation).

Upload a PDF, and Docwise will chunk the text, generate embeddings, and let you ask questions about the document — all running locally on your machine, no API keys needed.

## Architecture

```
Browser (React) ──► Go API ──► Ollama (LLM + Embeddings)
                      │
                 PostgreSQL + pgvector
```

**How it works:**

1. **Upload** — PDF text is extracted, split into overlapping chunks, embedded via Ollama (nomic-embed-text), and stored in PostgreSQL with pgvector
2. **Ask** — Your question is embedded, the most similar chunks are retrieved via cosine similarity, injected into a prompt, and sent to Llama 3.1 for a streamed answer

## Tech Stack

| Layer | Technology |
|-------|-----------|
| LLM | [Ollama](https://ollama.com) + Llama 3.1 8B |
| Embeddings | nomic-embed-text (768 dimensions) |
| Vector Store | PostgreSQL + [pgvector](https://github.com/pgvector/pgvector) |
| Backend | Go + [Chi](https://github.com/go-chi/chi) |
| Frontend | React + TypeScript + Vite |
| Infrastructure | Docker Compose |

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/)
- [Docker](https://docs.docker.com/get-docker/) (for PostgreSQL)
- [Ollama](https://ollama.com/download)
- [poppler](https://poppler.freedesktop.org/) (for PDF text extraction)

### Install prerequisites on macOS

```bash
brew install ollama poppler
```

## Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/nicoguedes/docwise.git
cd docwise
```

### 2. Start PostgreSQL with pgvector

```bash
docker compose up postgres -d
```

### 3. Start Ollama and pull models

```bash
# Start Ollama (runs in background)
ollama serve &

# Pull the required models
ollama pull llama3.1:8b
ollama pull nomic-embed-text
```

### 4. Set up environment

```bash
cp .env.example .env
```

The defaults work out of the box for local development.

### 5. Run the backend

```bash
cd backend
go run ./cmd/server
```

The API will start on `http://localhost:8080`. It automatically runs database migrations on startup.

### 6. Run the frontend

In a separate terminal:

```bash
cd frontend
npm install
npm run dev
```

The UI will be available at `http://localhost:5173`.

### 7. Use it

1. Open `http://localhost:5173` in your browser
2. Upload a PDF using the sidebar
3. Wait for it to finish processing (status changes to "ready")
4. Select the document and start asking questions

## Running with Docker Compose (full stack)

To run everything in containers:

```bash
make up
make ollama-pull
```

Then open `http://localhost:3000`.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/health` | Health check |
| `POST` | `/api/documents` | Upload a PDF (multipart/form-data) |
| `GET` | `/api/documents` | List all documents |
| `GET` | `/api/documents/{id}` | Get a document |
| `DELETE` | `/api/documents/{id}` | Delete a document |
| `POST` | `/api/chat` | Ask a question (SSE streaming response) |

## Project Structure

```
docwise/
├── backend/
│   ├── cmd/server/          # Entry point
│   ├── internal/
│   │   ├── config/          # Environment-based config
│   │   ├── handler/         # HTTP handlers
│   │   ├── middleware/      # CORS
│   │   ├── model/           # Data types
│   │   ├── service/         # PDF extraction, chunking, embeddings, RAG
│   │   └── store/           # PostgreSQL + pgvector queries
│   └── migrations/          # SQL migrations
├── frontend/
│   └── src/
│       ├── api/             # API client
│       ├── components/      # React components
│       ├── hooks/           # useChat, useDocuments
│       └── pages/           # Home page
├── docker-compose.yml
└── Makefile
```

## License

MIT
