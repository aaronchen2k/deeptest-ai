package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	middleware "github.com/deeptest-com/deeptest-next/internal/pkg/core/middleware"
	"github.com/kataras/iris/v12"
)

type UserModule struct {
	UserCtrl *handler.UserCtrl `inject:""`
}

func (m *UserModule) Party() func(index iris.Party) {
	return func(index iris.Party) {
		index.Use(middleware.JwtHandler(), middleware.Casbin())

		index.Get("/", m.UserCtrl.Paginate).Name = "用户列表"
		index.Get("/{id:uint}", m.UserCtrl.Get).Name = "用户详情"
		index.Post("/", m.UserCtrl.Create).Name = "创建用户"
		index.Post("/{id:uint}", m.UserCtrl.Update).Name = "编辑用户"
		index.Delete("/{id:uint}", m.UserCtrl.Delete).Name = "删除用户"
	}
}
