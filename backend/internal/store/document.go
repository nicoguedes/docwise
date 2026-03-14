package store

import (
	"context"
	"fmt"

	"github.com/viniciusguedes/docwise/backend/internal/model"
)

func (s *Store) CreateDocument(ctx context.Context, doc *model.Document) error {
	err := s.pool.QueryRow(ctx,
		`INSERT INTO documents (filename, file_size, status)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at, updated_at`,
		doc.Filename, doc.FileSize, doc.Status,
	).Scan(&doc.ID, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return fmt.Errorf("creating document: %w", err)
	}
	return nil
}

func (s *Store) GetDocument(ctx context.Context, id string) (*model.Document, error) {
	doc := &model.Document{}
	err := s.pool.QueryRow(ctx,
		`SELECT id, filename, file_size, page_count, status, created_at, updated_at
		 FROM documents WHERE id = $1`, id,
	).Scan(&doc.ID, &doc.Filename, &doc.FileSize, &doc.PageCount, &doc.Status, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting document: %w", err)
	}
	return doc, nil
}

func (s *Store) ListDocuments(ctx context.Context) ([]model.Document, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, filename, file_size, page_count, status, created_at, updated_at
		 FROM documents ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("listing documents: %w", err)
	}
	defer rows.Close()

	var docs []model.Document
	for rows.Next() {
		var doc model.Document
		if err := rows.Scan(&doc.ID, &doc.Filename, &doc.FileSize, &doc.PageCount, &doc.Status, &doc.CreatedAt, &doc.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning document: %w", err)
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func (s *Store) UpdateDocumentStatus(ctx context.Context, id string, status model.DocumentStatus, pageCount int) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE documents SET status = $1, page_count = $2, updated_at = now() WHERE id = $3`,
		status, pageCount, id,
	)
	if err != nil {
		return fmt.Errorf("updating document status: %w", err)
	}
	return nil
}

func (s *Store) DeleteDocument(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM documents WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting document: %w", err)
	}
	return nil
}
