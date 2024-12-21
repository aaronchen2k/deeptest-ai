package v1

type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResp struct {
	BaseDomain

	Password     string   `json:"password"`
	AuthorityIds []string `gorm:"-" json:"authorityIds"`
}

type UpdateUserProjectReq struct {
	ProjectId uint `json:"projectId" validate:"required"`
}
