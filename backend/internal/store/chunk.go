package store

import (
	"context"
	"fmt"

	pgvector "github.com/pgvector/pgvector-go"
	"github.com/viniciusguedes/docwise/backend/internal/model"
)

func (s *Store) InsertChunks(ctx context.Context, chunks []model.Chunk) error {
	for _, chunk := range chunks {
		_, err := s.pool.Exec(ctx,
			`INSERT INTO chunks (document_id, content, chunk_index, page_number, embedding)
			 VALUES ($1, $2, $3, $4, $5)`,
			chunk.DocumentID, chunk.Content, chunk.ChunkIndex, chunk.PageNumber,
			pgvector.NewVector(chunk.Embedding),
		)
		if err != nil {
			return fmt.Errorf("inserting chunk %d: %w", chunk.ChunkIndex, err)
		}
	}
	return nil
}

type ChunkResult struct {
	Content    string  `json:"content"`
	Similarity float64 `json:"similarity"`
}

func (s *Store) FindSimilarChunks(ctx context.Context, documentID string, queryEmbedding []float32, limit int) ([]ChunkResult, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT content, 1 - (embedding <=> $1) AS similarity
		 FROM chunks
		 WHERE document_id = $2
		 ORDER BY embedding <=> $1
		 LIMIT $3`,
		pgvector.NewVector(queryEmbedding), documentID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("finding similar chunks: %w", err)
	}
	defer rows.Close()

	var results []ChunkResult
	for rows.Next() {
		var r ChunkResult
		if err := rows.Scan(&r.Content, &r.Similarity); err != nil {
			return nil, fmt.Errorf("scanning chunk result: %w", err)
		}
		results = append(results, r)
	}
	return results, nil
}
