package model

var (
	Models = []interface{}{
		&TestCase{},
		&TestSet{},
		&TestPlan{},

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
