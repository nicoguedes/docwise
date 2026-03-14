package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/viniciusguedes/docwise/backend/internal/model"
	"github.com/viniciusguedes/docwise/backend/internal/store"
)

type RAGService struct {
	pdf      *PDFService
	chunker  *ChunkerService
	embedder *EmbedderService
	store    *store.Store
	ollamaURL string
	chatModel string
}

func NewRAGService(pdf *PDFService, chunker *ChunkerService, embedder *EmbedderService, store *store.Store, ollamaURL, chatModel string) *RAGService {
	return &RAGService{
		pdf:       pdf,
		chunker:   chunker,
		embedder:  embedder,
		store:     store,
		ollamaURL: ollamaURL,
		chatModel: chatModel,
	}
}

func (s *RAGService) ProcessDocument(ctx context.Context, docID, filePath string) {
	// Update status to processing
	if err := s.store.UpdateDocumentStatus(ctx, docID, model.StatusProcessing, 0); err != nil {
		log.Printf("Error updating document status: %v", err)
		return
	}

	// Extract text from PDF
	result, err := s.pdf.ExtractText(filePath)
	if err != nil {
		log.Printf("Error extracting text from PDF: %v", err)
		s.store.UpdateDocumentStatus(ctx, docID, model.StatusError, 0)
		return
	}

	// Chunk the text
	textChunks := s.chunker.Chunk(result.Pages)
	if len(textChunks) == 0 {
		log.Printf("No chunks generated from document %s", docID)
		s.store.UpdateDocumentStatus(ctx, docID, model.StatusError, 0)
		return
	}

	// Generate embeddings
	texts := make([]string, len(textChunks))
	for i, tc := range textChunks {
		texts[i] = tc.Content
	}

	embeddings, err := s.embedder.EmbedBatch(ctx, texts)
	if err != nil {
		log.Printf("Error generating embeddings: %v", err)
		s.store.UpdateDocumentStatus(ctx, docID, model.StatusError, 0)
		return
	}

	// Store chunks with embeddings
	chunks := make([]model.Chunk, len(textChunks))
	for i, tc := range textChunks {
		chunks[i] = model.Chunk{
			DocumentID: docID,
			Content:    tc.Content,
			ChunkIndex: i,
			PageNumber: tc.PageNumber,
			Embedding:  embeddings[i],
		}
	}

	if err := s.store.InsertChunks(ctx, chunks); err != nil {
		log.Printf("Error inserting chunks: %v", err)
		s.store.UpdateDocumentStatus(ctx, docID, model.StatusError, 0)
		return
	}

	// Update status to ready
	if err := s.store.UpdateDocumentStatus(ctx, docID, model.StatusReady, result.PageCount); err != nil {
		log.Printf("Error updating document status to ready: %v", err)
	}

	log.Printf("Document %s processed: %d pages, %d chunks", docID, result.PageCount, len(chunks))
}

type generateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

func (s *RAGService) Ask(ctx context.Context, documentID, question string) (io.ReadCloser, error) {
	// Embed the question
	queryEmbedding, err := s.embedder.Embed(ctx, question)
	if err != nil {
		return nil, fmt.Errorf("embedding question: %w", err)
	}

	// Find similar chunks
	results, err := s.store.FindSimilarChunks(ctx, documentID, queryEmbedding, 5)
	if err != nil {
		return nil, fmt.Errorf("finding similar chunks: %w", err)
	}

	// Build context from chunks
	var contextBuilder strings.Builder
	for i, r := range results {
		fmt.Fprintf(&contextBuilder, "--- Excerpt %d (similarity: %.2f) ---\n%s\n\n", i+1, r.Similarity, r.Content)
	}

	// Build prompt
	prompt := fmt.Sprintf(`You are a helpful assistant that answers questions based on the provided document excerpts. Use only the information from the excerpts to answer. If the answer is not in the excerpts, say so.

Document excerpts:
%s

Question: %s

Answer:`, contextBuilder.String(), question)

	// Call Ollama generate API with streaming
	body, err := json.Marshal(generateRequest{
		Model:  s.chatModel,
		Prompt: prompt,
		Stream: true,
	})
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.ollamaURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling Ollama generate API: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("Ollama generate API returned %d: %s", resp.StatusCode, string(respBody))
	}

	return resp.Body, nil
}
