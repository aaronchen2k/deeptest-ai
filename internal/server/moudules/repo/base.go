package repo

import (
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/database"
	_str "github.com/deeptest-com/deeptest-next/pkg/libs/string"

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

func (r *BaseRepo) PaginateScope(page, pageSize int, sort, orderBy string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize < 0:
			pageSize = -1
			//case pageSize == 0:
			//	pageSize = 10
		}

		if sort == "" {
			sort = "desc"
		}
		if orderBy == "" {
			orderBy = "created_at"
		}

		offset := (page - 1) * pageSize
		if page < 0 {
			offset = -1
		}
		db = db.Order(_str.Join(orderBy, " ", sort)).Offset(offset)
		if pageSize > 0 {
			db = db.Limit(pageSize)
		}
		return db
	}
}

func (r *BaseRepo) GetAdminRoleName() (roleName consts.RoleType) {
	roleName = consts.Admin
	return
}

func (r *BaseRepo) IsDisable(enable string) bool {
	if enable == "1" || enable == "true" {
		return false
	} else {
		return true
	}
}

func (r *BaseRepo) ArrayRemoveUintDuplication(arr []uint) []uint {
	set := make(map[uint]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}

	return arr[:j]
}

func GetDbInstance() (db *gorm.DB) {
	return database.GetInstance()
}
