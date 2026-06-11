package service

import (
	"encoding/json"
	"errors"

	"legalpermit/internal/model"
	"legalpermit/internal/repository"

	"gorm.io/datatypes"
)

type SettingService struct {
	settings *repository.SettingRepository
}

func NewSettingService(settings *repository.SettingRepository) *SettingService {
	return &SettingService{settings: settings}
}

// EnsureDefaults seeds DACI and notification settings on first run.
func (s *SettingService) EnsureDefaults() error {
	if _, err := s.settings.Get(model.SettingDACI); errors.Is(err, repository.ErrNotFound) {
		if err := s.SetDACI(model.DefaultDACI()); err != nil {
			return err
		}
	}
	if _, err := s.settings.Get(model.SettingNotification); errors.Is(err, repository.ErrNotFound) {
		if err := s.SetNotification(model.DefaultNotification()); err != nil {
			return err
		}
	}
	return nil
}

func (s *SettingService) GetDACI() (model.DACIConfig, error) {
	var cfg model.DACIConfig
	raw, err := s.settings.Get(model.SettingDACI)
	if errors.Is(err, repository.ErrNotFound) {
		return model.DefaultDACI(), nil
	}
	if err != nil {
		return cfg, err
	}
	return cfg, json.Unmarshal(raw, &cfg)
}

func (s *SettingService) SetDACI(cfg model.DACIConfig) error {
	raw, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return s.settings.Upsert(model.SettingDACI, datatypes.JSON(raw))
}

func (s *SettingService) GetNotification() (model.NotificationConfig, error) {
	var cfg model.NotificationConfig
	raw, err := s.settings.Get(model.SettingNotification)
	if errors.Is(err, repository.ErrNotFound) {
		return model.DefaultNotification(), nil
	}
	if err != nil {
		return cfg, err
	}
	return cfg, json.Unmarshal(raw, &cfg)
}

func (s *SettingService) SetNotification(cfg model.NotificationConfig) error {
	raw, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return s.settings.Upsert(model.SettingNotification, datatypes.JSON(raw))
}
