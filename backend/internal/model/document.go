package model

import "time"

type DocumentStatus string

const (
	StatusPending    DocumentStatus = "pending"
	StatusProcessing DocumentStatus = "processing"
	StatusReady      DocumentStatus = "ready"
	StatusError      DocumentStatus = "error"
)

type Document struct {
	ID        string         `json:"id"`
	Filename  string         `json:"filename"`
	FileSize  int64          `json:"file_size"`
	PageCount int            `json:"page_count"`
	Status    DocumentStatus `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
