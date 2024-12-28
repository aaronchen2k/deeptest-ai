package service

import (
	"encoding/json"
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	httpUtils "github.com/deeptest-com/deeptest-next/internal/pkg/libs/http"
	_file "github.com/deeptest-com/deeptest-next/pkg/libs/file"
	"github.com/deeptest-com/deeptest-next/pkg/libs/http"
	"github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"github.com/snowlyg/helper/dir"
	"os"
	"path/filepath"
	"strings"
)

type KnowledgeBaseService struct {
}

var (
	defaultDb = "b0b12d74-2f56-49a8-9fad-8f5c6919b85e"

	kbCreateDocUri = "datasets/%s/document/create-by-file"
	kbQueryDocUri  = "datasets/%s/documents"
	kbRemoveDocUri = "datasets/%s/documents/%s"
)

func (s *KnowledgeBaseService) UnzipAndUploadFiles(zipPath, kb string) (err error) {
	unzipDir, err := _file.Unzip(zipPath, filepath.Join(consts.WorkDir, "_temp"))
	if err != nil {
		return
	}

	uploadFiles, _ := s.backupKbFiles(unzipDir, unzipDir, "img")
	for _, filePath := range uploadFiles {
		err := s.uploadDocToKnowledgeBase(filePath, kb)

		if err != nil {
			continue
		}
	}

	return
}

func (s *KnowledgeBaseService) uploadDocToKnowledgeBase(pth, kb string) (err error) {
	if kb == "" {
		kb = defaultDb
	}

	url := ""
	if config.CONFIG.Ai.PlatformType == consts.Dify {
		url = _http.AddSepIfNeeded(config.CONFIG.Ai.PlatformUrl) +
			fmt.Sprintf(kbCreateDocUri, kb)
	}
	_logs.Infof("%s url = %s", config.CONFIG.Ai.PlatformType, url)

	data := s.getData()

	bts, err := httpUtils.PostFile(url, data, pth, s.getHeaders())
	if err != nil {
		return
	}
	_logs.Infof("create doc resp %s", string(bts))

	return
}

func (s *KnowledgeBaseService) ClearAll(kb string) (err error) {
	if kb == "" {
		kb = defaultDb
	}

	queryUrl := ""
	if config.CONFIG.Ai.PlatformType == consts.Dify {
		queryUrl = _http.AddSepIfNeeded(config.CONFIG.Ai.PlatformUrl) +
			fmt.Sprintf(kbQueryDocUri, kb)
		queryUrl += "?limit=100"
	}
	_logs.Infof("%s queryUrl = %s", config.CONFIG.Ai.PlatformType, queryUrl)

	headers := s.getHeaders()
	bts, err := httpUtils.Get(queryUrl, headers)
	if err != nil {
		return
	}

	docs := domain.KbQueryResult{}
	json.Unmarshal(bts, &docs)

	for _, doc := range docs.Data {
		removeUrl := ""
		if config.CONFIG.Ai.PlatformType == consts.Dify {
			removeUrl = _http.AddSepIfNeeded(config.CONFIG.Ai.PlatformUrl) +
				fmt.Sprintf(kbRemoveDocUri, kb, doc.Id)
		}
		_logs.Infof("%s removeUrl = %s", config.CONFIG.Ai.PlatformType, removeUrl)

		bts, err = httpUtils.Delete(removeUrl, headers)
		if err != nil {
			continue
		}
	}

	return
}

func (s *KnowledgeBaseService) getHeaders() (ret map[string]string) {
	ret = map[string]string{"Authorization": "Bearer " + os.Getenv("AI_DATASET_API_KEY")}
	return
}

func (s *KnowledgeBaseService) backupKbFiles(root, dirName, exclude string) (ret []string, err error) {
	dirKbDocs := filepath.Join(consts.WorkDir, consts.DirKbDocs)

	dirName = strings.TrimSuffix(dirName, _file.PathSep)

	infos, err := os.ReadDir(dirName)
	if err != nil {
		return
	}

	for _, info := range infos {
		path := filepath.Join(dirName, info.Name())
		realInfo, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		if realInfo.Name() == "__MACOSX" {
			continue
		}

		if realInfo.IsDir() {
			children, err := s.backupKbFiles(root, path, exclude)
			if err != nil {
				return nil, err
			}
			ret = append(ret, children...)

		} else {
			if info.Name() != exclude {
				ret = append(ret, path)
			}

			dist := filepath.Join(dirKbDocs, strings.TrimPrefix(path, root))

			_file.InsureDir(dir.Dir(dist))
			err = _file.CopyFile(path, dist)
			if err != nil {
				continue
			}
		}

	}
	return
}

func (s *KnowledgeBaseService) getData() (ret domain.KbCreateReq) {
	rules := domain.Rules{Segmentation: domain.Segmentation{
		Separator: "###", MaxTokens: 500,
	}}
	rules.PreProcessingRules = append(rules.PreProcessingRules, domain.PreProcessingRule{
		Id:      "remove_extra_spaces",
		Enabled: true,
	})
	rules.PreProcessingRules = append(rules.PreProcessingRules, domain.PreProcessingRule{
		Id:      "remove_urls_emails",
		Enabled: true,
	})

	ret = domain.KbCreateReq{
		IndexingTechnique: "high_quality",
		ProcessRule: domain.ProcessRule{
			Rules: rules,
			Mode:  "custom",
		},
	}

	return
}
