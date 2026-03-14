package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/viniciusguedes/docwise/backend/internal/model"
	"github.com/viniciusguedes/docwise/backend/internal/service"
)

type ChatHandler struct {
	ragService *service.RAGService
}

func NewChatHandler(ragService *service.RAGService) *ChatHandler {
	return &ChatHandler{ragService: ragService}
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func (h *ChatHandler) Ask(w http.ResponseWriter, r *http.Request) {
	var req model.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.DocumentID == "" || req.Question == "" {
		http.Error(w, "document_id and question are required", http.StatusBadRequest)
		return
	}

	// Get streaming response from Ollama
	body, err := h.ragService.Ask(r.Context(), req.DocumentID, req.Question)
	if err != nil {
		log.Printf("Error asking question: %v", err)
		http.Error(w, "Failed to process question", http.StatusInternalServerError)
		return
	}
	defer body.Close()

	// Stream response as SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		var ollamaResp ollamaGenerateResponse
		if err := json.Unmarshal(scanner.Bytes(), &ollamaResp); err != nil {
			continue
		}

		// Send SSE event
		fmt.Fprintf(w, "data: %s\n\n", mustJSON(map[string]any{
			"content": ollamaResp.Response,
			"done":    ollamaResp.Done,
		}))
		flusher.Flush()

		if ollamaResp.Done {
			break
		}
	}
}

func mustJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
