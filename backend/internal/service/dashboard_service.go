package service

import (
	"fmt"
	"time"

	"legalpermit/internal/repository"
)

// WarningSeverity ranks an early-warning item.
type WarningSeverity string

const (
	SeverityCritical WarningSeverity = "critical" // overdue
	SeverityWarning  WarningSeverity = "warning"  // due soon
	SeverityInfo     WarningSeverity = "info"     // missing required input
)

// Warning is one early-warning entry for the dashboard.
type Warning struct {
	ProjectID   uint            `json:"project_id"`
	ProjectName string          `json:"project_name"`
	StepID      uint            `json:"step_id"`
	StepCode    string          `json:"step_code"`
	StepName    string          `json:"step_name"`
	Severity    WarningSeverity `json:"severity"`
	Message     string          `json:"message"`
	DueDate     *time.Time      `json:"due_date"`
}

// DashboardService produces the AI-style early-warning feed. The rules below are
// deterministic; an LLM provider can later enrich Message via the OCR/AI hook.
type DashboardService struct {
	steps     *repository.StepRepository
	deadlines *DeadlineService
}

func NewDashboardService(steps *repository.StepRepository, deadlines *DeadlineService) *DashboardService {
	return &DashboardService{steps: steps, deadlines: deadlines}
}

const dueSoonWindow = 3 * 24 * time.Hour

// EarlyWarnings scans all open steps and flags SLA breaches and missing inputs.
func (s *DashboardService) EarlyWarnings() ([]Warning, error) {
	open, err := s.steps.OpenSteps()
	if err != nil {
		return nil, err
	}
	rules, err := s.deadlines.Map()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	warnings := make([]Warning, 0)

	for _, st := range open {
		base := Warning{
			ProjectID:   st.ProjectID,
			ProjectName: st.ProjectName,
			StepID:      st.ID,
			StepCode:    st.Code,
			StepName:    st.Name,
			DueDate:     st.DueDate,
		}

		// Deadline alerts only fire when the Master Deadline rule enables alerts
		// for this step. Missing-input warnings below are independent.
		alertOn := true
		if rule, ok := rules[st.Code]; ok {
			alertOn = rule.AlertEnabled
		}

		// SLA rules.
		if alertOn && st.DueDate != nil {
			switch {
			case now.After(*st.DueDate):
				days := int(now.Sub(*st.DueDate).Hours() / 24)
				w := base
				w.Severity = SeverityCritical
				w.Message = fmt.Sprintf("Terlambat %d hari dari deadline SLA.", days)
				warnings = append(warnings, w)
			case st.DueDate.Sub(now) <= dueSoonWindow:
				days := int(st.DueDate.Sub(now).Hours()/24) + 1
				w := base
				w.Severity = SeverityWarning
				w.Message = fmt.Sprintf("Mendekati deadline (sisa ~%d hari).", days)
				warnings = append(warnings, w)
			}
		}

		// Missing-input rules (block completion later).
		if st.RequiresPrice && st.PriceFix <= 0 {
			w := base
			w.Severity = SeverityInfo
			w.Message = fmt.Sprintf("%s belum diisi.", priceLabel(st.PriceLabel))
			warnings = append(warnings, w)
		}
		if st.RequiresSPK && st.SPKNumber == "" {
			w := base
			w.Severity = SeverityInfo
			w.Message = "Nomor SPK belum diisi."
			warnings = append(warnings, w)
		}
	}
	return warnings, nil
}

func priceLabel(label string) string {
	if label == "" {
		return "Harga Fix"
	}
	return label
}
