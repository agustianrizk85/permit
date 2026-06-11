package ocr

import (
	"context"
	"strings"
)

// MockProvider returns deterministic placeholder fields for the document's
// schema. It marks every value as "[OCR mock]" so it is never mistaken for real
// extracted data. Swap it for a real Provider by wiring one in the router.
type MockProvider struct{}

func NewMockProvider() *MockProvider { return &MockProvider{} }

func (m *MockProvider) Name() string { return "mock" }

func (m *MockProvider) Extract(_ context.Context, docType, filename string, _ []byte) (*Result, error) {
	fields := map[string]string{}
	for _, key := range SchemaFor(docType) {
		fields[key] = "[OCR mock] isi otomatis dari " + filename
	}
	if len(fields) == 0 {
		fields["raw_text"] = "[OCR mock] dokumen " + strings.TrimSpace(docType)
	}
	return &Result{
		DocType:    docType,
		Fields:     fields,
		Confidence: 0.0,
		Provider:   m.Name(),
	}, nil
}
