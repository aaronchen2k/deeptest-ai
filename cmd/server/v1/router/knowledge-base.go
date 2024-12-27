package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	"github.com/kataras/iris/v12"
)

type KnowledgeBaseModule struct {
	KnowledgeBaseCtrl *handler.KnowledgeBaseCtrl `inject:""`
}

func (m *KnowledgeBaseModule) Party() func(public iris.Party) {
	return func(party iris.Party) {
		party.Post("/uploadDoc", iris.LimitRequestBodySize(100*iris.MB),
			m.KnowledgeBaseCtrl.UploadDoc).Name = "上传文件"
		party.Post("/clearAll", m.KnowledgeBaseCtrl.ClearAll).Name = "清除知识库文件"
	}
}
