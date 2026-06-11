package service

import (
	"encoding/json"
	"errors"
	"time"

	"legalpermit/internal/dto"
	"legalpermit/internal/model"
	"legalpermit/internal/repository"

	"gorm.io/datatypes"
)

var (
	ErrPriceRequired = errors.New("harga fix wajib diisi sebelum step ini diselesaikan")
	ErrSPKRequired   = errors.New("nomor SPK wajib diisi sebelum step ini diselesaikan")
)

type StepService struct {
	steps *repository.StepRepository
}

func NewStepService(steps *repository.StepRepository) *StepService {
	return &StepService{steps: steps}
}

func (s *StepService) Get(id uint) (*model.ProcessStep, error) {
	return s.steps.FindByID(id)
}

// Update applies a partial update and enforces the spec's completion rules.
func (s *StepService) Update(id uint, req dto.UpdateStepRequest, actor uint) (*model.ProcessStep, error) {
	step, err := s.steps.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.PriceFix != nil {
		step.PriceFix = *req.PriceFix
	}
	if req.SPKNumber != nil {
		step.SPKNumber = *req.SPKNumber
	}
	if req.Notes != nil {
		step.Notes = *req.Notes
	}
	if req.Metadata != nil {
		raw, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, err
		}
		step.Metadata = datatypes.JSON(raw)
	}
	if req.SLADays != nil {
		step.SLADays = *req.SLADays
	}
	if req.DueDate != nil {
		step.DueDate = req.DueDate
	}

	if req.Status != nil {
		if err := s.applyStatus(step, *req.Status, actor); err != nil {
			return nil, err
		}
	}

	if err := s.steps.Save(step); err != nil {
		return nil, err
	}
	return step, nil
}

// applyStatus validates the completion guardrails before changing status.
func (s *StepService) applyStatus(step *model.ProcessStep, status model.StepStatus, actor uint) error {
	if status == model.StatusDone {
		if step.RequiresPrice && step.PriceFix <= 0 {
			return ErrPriceRequired
		}
		if step.RequiresSPK && step.SPKNumber == "" {
			return ErrSPKRequired
		}
		now := time.Now().UTC()
		step.CompletedAt = &now
		step.CompletedBy = &actor
	} else {
		step.CompletedAt = nil
		step.CompletedBy = nil
	}
	step.Status = status
	return nil
}
