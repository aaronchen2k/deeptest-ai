package repo

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"gorm.io/gorm"
)

type CaseRepo struct {
	DB          *gorm.DB `inject:""`
	*BaseRepo   `inject:""`
	ProjectRepo *ProjectRepo `inject:""`
	UserRepo    *UserRepo    `inject:""`
}

func (r *CaseRepo) LoadTree(projectId uint) (root *domain.CaseNode, err error) {
	pos, err := r.ListByProject(projectId)
	if err != nil {
		return
	}

	if len(pos) == 0 {
		return
	}

	tos := r.toTos(pos)

	root = &domain.CaseNode{}
	r.makeTree(tos, root)
	r.mountCount(root)

	return
}

func (r *CaseRepo) ListByProject(projectId uint) (pos []*model.TestCase, err error) {
	db := r.DB.
		Where("project_id=?", projectId).
		Where("NOT deleted")

	err = db.
		Order("parent_id ASC, ordr ASC").
		Find(&pos).Error

	return
}

func (r *CaseRepo) Get(id uint) (scenario model.TestCase, err error) {
	err = r.DB.Model(&model.TestCase{}).Where("id = ?", id).First(&scenario).Error
	if err != nil {
		return
	}

	return
}

func (r *CaseRepo) Create(po model.TestCase) (ret model.TestCase, err error) {
	if po.ID == 0 {
		po.Ordr = r.GetMaxOrder(po.ParentId)
	}

	err = r.DB.Model(&model.TestCase{}).Create(&po).Error
	if err != nil {
		return
	}

	ret = po

	return
}

func (r *CaseRepo) Update(req model.TestCase) (ret model.TestCase, err error) {
	values := map[string]interface{}{
		"title":  req.Title,
		"desc":   req.Desc,
		"status": req.Status,

		"disabled": req.Disabled,
		"deleted":  req.Deleted,
	}
	err = r.DB.Model(&req).Where("id = ?", req.ID).Updates(values).Error
	if err != nil {
		return
	}

	ret = req

	return
}

func (r *CaseRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.TestCase{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		return
	}

	return
}

func (r *CaseRepo) makeTree(findIn []*domain.CaseNode, parent *domain.CaseNode) { //参数为父节点，添加父节点的子节点指针切片
	children, _ := r.hasChild(findIn, parent) // 判断节点是否有子节点并返回

	if children != nil {
		parent.Children = append(parent.Children, children[0:]...) // 添加子节点

		for _, child := range children { // 查询子节点的子节点，并添加到子节点
			_, has := r.hasChild(findIn, child)
			if has {
				r.makeTree(findIn, child) // 递归添加节点
			}
		}
	}
}

func (r *CaseRepo) hasChild(categories []*domain.CaseNode, parent *domain.CaseNode) (
	ret []*domain.CaseNode, yes bool) {

	for _, item := range categories {
		if item.ParentId == parent.Id {
			//item.Parent = parent // loop json
			ret = append(ret, item)
		}
	}

	if ret != nil {
		yes = true
	}

	return
}

func (r *CaseRepo) toTos(pos []*model.TestCase) (tos []*domain.CaseNode) {
	for _, po := range pos {
		to := r.ToTo(po)

		tos = append(tos, to)
	}

	return
}
func (r *CaseRepo) ToTo(po *model.TestCase) (to *domain.CaseNode) {
	to = &domain.CaseNode{
		Id:       po.ID,
		Title:    po.Title,
		Desc:     po.Desc,
		Type:     po.Type,
		Ordr:     po.Ordr,
		ParentId: po.ParentId,
		IsDir:    true,
	}

	if po.Type == consts.NodeLeaf {
		to.IsDir = false
		to.Count = 1
	}

	return
}

func (r *CaseRepo) mountCount(node *domain.CaseNode) (count int) {
	for _, child := range node.Children {
		node.Count += r.mountCount(child)
	}
	return node.Count

}

func (r *CaseRepo) GetMaxOrder(parentId uint) (order int) {
	node := model.TestCase{}

	err := r.DB.Where("parent_id=?", parentId).
		Order("ordr DESC").
		First(&node).Error

	if err == nil {
		order = node.Ordr + 1
	}

	return
}

func (r *CaseRepo) UpdateOrder(pos consts.DropPos, targetId uint, projectId uint) (
	parentId uint, ordr int) {
	if pos == consts.Inner {
		parentId = targetId

		var preChild model.TestCase
		r.DB.Where("parent_id=? AND project_id = ?", parentId, projectId).
			Order("ordr DESC").Limit(1).
			First(&preChild)

		ordr = preChild.Ordr + 1

	} else if pos == consts.Before {
		brother, _ := r.Get(targetId)
		parentId = brother.ParentId

		r.DB.Model(&model.TestCase{}).
			Where("NOT deleted AND parent_id=? AND project_id = ? AND ordr >= ?",
				parentId, projectId, brother.Ordr).
			Update("ordr", gorm.Expr("ordr + 1"))

		ordr = brother.Ordr

	} else if pos == consts.After {
		brother, _ := r.Get(targetId)
		parentId = brother.ParentId

		r.DB.Model(&model.TestCase{}).
			Where("NOT deleted AND parent_id=? AND project_id = ? AND ordr > ?",
				parentId, parentId, brother.Ordr).
			Update("ordr", gorm.Expr("ordr + 1"))

		ordr = brother.Ordr + 1

	}

	return
}
func (r *CaseRepo) UpdateOrdAndParent(node model.TestCase) (err error) {
	err = r.DB.Model(&node).
		Updates(model.TestCase{Ordr: node.Ordr, ParentId: node.ParentId}).
		Error

	return
}
