package repository

import (
	"errors"

	"legalpermit/internal/model"

	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(d *model.Document) error {
	return r.db.Create(d).Error
}

func (r *DocumentRepository) FindByID(id uint) (*model.Document, error) {
	var d model.Document
	err := r.db.First(&d, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DocumentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Document{}, id).Error
}

// Search matches documents by type or original name (case-insensitive),
// optionally scoped to one project.
func (r *DocumentRepository) Search(query string, projectID *uint) ([]model.Document, error) {
	q := r.db.Model(&model.Document{}).Order("created_at desc")
	if projectID != nil {
		q = q.Where("project_id = ?", *projectID)
	}
	if query != "" {
		like := "%" + query + "%"
		q = q.Where("LOWER(doc_type) LIKE LOWER(?) OR LOWER(original_name) LIKE LOWER(?)", like, like)
	}
	var docs []model.Document
	err := q.Limit(200).Find(&docs).Error
	return docs, err
}
