package domain

type BasePermission struct {
	Name        string `gorm:"index:perm_index,unique;not null ;type:varchar(256)" json:"name" validate:"required,gte=4,lte=50"`
	Act         string `gorm:"index:perm_index;type:varchar(256)" json:"act" validate:"required"`
	DisplayName string `gorm:"type:varchar(256)" json:"displayName"`
	Description string `gorm:"type:varchar(256)" json:"description"`
}
