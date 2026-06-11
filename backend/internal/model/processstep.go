package model

import (
	"time"

	"gorm.io/datatypes"
)

// StepStatus is the checklist state of a single process step.
type StepStatus string

const (
	StatusPending    StepStatus = "pending"
	StatusInProgress StepStatus = "in_progress"
	StatusDone       StepStatus = "done"
)

// ProcessStep is one item in a process checklist (e.g. A1 "Cek SHM").
// The same generic shape is reused for every macro phase (A..I): some steps
// require a fixed price input, some require an SPK number, and structured
// sub-data (RT/RW, addresses, etc.) is kept in the flexible Metadata field so
// new phases can be added without schema changes.
type ProcessStep struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProjectID uint   `gorm:"index;not null" json:"project_id"`
	Code      string `gorm:"size:8;index;not null" json:"code"`     // e.g. "A1"
	Category  string `gorm:"size:2;index;not null" json:"category"` // e.g. "A"
	Name      string `gorm:"size:200;not null" json:"name"`
	Sequence  int    `gorm:"not null" json:"sequence"`

	Status StepStatus `gorm:"size:16;not null;default:pending" json:"status"`

	// Business rules captured from the spec.
	RequiresPrice bool   `json:"requires_price"`
	RequiresSPK   bool   `json:"requires_spk"`
	PriceLabel    string `gorm:"size:32" json:"price_label"` // e.g. "Harga Fix", "UTJ"

	// Output hints from the spec.
	NotifyDepartments  bool `json:"notify_departments"`  // C10 Siteplan -> Perencanaan & Teknik
	ConfidentialOutput bool `json:"confidential_output"` // C12/D1/D2 -> watermark for Sales

	// SLA / deadline tracking (default from template, overridable by KADEP).
	SLADays int        `json:"sla_days"`
	DueDate *time.Time `json:"due_date"`

	// Captured values.
	PriceFix  int64  `json:"price_fix"` // rupiah, integer to avoid float rounding
	SPKNumber string `gorm:"size:64" json:"spk_number"`
	Notes     string `gorm:"type:text" json:"notes"`

	// Flexible structured data per step (e.g. A5: rt_rw, alamat, tanggal).
	Metadata datatypes.JSON `gorm:"type:jsonb" json:"metadata"`

	CompletedBy *uint      `json:"completed_by"`
	CompletedAt *time.Time `json:"completed_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Documents []Document `gorm:"foreignKey:ProcessStepID" json:"documents,omitempty"`
}
