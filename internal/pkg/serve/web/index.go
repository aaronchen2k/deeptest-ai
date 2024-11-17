package web

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/viper_server"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
)

// init
func init() {
	viper_server.Init(config.GetViperConfig())
}

type WebBaseFunc interface {
	AddWebStatic(staticAbsPath, webPrefix string, paths ...string)
	AddUploadStatic(staticAbsPath, webPrefix string)
	InitRouter() error
	Run()
}

type WebFunc interface {
	WebBaseFunc
}

// Start
func Start(wf WebFunc) {
	err := wf.InitRouter()
	if err != nil {
		_logUtils.Error(err.Error())
		return
	}
	wf.Run()
}

func StartTest(wf WebFunc) {
	err := wf.InitRouter()
	if err != nil {
		_logUtils.Error(err.Error())
	}
}
