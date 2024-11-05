package repo

import (
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/database"

	"gorm.io/gorm"
)

type BaseRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *BaseRepo) GetAncestorIds(id uint, tableName string) (ids []uint, err error) {
	sql := `
		WITH RECURSIVE temp AS
		(
			SELECT id, parent_id, name from %s a where a.id = %d
		
			UNION ALL
		
			SELECT b.id, b.parent_id, b.name 
				from temp c
				inner join %s b on b.id = c.parent_id
		) 
		select id from temp e;
`

	sql = fmt.Sprintf(sql, tableName, id, tableName)

	err = r.DB.Raw(sql).Scan(&ids).Error
	if err != nil {
		return
	}

	return
}

func (r *BaseRepo) GetDescendantIds(id uint, tableName string, typ consts.CategoryDiscriminator, projectId int) (
	ids []uint, err error) {
	sql := `
		WITH RECURSIVE temp AS
		(
			SELECT id, parent_id from %s a 
				WHERE a.id = %d AND type='%s' AND project_id=%d AND NOT a.deleted
		
			UNION ALL
		
			SELECT b.id, b.parent_id 
				from temp c
				inner join %s b on b.parent_id = c.id
				WHERE type='%s' AND project_id=%d AND NOT b.deleted
		) 
		select id from temp e;
`
	sql = fmt.Sprintf(sql, tableName,
		id, typ, projectId,
		tableName,
		typ, projectId)

	err = r.DB.Raw(sql).Scan(&ids).Error
	if err != nil {
		return
	}

	return
}

func (r *BaseRepo) GetAllChildIdsSimple(id uint, tableName string) (
	ids []uint, err error) {
	sql := `
		WITH RECURSIVE temp AS
		(
			SELECT id, parent_id from %s a 
				WHERE a.id = %d AND NOT a.deleted
		
			UNION ALL
		
			SELECT b.id, b.parent_id 
				from temp c
				inner join %s b on b.parent_id = c.id
				WHERE NOT b.deleted
		) 
		select id from temp e;
`
	sql = fmt.Sprintf(sql, tableName, id, tableName)

	err = r.DB.Raw(sql).Scan(&ids).Error
	if err != nil {
		return
	}

	return
}

func (r *BaseRepo) GetAdminRoleName() (roleName consts.RoleType) {
	roleName = consts.Admin
	return
}

func GetDbInstance() (db *gorm.DB) {
	return database.GetInstance()
}
