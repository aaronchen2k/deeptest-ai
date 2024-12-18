package modules

// InitDBFunc 数据化初始化接口
type InitDBFunc interface {
	Init() (err error)
}
