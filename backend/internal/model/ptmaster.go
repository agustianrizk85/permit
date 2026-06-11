package model

import (
	"time"

	"gorm.io/datatypes"
)

// PTMaster is reusable company (PT) master data (Proses E). It is shared across
// projects and carries its own legal documents.
type PTMaster struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:160;not null" json:"name"`
	NPWP      string    `gorm:"size:32" json:"npwp"`
	NIB       string    `gorm:"size:32" json:"nib"`
	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Documents []PTDocument `gorm:"foreignKey:PTMasterID" json:"documents,omitempty"`
}

// PTDocument is a file attached to a PT (Akta PT, SK Kemenkumham, NIB, KTP &
// NPWP pemilik lahan, NPWP PT).
type PTDocument struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	PTMasterID   uint           `gorm:"index;not null" json:"pt_master_id"`
	DocType      string         `gorm:"size:48;not null" json:"doc_type"`
	OriginalName string         `gorm:"size:255;not null" json:"original_name"`
	StoredName   string         `gorm:"size:255;not null" json:"stored_name"`
	Path         string         `gorm:"size:512;not null" json:"-"`
	MimeType     string         `gorm:"size:128" json:"mime_type"`
	SizeBytes    int64          `json:"size_bytes"`
	OCRData      datatypes.JSON `gorm:"type:jsonb" json:"ocr_data"`
	UploadedBy   uint           `json:"uploaded_by"`
	CreatedAt    time.Time      `json:"created_at"`
}

// PTDocTypes lists the expected document slots for a PT (Proses E).
var PTDocTypes = []string{
	"Akta PT", "SK Kemenkumham", "NIB",
	"KTP Pemilik Lahan", "NPWP Pemilik Lahan", "NPWP PT",
}
