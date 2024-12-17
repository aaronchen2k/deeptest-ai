package model

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
)

type SysOplog struct {
	BaseModel
	domain.BaseOplog
}

func (SysOplog) TableName() string {
	return "sys_oplogs"
}
