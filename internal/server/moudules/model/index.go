package model

var (
	Models = []interface{}{
		&SysPerm{},
		&SysRole{},
		&SysUser{},
		&SysUserProfile{},

		&ProjectRole{},
		&Org{},
		&Project{},
		&ProjectMember{},
	}
)
