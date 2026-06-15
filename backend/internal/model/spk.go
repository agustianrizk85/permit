package model

import "time"

// SPKStatus is the approval state of an SPK (Proses J). Kadep creates a draft;
// Dirops approves or rejects from the dashboard.
type SPKStatus string

const (
	SPKDraft    SPKStatus = "draft"
	SPKApproved SPKStatus = "approved"
	SPKRejected SPKStatus = "rejected"
)

// Pricing modes determine how Total is derived from UnitPrice and Quantity.
const (
	PricingPerMeter = "per_meter" // Total = UnitPrice * Quantity (luas m²)
	PricingPerUnit  = "per_unit"  // Total = UnitPrice * Quantity (jumlah unit)
	PricingLumpsum  = "lumpsum"   // Total = UnitPrice
)

// SPK is a "Surat Perintah Kerja" Legal Permit (Proses J). The number is
// generated automatically (SPK/LP/0001/VI/2026); the vendor is referenced from
// the Vendor master (Proses I-1).
type SPK struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Number   string `gorm:"size:48;uniqueIndex;not null" json:"number"`
	Type     string `gorm:"size:32;not null" json:"type"`
	TypeName string `gorm:"size:120" json:"type_name"`

	ProjectID *uint    `gorm:"index" json:"project_id"`
	Project   *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	VendorID  uint     `gorm:"index;not null" json:"vendor_id"`
	Vendor    *Vendor  `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`

	PricingMode string  `gorm:"size:16;not null" json:"pricing_mode"`
	UnitPrice   int64   `json:"unit_price"` // Rp per meter / per unit / lumpsum
	Quantity    float64 `json:"quantity"`   // luas (m²) atau total unit
	Total       int64   `json:"total"`      // nilai SPK (Rp)

	CompletionTime string `gorm:"size:120" json:"completion_time"` // waktu penyelesaian
	PaymentTerms   string `gorm:"type:text" json:"payment_terms"`  // termin pembayaran
	ScopeNote      string `gorm:"type:text" json:"scope_note"`

	Status       SPKStatus  `gorm:"size:16;not null;default:'draft'" json:"status"`
	CreatedBy    uint       `json:"created_by"`
	ApprovedBy   *uint      `json:"approved_by"`
	ApprovedAt   *time.Time `json:"approved_at"`
	DecisionNote string     `gorm:"type:text" json:"decision_note"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SPKType is one of the eight SPK templates (Proses J-1..J-8).
type SPKType struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	PricingMode string `json:"pricing_mode"`
}

// SPKTypes enumerates the SPK catalog from the spec (Proses J).
var SPKTypes = []SPKType{
	{Code: "bpn", Name: "SPK Pengukuran BPN", PricingMode: PricingPerMeter},
	{Code: "perizinan", Name: "SPK ITR, KKPR, Siteplan, IMB/PBG", PricingMode: PricingPerUnit},
	{Code: "tgt", Name: "SPK Aspek TGT", PricingMode: PricingLumpsum},
	{Code: "limbah", Name: "SPK Rekom Limbah", PricingMode: PricingLumpsum},
	{Code: "ukl", Name: "SPK Rekom UKL-UPL, Peil Banjir, Andalalin", PricingMode: PricingLumpsum},
	{Code: "shm", Name: "SPK Pemecahan SHM", PricingMode: PricingPerUnit},
	{Code: "pbb", Name: "SPK Pemecahan PBB", PricingMode: PricingPerUnit},
	{Code: "lainnya", Name: "SPK Lainnya", PricingMode: PricingLumpsum},
}

// FindSPKType returns the catalog entry for a type code.
func FindSPKType(code string) (SPKType, bool) {
	for _, t := range SPKTypes {
		if t.Code == code {
			return t, true
		}
	}
	return SPKType{}, false
}
