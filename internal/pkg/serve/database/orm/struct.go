package orm

type ErrMsg struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
}

// Model
type Model struct {
	Id        uint   `json:"id" uri:"id" form:"id" param:"id"`
	UpdatedAt string `json:"updatedAt" uri:"updatedAt" form:"updatedAt" param:"updatedAt"`
	CreatedAt string `json:"createdAt" uri:"createdAt" form:"createdAt" param:"createdAt"`
	DeletedAt string `json:"deletedAt" uri:"deletedAt" form:"deletedAt" param:"deletedAt"`
}
