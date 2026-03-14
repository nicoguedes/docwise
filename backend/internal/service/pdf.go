package service

import (
	"fmt"
	"os/exec"
	"strings"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

type PDFResult struct {
	Text      string
	PageCount int
	Pages     []string
}

// ExtractText uses pdftotext (from poppler-utils) to extract text from a PDF.
// In Docker, poppler-utils is installed in the image.
// For local dev, install with: brew install poppler (macOS) or apt install poppler-utils (Linux)
func (s *PDFService) ExtractText(filePath string) (*PDFResult, error) {
	// Extract text using pdftotext
	cmd := exec.Command("pdftotext", "-layout", filePath, "-")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("running pdftotext: %w (make sure poppler-utils is installed)", err)
	}

	text := string(output)

	// Split by form feed character (page separator in pdftotext output)
	pages := strings.Split(text, "\f")
	// Remove empty last page if present
	if len(pages) > 0 && strings.TrimSpace(pages[len(pages)-1]) == "" {
		pages = pages[:len(pages)-1]
	}

	pageCount := len(pages)
	if pageCount == 0 {
		pages = []string{text}
		pageCount = 1
	}

	return &PDFResult{
		Text:      text,
		PageCount: pageCount,
		Pages:     pages,
	}, nil
}
