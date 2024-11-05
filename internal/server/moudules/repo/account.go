package repo

import (
	"gorm.io/gorm"
)

type AccountRepo struct {
	*BaseRepo `inject:""`
	DB        *gorm.DB `inject:""`
}
