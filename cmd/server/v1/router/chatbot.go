package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	"github.com/kataras/iris/v12"
)

type ChatbotModule struct {
	ChatbotCtrl *handler.ChatbotCtrl `inject:""`
}

func (m *ChatbotModule) Party() func(public iris.Party) {
	return func(party iris.Party) {
		party.Post("/chat", m.ChatbotCtrl.Chat).Name = "聊天"
	}
}
