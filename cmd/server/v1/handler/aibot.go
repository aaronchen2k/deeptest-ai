package handler

import (
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/service"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"

	"github.com/kataras/iris/v12"
)

type ChatbotCtrl struct {
	BaseCtrl
	ChatbotService *service.ChatbotService `inject:""`
}

func (c *ChatbotCtrl) ChatCompletion(ctx iris.Context) {
	flusher, ok := ctx.ResponseWriter().Flusher()
	if !ok {
		ctx.StopWithText(iris.StatusHTTPVersionNotSupported, "Streaming unsupported!")
		return
	}

	ctx.ContentType("text/event-stream")
	//ctx.Header("content-type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")

	req := v1.ChatReq{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	c.ChatbotService.ChatCompletion(req, flusher, ctx)
}
