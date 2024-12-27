package service

import (
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/pkg/libs/http"
	"github.com/deeptest-com/deeptest-next/pkg/libs/log"
)

type KnowledgeBaseService struct {
}

var (
	kbCreateDocUri = "/datasets/%s/document/create-by-file"
	kbQueryDocUri  = "/datasets/%s/documents"
	kbRemoveDocUri = "/datasets/%s/documents/%s"
)

func (s *KnowledgeBaseService) UploadDoc(pth, kb string) (err error) {
	url := ""
	if config.CONFIG.Ai.PlatformType == consts.Dify {
		url = _http.AddSepIfNeeded(config.CONFIG.Ai.PlatformUrl) +
			fmt.Sprintf(kbCreateDocUri, kb)
	}
	_logs.Infof("%s url = %s", config.CONFIG.Ai.PlatformType, url)

	return
}

func (s *KnowledgeBaseService) ClearAll(kb string) (err error) {
	queryUrl := ""
	if config.CONFIG.Ai.PlatformType == consts.Dify {
		queryUrl = _http.AddSepIfNeeded(config.CONFIG.Ai.PlatformUrl) +
			fmt.Sprintf(kbQueryDocUri, kb)
	}
	_logs.Infof("%s queryUrl = %s", config.CONFIG.Ai.PlatformType, queryUrl)

	removeUrl := ""
	if config.CONFIG.Ai.PlatformType == consts.Dify {
		queryUrl = _http.AddSepIfNeeded(config.CONFIG.Ai.PlatformUrl) +
			fmt.Sprintf(kbRemoveDocUri, kb)
	}
	_logs.Infof("%s removeUrl = %s", config.CONFIG.Ai.PlatformType, removeUrl)

	return
}
