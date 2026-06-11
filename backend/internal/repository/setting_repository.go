package repository

import (
	"errors"

	"legalpermit/internal/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

// Get returns the raw JSON for a key, or ErrNotFound.
func (r *SettingRepository) Get(key string) (datatypes.JSON, error) {
	var s model.Setting
	err := r.db.First(&s, "key = ?", key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return s.Value, nil
}

// Upsert inserts or updates a setting value.
func (r *SettingRepository) Upsert(key string, value datatypes.JSON) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&model.Setting{Key: key, Value: value}).Error
}
