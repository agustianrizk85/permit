package service

import (
	"legalpermit/internal/model"
	"legalpermit/internal/repository"
)

type DeadlineService struct {
	deadlines *repository.DeadlineRepository
}

func NewDeadlineService(deadlines *repository.DeadlineRepository) *DeadlineService {
	return &DeadlineService{deadlines: deadlines}
}

// EnsureDefaults seeds one Master Deadline rule per catalog step on first run.
// The default mirrors the previous behaviour (deadline + alert on for steps that
// had an SLA), but every rule is now editable by KADEP.
func (s *DeadlineService) EnsureDefaults() error {
	n, err := s.deadlines.Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	templates := Catalog()
	rules := make([]model.DeadlineRule, 0, len(templates))
	for _, t := range templates {
		on := t.SLADays > 0
		rules = append(rules, model.DeadlineRule{
			Code:            t.Code,
			Name:            t.Name,
			Category:        t.Category,
			DeadlineEnabled: on,
			AlertEnabled:    on,
			SLADays:         t.SLADays,
		})
	}
	return s.deadlines.Upsert(rules)
}

func (s *DeadlineService) List() ([]model.DeadlineRule, error) {
	return s.deadlines.List()
}

func (s *DeadlineService) Map() (map[string]model.DeadlineRule, error) {
	return s.deadlines.Map()
}

// Update applies edits from the Master Deadline editor. Negative SLA is clamped
// to zero; a disabled deadline implies no alert.
func (s *DeadlineService) Update(rules []model.DeadlineRule) ([]model.DeadlineRule, error) {
	for i := range rules {
		if rules[i].SLADays < 0 {
			rules[i].SLADays = 0
		}
		if !rules[i].DeadlineEnabled {
			rules[i].AlertEnabled = false
		}
	}
	if err := s.deadlines.Upsert(rules); err != nil {
		return nil, err
	}
	return s.deadlines.List()
}
