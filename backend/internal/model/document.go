package model

import (
	"time"

	"gorm.io/datatypes"
)

// Document is an uploaded file attached to a process step (and project).
// Documents that must be shared with other departments (Sales, etc.) can be
// flagged Confidential so the system applies a watermark / black-and-white
// rendition on export (IMB/PBG, SHM, PBB per the spec).
type Document struct {
	ID            uint  `gorm:"primaryKey" json:"id"`
	ProjectID     uint  `gorm:"index;not null" json:"project_id"`
	ProcessStepID *uint `gorm:"index" json:"process_step_id"`

	DocType      string `gorm:"size:48;not null" json:"doc_type"` // KTP, KK, NPWP, SHM, PBB, SPK, ...
	OriginalName string `gorm:"size:255;not null" json:"original_name"`
	StoredName   string `gorm:"size:255;not null" json:"stored_name"`
	Path         string `gorm:"size:512;not null" json:"-"`
	MimeType     string `gorm:"size:128" json:"mime_type"`
	SizeBytes    int64  `json:"size_bytes"`

	Confidential bool `gorm:"default:false" json:"confidential"`

	// Result of optional OCR AI extraction (KTP fields, etc.).
	OCRData datatypes.JSON `gorm:"type:jsonb" json:"ocr_data"`

	UploadedBy uint      `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
}
