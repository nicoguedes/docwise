package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/viniciusguedes/docwise/backend/internal/config"
	"github.com/viniciusguedes/docwise/backend/internal/handler"
	"github.com/viniciusguedes/docwise/backend/internal/middleware"
	"github.com/viniciusguedes/docwise/backend/internal/service"
	"github.com/viniciusguedes/docwise/backend/internal/store"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	// Initialize store
	db, err := store.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize services
	pdfService := service.NewPDFService()
	chunkerService := service.NewChunkerService(cfg.ChunkSize, cfg.ChunkOverlap)
	embedderService := service.NewEmbedderService(cfg.OllamaURL, cfg.EmbedModel)
	ragService := service.NewRAGService(pdfService, chunkerService, embedderService, db, cfg.OllamaURL, cfg.ChatModel)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	documentHandler := handler.NewDocumentHandler(db, ragService, "uploads")
	chatHandler := handler.NewChatHandler(ragService)

	// Setup router
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(middleware.CORS()))

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", healthHandler.Check)

		r.Post("/documents", documentHandler.Upload)
		r.Get("/documents", documentHandler.List)
		r.Get("/documents/{id}", documentHandler.Get)
		r.Delete("/documents/{id}", documentHandler.Delete)

		r.Post("/chat", chatHandler.Ask)
	})

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
