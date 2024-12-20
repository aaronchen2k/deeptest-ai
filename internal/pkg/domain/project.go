package domain

type ProjectBase struct {
	Name string `json:"name"`
	Desc string `json:"desc" gorm:"type:text"`

	SchemaId       uint   `json:"schemaId"`
	OrgId          uint   `json:"orgId"`
	Logo           string `json:"logo"`
	ShortName      string `json:"shortName"`
	IncludeExample bool   `json:"includeExample"`
}
