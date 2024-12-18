package model

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	"gorm.io/gorm"
)

type PermCollection []SysPerm

// SysPerm 权鉴模块
type SysPerm struct {
	BaseModel
	domain.BasePermission
}

// Create 添加
func (item *SysPerm) Create(db *gorm.DB) (id uint, err error) {
	//if db == nil {
	//	return 0, gorm.ErrInvalidDB
	//}
	//if !service.CheckNameAndAct(NameScope(item.Name), ActScope(item.Act)) {
	//	return item.ID, errors.New(str.Join("权限[", item.Name, "-", item.Act, "]已存在"))
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
func (item *SysPerm) Update(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) (err error) {
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
func (item *SysPerm) Delete(db *gorm.DB, scopes ...func(db *gorm.DB) *gorm.DB) (err error) {
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

func (SysPerm) TableName() string {
	return "sys_perms"
}
