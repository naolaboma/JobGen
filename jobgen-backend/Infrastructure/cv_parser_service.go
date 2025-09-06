package infrastructure

import (
	"bytes"
	"io"
	"os"

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
	// Write incoming reader to a temp file so parsers that require a path can read it
	tempFile, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, reader); err != nil {
		return "", err
	}

	// Open the PDF using ledongthuc/pdf
	f, r, err := pdf.Open(tempFile.Name())
	if err != nil {
		return "", err
	}
	defer f.Close()

	plain, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(plain); err != nil {
		return "", err
	}
	return buf.String(), nil
}
