package handler

import (
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/service"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	_logs "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

type FileCtrl struct {
	FileService    *service.FileService    `inject:""`
	ChatbotService *service.ChatbotService `inject:""`
}

func (c *FileCtrl) Upload(ctx iris.Context) {
	f, fh, err := ctx.FormFile("file")
	if err != nil {
		_logs.Errorf("文件上传失败", zap.String("ctx.FormFile(\"file\")", err.Error()))
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}
	defer f.Close()

	name := fh.Filename
	pth, err := c.FileService.UploadFile(ctx, fh)
	if err != nil {
		ctx.JSON(_domain.Response{Code: _domain.SystemErr.Code, Msg: err.Error()})
		return
	}

	c.ChatbotService.UploadToKb(pth)

	ctx.JSON(_domain.Response{Code: _domain.Success.Code, Data: iris.Map{"path": pth, "name": name}})
}
