package repository

import (
	"errors"

	"legalpermit/internal/model"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// CreateWithSteps persists the project and its seeded steps in one transaction.
func (r *ProjectRepository) CreateWithSteps(p *model.Project, steps []model.ProcessStep) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		for i := range steps {
			steps[i].ProjectID = p.ID
		}
		if len(steps) > 0 {
			if err := tx.Create(&steps).Error; err != nil {
				return err
			}
		}
		p.Steps = steps
		return nil
	})
}

func (r *ProjectRepository) List() ([]model.Project, error) {
	var projects []model.Project
	err := r.db.Order("created_at desc").Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) FindByID(id uint) (*model.Project, error) {
	var p model.Project
	err := r.db.
		Preload("Steps", func(db *gorm.DB) *gorm.DB { return db.Order("process_steps.sequence asc") }).
		Preload("Steps.Documents").
		First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepository) UpdateStage(id uint, stage model.ProjectStage) error {
	return r.db.Model(&model.Project{}).Where("id = ?", id).Update("stage", stage).Error
}
