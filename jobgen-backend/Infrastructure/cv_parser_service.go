package infrastructure

import (
	"io"
	"os"
	"strings"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

type CVParserService interface {
	ExtractText(file io.Reader) (string, error)
}

type pdfParserService struct{}

func NewCVParserService() CVParserService {
	return &pdfParserService{}
}

func (p *pdfParserService) ExtractText(reader io.Reader) (string, error) {
	tempFile, err := os.CreateTemp("", "upload-*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, reader); err != nil {
		return "", err
	}

	f, err := os.Open(tempFile.Name())
	if err != nil {
		return "", err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return "", err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", err
	}

	var allText strings.Builder
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return "", err
		}
		ex, err := extractor.New(page)
		if err != nil {
			return "", err
		}
		pageText, err := ex.ExtractText()
		if err != nil {
			return "", err
		}
		allText.WriteString(pageText + "\n")
	}

	return allText.String(), nil
}
