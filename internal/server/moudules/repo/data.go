package repo

import (
	"gorm.io/gorm"
)

type DataRepo struct {
	*BaseRepo `inject:""`
	DB        *gorm.DB `inject:""`
}
