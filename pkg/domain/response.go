package _domain

type Response struct {
	Code int64 `json:"code"`

	Msg    string `json:"msg"`
	MsgKey string `json:"msgKey"` // show i118 msg on client side

	Data interface{} `json:"data"`
}

type PageData struct {
	Items interface{} `json:"items"`

	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func (d *PageData) Populate(result interface{}, total int64, page, pageSize int) {
	d.Items = result
	d.Total = int(total)
	d.Page = page
	d.PageSize = pageSize
}
