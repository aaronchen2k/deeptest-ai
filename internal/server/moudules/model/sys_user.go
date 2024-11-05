package model

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	"gorm.io/gorm"
)

type SysUser struct {
	BaseModel
	domain.BaseUser

	Password  string   `gorm:"type:varchar(250)" json:"password" validate:"required"`
	RoleNames []string `gorm:"-" json:"role_names"`
}

type Avatar struct {
	Avatar string `json:"avatar"`
}

// Create 添加
func (item *SysUser) Create(db *gorm.DB) (id uint, err error) {
	//if db == nil {
	//	return 0, gorm.ErrInvalidDB
	//}
	//err := db.Model(item).Create(item).Error
	//if err != nil {
	//	zap_server.ZAPLOG.Error(err.Error())
	//	return item.ID, err
	//}
	//return item.ID, nil

	return
}

// Update 更新
func (item *SysUser) Update(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) (err error) {
	//if db == nil {
	//	return gorm.ErrInvalidDB
	//}
	//err := db.Model(item).Scopes(scopes...).Updates(item).Error
	//if err != nil {
	//	zap_server.ZAPLOG.Error(err.Error())
	//	return err
	//}
	//return nil

	return
}

// Delete 删除
func (item *SysUser) Delete(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) (err error) {
	//if db == nil {
	//	return gorm.ErrInvalidDB
	//}
	//err := db.Model(item).Unscoped().Scopes(scopes...).Delete(item).Error
	//if err != nil {
	//	zap_server.ZAPLOG.Error(err.Error())
	//	return err
	//}
	//return nil

	return
}

func (SysUser) TableName() string {
	return "sys_users"
}
