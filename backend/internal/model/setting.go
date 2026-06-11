package model

import (
	"time"

	"gorm.io/datatypes"
)

// Setting is a single JSON-valued configuration row (DACI roster, notification
// preferences, …). Editable by KADEP per the spec's "DINAMIS SETTING KADEP".
type Setting struct {
	Key       string         `gorm:"primaryKey;size:64" json:"key"`
	Value     datatypes.JSON `gorm:"type:jsonb" json:"value"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// Setting keys.
const (
	SettingDACI         = "daci"
	SettingNotification = "notification"
)

// DACIConfig captures the DACI roster from the spec (Driver / Approver /
// Consulting / Informed) plus named drivers.
type DACIConfig struct {
	Drivers    []DACIDriver `json:"drivers"`
	Approver   []string     `json:"approver"`
	Consulting []string     `json:"consulting"`
	Informed   []string     `json:"informed"`
}

type DACIDriver struct {
	Code string `json:"code"` // e.g. "GP124"
	Name string `json:"name"` // e.g. "Kusnadi"
}

// NotificationConfig controls WA discipline reminders & audit chatbot toggles.
type NotificationConfig struct {
	WhatsAppEnabled bool   `json:"whatsapp_enabled"`
	AuditAIEnabled  bool   `json:"audit_ai_enabled"`
	ReminderHour    int    `json:"reminder_hour"`  // local hour to send daily discipline reminder
	WhatsAppAPIURL  string `json:"whatsapp_api_url"`
}

// DefaultDACI is seeded on first run from the spec.
func DefaultDACI() DACIConfig {
	return DACIConfig{
		Drivers: []DACIDriver{
			{Code: "GP124", Name: "Kusnadi"},
			{Code: "GP3", Name: "Fadil"},
		},
		Approver:   []string{"kadep", "dirops"},
		Consulting: []string{"kadep", "dirops", "ceo"},
		Informed:   []string{"dirops", "ceo"},
	}
}

// DefaultNotification is seeded on first run.
func DefaultNotification() NotificationConfig {
	return NotificationConfig{WhatsAppEnabled: false, AuditAIEnabled: false, ReminderHour: 9}
}
