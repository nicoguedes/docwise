package service

import (
	"strings"

	"github.com/viniciusguedes/docwise/backend/internal/model"
)

type ChunkerService struct {
	chunkSize int
	overlap   int
}

func NewChunkerService(chunkSize, overlap int) *ChunkerService {
	return &ChunkerService{
		chunkSize: chunkSize,
		overlap:   overlap,
	}
}

func (s *ChunkerService) Chunk(pages []string) []model.TextChunk {
	var chunks []model.TextChunk

	fullText := strings.Join(pages, "\n")

	// Build a map of character offset -> page number
	pageOffsets := make([]int, len(pages))
	offset := 0
	for i, page := range pages {
		pageOffsets[i] = offset
		offset += len(page) + 1 // +1 for the newline
	}

	for start := 0; start < len(fullText); {
		end := start + s.chunkSize
		if end > len(fullText) {
			end = len(fullText)
		}

		// Try to break at a sentence boundary
		if end < len(fullText) {
			if idx := lastSentenceBreak(fullText[start:end]); idx > 0 {
				end = start + idx + 1
			}
		}

		chunk := strings.TrimSpace(fullText[start:end])
		if len(chunk) > 0 {
			pageNum := findPage(start, pageOffsets)
			chunks = append(chunks, model.TextChunk{
				Content:    chunk,
				PageNumber: pageNum + 1,
			})
		}

		start = end - s.overlap
		if start < 0 {
			start = 0
		}
		if end >= len(fullText) {
			break
		}
	}

	return chunks
}

func lastSentenceBreak(text string) int {
	best := -1
	for _, sep := range []string{". ", ".\n", "? ", "!\n", "! ", "?\n"} {
		if idx := strings.LastIndex(text, sep); idx > best {
			best = idx
		}
	}
	return best
}

func findPage(offset int, pageOffsets []int) int {
	page := 0
	for i, po := range pageOffsets {
		if offset >= po {
			page = i
		}
	}
	return page
}
