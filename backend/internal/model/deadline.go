package model

import "time"

// DeadlineRule is the per-step "Master Deadline" configuration. It lets KADEP
// decide, for each process step, whether the step uses a deadline at all and
// whether it raises alerts — instead of every step being forced to have one.
type DeadlineRule struct {
	Code            string    `gorm:"primaryKey;size:8" json:"code"`
	Name            string    `gorm:"size:200" json:"name"`
	Category        string    `gorm:"size:2;index" json:"category"`
	DeadlineEnabled bool      `json:"deadline_enabled"` // assign a due date on project creation
	AlertEnabled    bool      `json:"alert_enabled"`    // raise early-warning alerts
	SLADays         int       `json:"sla_days"`
	UpdatedAt       time.Time `json:"updated_at"`
}
