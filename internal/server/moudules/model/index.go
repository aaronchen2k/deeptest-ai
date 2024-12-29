package model

var (
	Models = []interface{}{
		&TestCase{},
		&TestSet{},
		&TestPlan{},

		KbMaterial{},
		KbDoc{},

		&ProjectRole{},
		&Org{},
		&Project{},
		&ProjectMember{},

		&SysPerm{},
		&SysRole{},
		&SysUser{},
		&SysUserProfile{},
	}
)
