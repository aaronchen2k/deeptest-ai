package serve

import (
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	_file "github.com/deeptest-com/deeptest-next/pkg/libs/file"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/snowlyg/helper/dir"
	"path/filepath"
)

func AddStatic(app *iris.Application) {
	addWebUi(app)
	addUpload(app)
}

// addWebUi 添加前端页面访问
func addWebUi(app *iris.Application) {
	pth := filepath.Join(dir.GetCurrentAbPath(), "deeptest-ui")

	//fileUtils.MkDirIfNeeded(pth)
	_logUtils.Infof("*** ui dir: %s", pth)

	app.HandleDir("/", iris.Dir(pth), iris.DirOptions{
		IndexName: "index.html",
		ShowList:  false,
		SPA:       true,
	})
}

// addUpload 添加上传文件访问
func addUpload(app *iris.Application) {
	pth := filepath.Join(consts.WorkDir, consts.DirUpload)
	_file.InsureDir(pth)
	_logUtils.Infof("*** upload dir: %s", pth)

	app.HandleDir("/upload", iris.Dir(pth), router.DirOptions{Attachments: router.Attachments{Enable: true}})
}
