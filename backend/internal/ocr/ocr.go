// Package ocr provides a pluggable OCR/AI extraction layer. The spec asks for
// "otomatisasi data dari input/upload dokumen (KTP, dsb) -> OCR AI". A real
// provider (e.g. Claude vision, Google Document AI) implements Provider; the
// default MockProvider lets the rest of the system be built and demoed without
// external credentials.
package ocr

import (
	"context"
	"strings"
)

// Result is the structured extraction for a document.
type Result struct {
	DocType    string            `json:"doc_type"`
	Fields     map[string]string `json:"fields"`
	Confidence float64           `json:"confidence"`
	Provider   string            `json:"provider"`
}

// Provider extracts structured fields from a document's bytes.
type Provider interface {
	Name() string
	Extract(ctx context.Context, docType string, filename string, data []byte) (*Result, error)
}

// FieldSchema lists the fields we attempt to extract per document type, so the
// frontend can map them onto step metadata.
var FieldSchema = map[string][]string{
	"KTP":  {"nik", "nama", "tempat_lahir", "tanggal_lahir", "alamat", "agama", "pekerjaan"},
	"KK":   {"no_kk", "kepala_keluarga", "alamat"},
	"NPWP": {"npwp", "nama", "alamat"},
	"SHM":  {"no_shm", "luas", "atas_nama", "lokasi"},
	"PBB":  {"nop", "atas_nama", "alamat", "njop"},
}

// SchemaFor returns the field keys for a document type (empty if unknown).
func SchemaFor(docType string) []string {
	return FieldSchema[strings.ToUpper(strings.TrimSpace(docType))]
}
