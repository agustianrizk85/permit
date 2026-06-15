package repository

import (
	"errors"

	"legalpermit/internal/model"

	"gorm.io/gorm"
)

type VendorRepository struct {
	db *gorm.DB
}

func NewVendorRepository(db *gorm.DB) *VendorRepository {
	return &VendorRepository{db: db}
}

func (r *VendorRepository) Create(v *model.Vendor) error {
	return r.db.Create(v).Error
}

func (r *VendorRepository) Update(v *model.Vendor) error {
	return r.db.Save(v).Error
}

func (r *VendorRepository) List() ([]model.Vendor, error) {
	var vendors []model.Vendor
	err := r.db.Order("name asc").Find(&vendors).Error
	return vendors, err
}

func (r *VendorRepository) FindByID(id uint) (*model.Vendor, error) {
	var v model.Vendor
	err := r.db.First(&v, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}
