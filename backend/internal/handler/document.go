package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/viniciusguedes/docwise/backend/internal/model"
	"github.com/viniciusguedes/docwise/backend/internal/service"
	"github.com/viniciusguedes/docwise/backend/internal/store"
)

type DocumentHandler struct {
	store      *store.Store
	ragService *service.RAGService
	uploadDir  string
}

func NewDocumentHandler(store *store.Store, ragService *service.RAGService, uploadDir string) *DocumentHandler {
	return &DocumentHandler{
		store:      store,
		ragService: ragService,
		uploadDir:  uploadDir,
	}
}

func (h *DocumentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// Limit upload size to 50MB
	r.ParseMultipartForm(50 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	if filepath.Ext(header.Filename) != ".pdf" {
		http.Error(w, "Only PDF files are accepted", http.StatusBadRequest)
		return
	}

	// Create document record
	doc := &model.Document{
		Filename: header.Filename,
		FileSize: header.Size,
		Status:   model.StatusPending,
	}

	if err := h.store.CreateDocument(r.Context(), doc); err != nil {
		log.Printf("Error creating document: %v", err)
		http.Error(w, "Failed to create document record", http.StatusInternalServerError)
		return
	}

	// Save file to disk
	filePath := filepath.Join(h.uploadDir, fmt.Sprintf("%s.pdf", doc.ID))
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("Error saving file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Process document in background with a detached context
	go h.ragService.ProcessDocument(context.Background(), doc.ID, filePath)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(doc)
}

func (h *DocumentHandler) List(w http.ResponseWriter, r *http.Request) {
	docs, err := h.store.ListDocuments(r.Context())
	if err != nil {
		log.Printf("Error listing documents: %v", err)
		http.Error(w, "Failed to list documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

func (h *DocumentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	doc, err := h.store.GetDocument(r.Context(), id)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

func (h *DocumentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Delete file from disk
	filePath := filepath.Join(h.uploadDir, fmt.Sprintf("%s.pdf", id))
	os.Remove(filePath)

	if err := h.store.DeleteDocument(r.Context(), id); err != nil {
		log.Printf("Error deleting document: %v", err)
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
