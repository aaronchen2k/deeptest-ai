package repo

import (
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"gorm.io/gorm"
)

type SettingsRepo struct {
	*BaseRepo `inject:""`
	DB        *gorm.DB `inject:""`
}

func (r *SettingsRepo) Get(key string) (config model.Settings, err error) {
	err = r.DB.
		Where("k = ?", key).
		First(&config).Error
	return
}

func (r *SettingsRepo) Save(req model.Settings) (err error) {
	config, err := r.Get(req.Key)
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	if config.Key == "" || err == gorm.ErrRecordNotFound {
		if err = r.DB.Model(&req).Create(req).Error; err != nil {
			return err
		}
	}

	err = r.DB.Model(&model.Settings{}).
		Where("k = ?", req.Key).
		Update("v", req.Value).Error

	return
}
