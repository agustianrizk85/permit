package model

import "time"

// Vendor is reusable third-party (Pihak Ketiga) master data (Proses I). It is
// referenced by SPK Legal Permit (Proses J) as the executing vendor.
type Vendor struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Category      string    `gorm:"size:48;not null;default:'Legal Permit'" json:"category"`
	Name          string    `gorm:"size:160;not null" json:"name"`
	Address       string    `gorm:"type:text" json:"address"`
	KTPNumber     string    `gorm:"size:32" json:"ktp_number"`
	AccountNumber string    `gorm:"size:48" json:"account_number"`
	BankName      string    `gorm:"size:80" json:"bank_name"`
	AccountHolder string    `gorm:"size:160" json:"account_holder"`
	Notes         string    `gorm:"type:text" json:"notes"`
	CreatedBy     uint      `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// VendorCategories lists the supported vendor categories (Proses I-1).
var VendorCategories = []string{"Legal Permit"}
