package service

import (
	"errors"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	_date "github.com/deeptest-com/deeptest-next/pkg/libs/date"
	_file "github.com/deeptest-com/deeptest-next/pkg/libs/file"
	_logs "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/dir"
)

var (
	ErrEmpty = errors.New("请上传正确的文件")
)

type FileService struct {
}

// UploadFile 上传文件
func (s *FileService) UploadFile(ctx iris.Context, fh *multipart.FileHeader) (ret string, err error) {
	filename, err := _file.GetUploadFileName(fh.Filename)
	if err != nil {
		_logs.Errorf("获取文件名失败，错误%s", err.Error())
		return
	}

	targetDir := filepath.Join(consts.DirUpload, _date.DateStr(time.Now()))
	absDir := filepath.Join(consts.WorkDir, targetDir)
	dir.InsureDir(absDir)

	err = dir.InsureDir(targetDir)
	if err != nil {
		_logs.Errorf("文件上传失败，错误：%s", err.Error())
		return
	}

	pth := filepath.Join(absDir, filename)
	_, err = ctx.SaveFormFile(fh, pth)
	if err != nil {
		_logs.Errorf("文件保存到本地失败，错误：%s", err.Error())
		return
	}

	ret = filepath.Join(targetDir, filename)

	return
}
