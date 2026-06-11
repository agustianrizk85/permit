package service

import (
	"time"

	"legalpermit/internal/dto"
	"legalpermit/internal/model"
	"legalpermit/internal/repository"
)

type ProjectService struct {
	projects  *repository.ProjectRepository
	steps     *repository.StepRepository
	deadlines *DeadlineService
}

func NewProjectService(projects *repository.ProjectRepository, steps *repository.StepRepository, deadlines *DeadlineService) *ProjectService {
	return &ProjectService{projects: projects, steps: steps, deadlines: deadlines}
}

// Create persists a new project and seeds its checklist from the catalog.
func (s *ProjectService) Create(req dto.CreateProjectRequest, createdBy uint) (*model.Project, error) {
	project := &model.Project{
		Name:      req.Name,
		Location:  req.Location,
		OwnerName: req.OwnerName,
		PTName:    req.PTName,
		Stage:     model.StagePraAkad,
		CreatedBy: createdBy,
	}

	base := time.Now().UTC()
	rules, err := s.deadlines.Map()
	if err != nil {
		return nil, err
	}

	templates := Catalog()
	steps := make([]model.ProcessStep, 0, len(templates))
	for i, t := range templates {
		step := t.toModel(0, i+1)
		// Apply the Master Deadline configuration: a due date is only set when
		// the rule for this step is enabled.
		if rule, ok := rules[t.Code]; ok {
			step.SLADays = rule.SLADays
			if rule.DeadlineEnabled && rule.SLADays > 0 {
				due := base.AddDate(0, 0, rule.SLADays)
				step.DueDate = &due
			}
		}
		steps = append(steps, step)
	}

	if err := s.projects.CreateWithSteps(project, steps); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) List() ([]model.Project, error) {
	return s.projects.List()
}

func (s *ProjectService) Get(id uint) (*model.Project, error) {
	return s.projects.FindByID(id)
}

// Progress computes the checklist completion summary for the dashboard.
func (s *ProjectService) Progress(id uint) (*dto.ProjectProgress, error) {
	counts, err := s.steps.CountByStatus(id)
	if err != nil {
		return nil, err
	}
	var total, done int64
	for status, n := range counts {
		total += n
		if status == model.StatusDone {
			done += n
		}
	}
	pct := 0
	if total > 0 {
		pct = int(done * 100 / total)
	}
	return &dto.ProjectProgress{
		ProjectID:  id,
		Total:      total,
		Done:       done,
		Percentage: pct,
		ByStatus:   counts,
	}, nil
}
