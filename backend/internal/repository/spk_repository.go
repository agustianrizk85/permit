package repository

import (
	"errors"
	"strconv"

	"legalpermit/internal/model"

	"gorm.io/gorm"
)

type SPKRepository struct {
	db *gorm.DB
}

func NewSPKRepository(db *gorm.DB) *SPKRepository {
	return &SPKRepository{db: db}
}

func (r *SPKRepository) Create(s *model.SPK) error {
	return r.db.Create(s).Error
}

func (r *SPKRepository) Update(s *model.SPK) error {
	return r.db.Save(s).Error
}

// List returns SPKs newest-first, optionally filtered by status, with vendor and
// project preloaded for display.
func (r *SPKRepository) List(status string) ([]model.SPK, error) {
	var spks []model.SPK
	q := r.db.Preload("Vendor").Preload("Project").Order("created_at desc")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	err := q.Find(&spks).Error
	return spks, err
}

func (r *SPKRepository) FindByID(id uint) (*model.SPK, error) {
	var s model.SPK
	err := r.db.Preload("Vendor").Preload("Project").First(&s, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// CountByYear counts SPKs whose number was issued in the given year, used to
// derive the next running sequence (numbers contain "/<year>").
func (r *SPKRepository) CountByYear(year int) (int64, error) {
	var n int64
	// Numbers look like SPK/LP/0001/VI/2026 — match the trailing /<year>.
	err := r.db.Model(&model.SPK{}).Where("number LIKE ?", "%/"+strconv.Itoa(year)).Count(&n).Error
	return n, err
}
