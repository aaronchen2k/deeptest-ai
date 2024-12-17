package service

import (
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	"gorm.io/gorm"
)

type SettingsService struct {
	SettingsRepo *repo.SettingsRepo `inject:""`
}

func (s *SettingsService) Get(key string) (value string, err error) {
	config, err := s.SettingsRepo.Get(key)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	value = config.Value
	return
}

func (s *SettingsService) Save(req model.Settings) (err error) {
	err = s.SettingsRepo.Save(req)
	return
}
