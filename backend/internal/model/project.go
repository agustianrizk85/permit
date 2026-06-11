package model

import "time"

// ProjectStage tracks how far a land/project has progressed through the
// macro phases (A Pra-Akad -> B Akad -> C Permit -> D Legal).
type ProjectStage string

const (
	StagePraAkad ProjectStage = "pra_akad" // Proses A
	StageAkad    ProjectStage = "akad"     // Proses B
	StagePermit  ProjectStage = "permit"   // Proses C
	StageLegal   ProjectStage = "legal"    // Proses D
	StageDone    ProjectStage = "done"
)

// Project represents a parcel of land (lahan) being acquired and permitted.
type Project struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	Name       string       `gorm:"size:160;not null" json:"name"`
	Location   string       `gorm:"size:255" json:"location"`
	OwnerName  string       `gorm:"size:160" json:"owner_name"` // pemilik lahan
	PTName     string       `gorm:"size:160" json:"pt_name"`    // PT yang dipakai
	Stage      ProjectStage `gorm:"size:32;not null;default:pra_akad" json:"stage"`
	CreatedBy  uint         `json:"created_by"`
	Creator    *User        `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`

	Steps []ProcessStep `gorm:"foreignKey:ProjectID" json:"steps,omitempty"`
}
