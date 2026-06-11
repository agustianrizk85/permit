package repository

import (
	"errors"

	"legalpermit/internal/model"

	"gorm.io/gorm"
)

type PTRepository struct {
	db *gorm.DB
}

func NewPTRepository(db *gorm.DB) *PTRepository {
	return &PTRepository{db: db}
}

func (r *PTRepository) Create(pt *model.PTMaster) error {
	return r.db.Create(pt).Error
}

func (r *PTRepository) List() ([]model.PTMaster, error) {
	var pts []model.PTMaster
	err := r.db.Order("created_at desc").Find(&pts).Error
	return pts, err
}

func (r *PTRepository) FindByID(id uint) (*model.PTMaster, error) {
	var pt model.PTMaster
	err := r.db.Preload("Documents").First(&pt, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &pt, nil
}

func (r *PTRepository) AddDocument(doc *model.PTDocument) error {
	return r.db.Create(doc).Error
}

func (r *PTRepository) FindDocument(id uint) (*model.PTDocument, error) {
	var d model.PTDocument
	err := r.db.First(&d, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}
