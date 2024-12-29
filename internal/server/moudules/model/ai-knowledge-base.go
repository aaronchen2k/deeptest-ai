package model

type KbMaterial struct {
	BaseModel

	File string `json:"file"`
}

func (KbMaterial) TableName() string {
	return "ai_materials"
}

type KbDoc struct {
	BaseModel

	Name     string `json:"name"`
	SrcPath  string `json:"srcPath"`
	DictPath string `json:"dictPath"`

	MaterialId uint `json:"materialId"`
}

func (KbDoc) TableName() string {
	return "ai_kb_docs"
}
