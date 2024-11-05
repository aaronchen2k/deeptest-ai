package router

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/handler"
	"github.com/kataras/iris/v12"
)

type AibotModule struct {
	AibotCtrl *handler.AibotCtrl `inject:""`
}

func (m *AibotModule) Party() func(public iris.Party) {
	return func(party iris.Party) {
		party.Post("/chat_completion", m.AibotCtrl.ChatCompletion).Name = "聊天"
		party.Post("/knowledge_base_chat", m.AibotCtrl.KnowledgeBaseChat).Name = "与知识库对话"

		party.Get("/list_valid_model", m.AibotCtrl.ListValidModel).Name = "列出可用的大模型"
		party.Get("/list_knowledge_base", m.AibotCtrl.ListKnowledgeBase).Name = "列出可用的知识库"
	}
}
