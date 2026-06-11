package repository

import (
	"errors"

	"legalpermit/internal/model"

	"gorm.io/gorm"
)

type StepRepository struct {
	db *gorm.DB
}

func NewStepRepository(db *gorm.DB) *StepRepository {
	return &StepRepository{db: db}
}

func (r *StepRepository) FindByID(id uint) (*model.ProcessStep, error) {
	var s model.ProcessStep
	err := r.db.Preload("Documents").First(&s, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StepRepository) Save(s *model.ProcessStep) error {
	return r.db.Save(s).Error
}

// OpenStep is a not-yet-done step joined with its project, for early warnings.
type OpenStep struct {
	model.ProcessStep
	ProjectName string `json:"project_name"`
}

// OpenSteps returns every step that is not "done", with the project name, so the
// early-warning engine can evaluate SLA breaches and missing inputs.
func (r *StepRepository) OpenSteps() ([]OpenStep, error) {
	var rows []OpenStep
	err := r.db.
		Table("process_steps").
		Select("process_steps.*, projects.name as project_name").
		Joins("JOIN projects ON projects.id = process_steps.project_id").
		Where("process_steps.status <> ?", model.StatusDone).
		Order("process_steps.due_date asc").
		Scan(&rows).Error
	return rows, err
}

// CountByStatus returns how many steps a project has in each status, used by
// the dashboard progress summary.
func (r *StepRepository) CountByStatus(projectID uint) (map[model.StepStatus]int64, error) {
	type row struct {
		Status model.StepStatus
		Count  int64
	}
	var rows []row
	err := r.db.Model(&model.ProcessStep{}).
		Select("status, count(*) as count").
		Where("project_id = ?", projectID).
		Group("status").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := map[model.StepStatus]int64{}
	for _, r := range rows {
		out[r.Status] = r.Count
	}
	return out, nil
}
