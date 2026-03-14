.PHONY: up down dev-backend dev-frontend migrate ollama-pull test lint clean

up:
	docker compose up -d --build

down:
	docker compose down

dev-backend:
	cd backend && go run ./cmd/server

dev-frontend:
	cd frontend && npm run dev

migrate:
	cd backend && go run ./cmd/server --migrate

ollama-pull:
	docker compose exec ollama ollama pull llama3.1:8b
	docker compose exec ollama ollama pull nomic-embed-text

test:
	cd backend && go test ./...

lint:
	cd backend && golangci-lint run
	cd frontend && npm run lint

clean:
	docker compose down -v
	rm -f backend/uploads/*.pdf
