package handler

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/service"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"

	"github.com/kataras/iris/v12"
)

type AibotCtrl struct {
	BaseCtrl
	AibotService *service.AibotService `inject:""`
}

func (c *AibotCtrl) ChatCompletion(ctx iris.Context) {
	flusher, ok := ctx.ResponseWriter().Flusher()
	if !ok {
		ctx.StopWithText(iris.StatusHTTPVersionNotSupported, "Streaming unsupported!")
		return
	}

	ctx.ContentType("text/event-stream")
	//ctx.Header("content-type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")

	req := v1.ChatCompletionReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	c.AibotService.ChatCompletion(req, flusher, ctx)
}

func (c *AibotCtrl) KnowledgeBaseChat(ctx iris.Context) {
	flusher, ok := ctx.ResponseWriter().Flusher()
	if !ok {
		ctx.StopWithText(iris.StatusHTTPVersionNotSupported, "Streaming unsupported!")
		return
	}

	ctx.ContentType("text/event-stream")
	//ctx.Header("content-type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")

	req := v1.KnowledgeBaseChatReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	c.AibotService.KnowledgeBaseChat(req, flusher, ctx)
}

func (c *AibotCtrl) ListValidModel(ctx iris.Context) {
	typ := ctx.URLParamDefault("type", "llm")

	data, err := c.AibotService.ListValidModel(typ)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: data})
}

func (c *AibotCtrl) ListKnowledgeBase(ctx iris.Context) {
	data, err := c.AibotService.ListKnowledgeBase()
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: data})
}
