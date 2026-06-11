package model

import "time"

// Role enumerates the system roles. Milestone 1 seeds only DIROPS and KADEP,
// but CEO and LEGAL_PERMIT are defined for the wider DACI flow (approver,
// consulting, informed).
type Role string

const (
	RoleCEO         Role = "ceo"
	RoleDirops      Role = "dirops"
	RoleKadep       Role = "kadep"
	RoleLegalPermit Role = "legal_permit"
)

// User is an application account.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:120;not null" json:"name"`
	Email        string    `gorm:"size:160;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         Role      `gorm:"size:32;not null" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
