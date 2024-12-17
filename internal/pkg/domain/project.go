package domain

type ProjectBase struct {
	Name string `json:"name"`
	Desc string `json:"desc" gorm:"column:descr;type:text"`

	SchemaId       uint   `json:"schemaId"`
	OrgId          uint   `json:"orgId"`
	Logo           string `json:"logo"`
	ShortName      string `json:"shortName"`
	IncludeExample bool   `json:"includeExample"`
	AdminId        uint   `json:"adminId"`
	AdminName      string `gorm:"-" json:"adminName"`
}
