package repository

import (
	"legalpermit/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeadlineRepository struct {
	db *gorm.DB
}

func NewDeadlineRepository(db *gorm.DB) *DeadlineRepository {
	return &DeadlineRepository{db: db}
}

func (r *DeadlineRepository) Count() (int64, error) {
	var n int64
	err := r.db.Model(&model.DeadlineRule{}).Count(&n).Error
	return n, err
}

func (r *DeadlineRepository) List() ([]model.DeadlineRule, error) {
	var rules []model.DeadlineRule
	err := r.db.Order("code asc").Find(&rules).Error
	return rules, err
}

// Map returns the rules keyed by step code for fast lookup during seeding and
// early-warning evaluation.
func (r *DeadlineRepository) Map() (map[string]model.DeadlineRule, error) {
	rules, err := r.List()
	if err != nil {
		return nil, err
	}
	out := make(map[string]model.DeadlineRule, len(rules))
	for _, rule := range rules {
		out[rule.Code] = rule
	}
	return out, nil
}

// Upsert inserts or updates a batch of rules (used by both seeding and the
// settings editor).
func (r *DeadlineRepository) Upsert(rules []model.DeadlineRule) error {
	if len(rules) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "category", "deadline_enabled", "alert_enabled", "sla_days", "updated_at"}),
	}).Create(&rules).Error
}
