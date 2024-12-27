package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	"github.com/kataras/iris/v12"
)

type FileModule struct {
	FileCtrl *handler.FileCtrl `inject:""`
}

func (m *FileModule) Party() func(public iris.Party) {
	return func(party iris.Party) {
		party.Post("/upload", iris.LimitRequestBodySize(100*iris.MB),
			m.FileCtrl.Upload).Name = "上传文件"
	}
}
