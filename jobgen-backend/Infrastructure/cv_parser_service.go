package infrastructure

import (
	"bytes"
	"io"
	"os"
	"strings"

	pdf "github.com/ledongthuc/pdf"
)

type CVParserService interface {
	ExtractText(file io.Reader) (string, error)
}

type pdfParserService struct{}

func NewCVParserService() CVParserService {
	return &pdfParserService{}
}

func (p *pdfParserService) ExtractText(reader io.Reader) (string, error) {
	// Write the incoming reader to a temporary file so we can use the parser APIs
	tmp, err := os.CreateTemp("", "cv-*.pdf")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	defer func() {
		tmp.Close()
		os.Remove(tmpPath)
	}()

	if _, err := io.Copy(tmp, reader); err != nil {
		return "", err
	}

	// Re-open via ledongthuc/pdf and extract plain text
	f, r, err := pdf.Open(tmpPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	br, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	if _, err := b.ReadFrom(br); err != nil {
		return "", err
	}

	// Normalize line endings lightly
	text := strings.ReplaceAll(b.String(), "\r\n", "\n")
	return text, nil
}
