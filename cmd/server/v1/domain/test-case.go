package v1

import "github.com/deeptest-com/deeptest-next/internal/pkg/consts"

type CaseMoveReq struct {
	DragKey int            `json:"dragKey"`
	DropKey int            `json:"dropKey"`
	DropPos consts.DropPos `json:"dropPos"`
}
