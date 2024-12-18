package v1

import (
	"time"
)

type BaseDomain struct {
	ID uint `json:"id"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	Disabled bool `json:"disabled,omitempty"`
}

type DataReq struct {
	ClearData bool    `json:"clearData"`
	Sys       DataSys `json:"sys"`
}
type DataSys struct {
	AdminPassword string `json:"adminPassword"`
}
