package service

import (
	"errors"
	"fmt"
	"time"

	"legalpermit/internal/model"
	"legalpermit/internal/repository"
)

var (
	// ErrInvalidSPKType is returned when an unknown SPK type code is supplied.
	ErrInvalidSPKType = errors.New("jenis SPK tidak dikenal")
	// ErrSPKNotDraft is returned when approving/rejecting a non-draft SPK.
	ErrSPKNotDraft = errors.New("SPK sudah diputuskan (bukan draft)")
)

type SPKService struct {
	spks    *repository.SPKRepository
	vendors *repository.VendorRepository
}

func NewSPKService(spks *repository.SPKRepository, vendors *repository.VendorRepository) *SPKService {
	return &SPKService{spks: spks, vendors: vendors}
}

type CreateSPKInput struct {
	Type           string  `json:"type" binding:"required"`
	VendorID       uint    `json:"vendor_id" binding:"required"`
	ProjectID      *uint   `json:"project_id"`
	UnitPrice      int64   `json:"unit_price"`
	Quantity       float64 `json:"quantity"`
	CompletionTime string  `json:"completion_time"`
	PaymentTerms   string  `json:"payment_terms"`
	ScopeNote      string  `json:"scope_note"`
}

// Create issues a draft SPK with an auto-generated number. Kadep only (enforced
// at the route). The total is derived from the type's pricing mode.
func (s *SPKService) Create(in CreateSPKInput, createdBy uint) (*model.SPK, error) {
	spkType, ok := model.FindSPKType(in.Type)
	if !ok {
		return nil, ErrInvalidSPKType
	}
	if _, err := s.vendors.FindByID(in.VendorID); err != nil {
		return nil, err
	}

	number, err := s.nextNumber()
	if err != nil {
		return nil, err
	}

	spk := &model.SPK{
		Number:         number,
		Type:           spkType.Code,
		TypeName:       spkType.Name,
		ProjectID:      in.ProjectID,
		VendorID:       in.VendorID,
		PricingMode:    spkType.PricingMode,
		UnitPrice:      in.UnitPrice,
		Quantity:       in.Quantity,
		Total:          computeTotal(spkType.PricingMode, in.UnitPrice, in.Quantity),
		CompletionTime: in.CompletionTime,
		PaymentTerms:   in.PaymentTerms,
		ScopeNote:      in.ScopeNote,
		Status:         model.SPKDraft,
		CreatedBy:      createdBy,
	}
	if err := s.spks.Create(spk); err != nil {
		return nil, err
	}
	return s.spks.FindByID(spk.ID)
}

func (s *SPKService) List(status string) ([]model.SPK, error) { return s.spks.List(status) }

func (s *SPKService) Get(id uint) (*model.SPK, error) { return s.spks.FindByID(id) }

// Decide approves or rejects a draft SPK. Dirops only (enforced at the route).
func (s *SPKService) Decide(id uint, approve bool, approverID uint, note string) (*model.SPK, error) {
	spk, err := s.spks.FindByID(id)
	if err != nil {
		return nil, err
	}
	if spk.Status != model.SPKDraft {
		return nil, ErrSPKNotDraft
	}
	now := time.Now().UTC()
	if approve {
		spk.Status = model.SPKApproved
	} else {
		spk.Status = model.SPKRejected
	}
	spk.ApprovedBy = &approverID
	spk.ApprovedAt = &now
	spk.DecisionNote = note
	if err := s.spks.Update(spk); err != nil {
		return nil, err
	}
	return s.spks.FindByID(spk.ID)
}

// nextNumber builds the running SPK number, e.g. SPK/LP/0001/VI/2026.
func (s *SPKService) nextNumber() (string, error) {
	now := time.Now()
	year := now.Year()
	count, err := s.spks.CountByYear(year)
	if err != nil {
		return "", err
	}
	seq := count + 1
	return fmt.Sprintf("SPK/LP/%04d/%s/%d", seq, romanMonth(now.Month()), year), nil
}

func computeTotal(mode string, unitPrice int64, qty float64) int64 {
	switch mode {
	case model.PricingPerMeter, model.PricingPerUnit:
		return int64(float64(unitPrice) * qty)
	default: // lumpsum
		return unitPrice
	}
}

var romanMonths = [...]string{
	"I", "II", "III", "IV", "V", "VI",
	"VII", "VIII", "IX", "X", "XI", "XII",
}

func romanMonth(m time.Month) string {
	if m < 1 || m > 12 {
		return ""
	}
	return romanMonths[m-1]
}
