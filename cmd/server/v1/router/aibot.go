package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	"github.com/kataras/iris/v12"
)

type ChatbotModule struct {
	ChatbotCtrl *handler.AibotCtrl `inject:""`
}

func (m *ChatbotModule) Party() func(public iris.Party) {
	return func(party iris.Party) {
		party.Post("/chat_completion", m.ChatbotCtrl.ChatCompletion).Name = "聊天"

		party.Get("/list_valid_model", m.ChatbotCtrl.ListValidModel).Name = "列出可用的大模型"
		party.Get("/list_knowledge_base", m.ChatbotCtrl.ListKnowledgeBase).Name = "列出可用的知识库"
	}
}
