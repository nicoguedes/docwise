package model

type Chunk struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	Content    string    `json:"content"`
	ChunkIndex int       `json:"chunk_index"`
	PageNumber int       `json:"page_number"`
	Embedding  []float32 `json:"-"`
}

type TextChunk struct {
	Content    string
	PageNumber int
}

type ChatRequest struct {
	DocumentID string `json:"document_id"`
	Question   string `json:"question"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
