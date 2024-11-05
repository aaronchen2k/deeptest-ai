package inits

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	_file "github.com/deeptest-com/deeptest-next/pkg/libs/file"
)

func Init() {
	consts.ExecDir = _file.GetExecDir()
	consts.WorkDir = _file.GetWorkDir()
}
