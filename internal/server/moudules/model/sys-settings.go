package model

type Settings struct {
	Key   string `gorm:"column:k;index:key_index,unique;not null" json:"key"`
	Value string `gorm:"column:v;type:text" json:"value"`
}

func (Settings) TableName() string {
	return "sys_settings"
}
