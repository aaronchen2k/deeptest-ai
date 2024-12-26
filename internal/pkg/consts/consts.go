package consts

import "errors"

const (
	AdminUserName     = "admin"
	AdminUserPassword = "P2ssw0rd"
	AdminRoleName     = "admin"

	ConfigType     = "json"
	CasbinFileName = "rbac_model.conf"

	DirUpload = "upload"
	DirUi     = "deeptest-ui"
)

var (
	System = "deeptest"
	App    = "server"

	ExecDir = ""
	WorkDir = ""

	ConfDir = "config"
)

var (
	ErrUserNameOrPassword = errors.New("用户名或密码错误")
	ErrUserNameInvalid    = errors.New("用户名名称已经被使用")
	ErrRoleNameInvalid    = errors.New("角色名称已经被使用")

	ErrParamValidate      = errors.New("参数验证失败")
	ErrPaginateParam      = errors.New("分页查询参数缺失")
	ErrUnSupportFramework = errors.New("不支持的框架")
)
