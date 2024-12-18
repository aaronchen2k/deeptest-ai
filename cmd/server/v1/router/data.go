package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	"github.com/kataras/iris/v12"
)

type DataModule struct {
	DataCtrl *handler.DataCtrl `inject:""`
}

func (m *DataModule) Party() func(public iris.Party) {
	return func(party iris.Party) {
		party.Post("/initdb", m.DataCtrl.Init).Name = "InitDB"
		party.Post("/checkdb", m.DataCtrl.Init).Name = "CheckDB"
	}
}
