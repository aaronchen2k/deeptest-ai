package _domain

type BizErr struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
}

var (
	Success       = BizErr{0, "请求成功"}
	AuthErr       = BizErr{401, "请重新登录"}
	AuthActionErr = BizErr{403, "权限不足"}

	NeedInitErr = BizErr{1000, "未初始化"}
	ParamErr    = BizErr{2000, "参数解析失败"}

	RequestErr = BizErr{3000, "请求失败"}
	FailErr    = BizErr{3500, "处理失败"}
	SystemErr  = BizErr{4000, "系统错误"}
	LoginErr   = BizErr{5000, "登录失败"}

	ErrNameExist = BizErr{10100, "同名记录已存在"}

	ErrNoUser             = BizErr{10100, "找不到用户"}
	ErrUsernameExist      = BizErr{10200, "用户名已占用"}
	ErrEmailExist         = BizErr{10300, "邮箱已存在"}
	ErrShortNameExist     = BizErr{10400, "英文缩写已存在"}
	ErrPasswordMustBeSame = BizErr{10500, "两次密码必须一样"}

	ErrProjectNotExist  = BizErr{10700, "项目不存在"}
	ErrUserNotInProject = BizErr{10600, "不是该项目的成员"}
)

func (e BizErr) Error() string {
	return e.Msg
}
