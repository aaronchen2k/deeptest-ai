package repo

import (
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"gorm.io/gorm"
)

type KnowledgeBaseRepo struct {
	DB        *gorm.DB `inject:""`
	*BaseRepo `inject:""`
}

func (r *KnowledgeBaseRepo) ListMaterial() (ret []model.KbMaterial, err error) {
	err = r.DB.Where("NOT deleted").
		Find(&ret).Error
	if err != nil {
		return
	}

	return
}

func (r *KnowledgeBaseRepo) ListDoc(materialId uint) (ret []model.KbDoc, err error) {
	err = r.DB.Where("NOT deleted AND where material_id = ?", materialId).
		Find(&ret).Error
	if err != nil {
		return
	}

	return
}

func (r *KnowledgeBaseRepo) SaveMaterial(po *model.KbMaterial) (err error) {
	err = r.DB.Save(po).Error

	return
}

func (r *KnowledgeBaseRepo) SaveDoc(po *model.KbDoc) (err error) {
	err = r.DB.Save(po).Error

	return
}

func (r *KnowledgeBaseRepo) DeleteMaterial(id uint) (err error) {
	err = r.DB.Model(&model.KbMaterial{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error

	return
}

func (r *KnowledgeBaseRepo) DeleteDoc(id uint) (err error) {
	err = r.DB.Model(&model.KbDoc{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error

	return
}
