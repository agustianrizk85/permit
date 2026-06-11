package dto

import (
	"time"

	"legalpermit/internal/model"
)

// --- Auth ---

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	ExpiresAt string      `json:"expires_at"`
	User      *model.User `json:"user"`
}

// --- Project ---

type CreateProjectRequest struct {
	Name      string `json:"name" binding:"required"`
	Location  string `json:"location"`
	OwnerName string `json:"owner_name"`
	PTName    string `json:"pt_name"`
}

// --- Step ---

// UpdateStepRequest carries partial updates. Pointer fields allow distinguishing
// "not provided" from a zero value (important for price = 0 / clearing notes).
type UpdateStepRequest struct {
	Status    *model.StepStatus `json:"status"`
	PriceFix  *int64            `json:"price_fix"`
	SPKNumber *string           `json:"spk_number"`
	Notes     *string           `json:"notes"`
	Metadata  map[string]any    `json:"metadata"`
	SLADays   *int              `json:"sla_days"`   // KADEP dynamic SLA
	DueDate   *time.Time        `json:"due_date"`   // explicit override
}

// --- Dashboard ---

type ProjectProgress struct {
	ProjectID  uint                       `json:"project_id"`
	Total      int64                      `json:"total"`
	Done       int64                      `json:"done"`
	Percentage int                        `json:"percentage"`
	ByStatus   map[model.StepStatus]int64 `json:"by_status"`
}
