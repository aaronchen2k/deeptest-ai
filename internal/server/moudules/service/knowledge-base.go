package service

import (
	"encoding/json"
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/internal/pkg/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/libs/http"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/repo"
	"github.com/deeptest-com/deeptest-next/pkg/libs/file"
	"github.com/deeptest-com/deeptest-next/pkg/libs/http"
	"github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"github.com/snowlyg/helper/dir"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type KnowledgeBaseService struct {
	KnowledgeBaseRepo *repo.KnowledgeBaseRepo `inject:""`
}

var (
	kbCreateDocUri = "datasets/%s/document/create-by-file"
	kbQueryDocUri  = "datasets/%s/documents"
	kbRemoveDocUri = "datasets/%s/documents/%s"
)

func (s *KnowledgeBaseService) UnzipAndUploadFiles(zipPath, dataset string) (err error) {
	material := model.KbMaterial{File: zipPath}
	s.KnowledgeBaseRepo.SaveMaterial(&material)

	tempDir := filepath.Join(consts.WorkDir, "_temp")
	unzipDir, err := _file.Unzip(zipPath, tempDir)
	if err != nil {
		return
	}

	uploadFiles, _ := s.backupKbFiles(unzipDir, unzipDir, "img", material.ID)

	for _, filePath := range uploadFiles {
		err := s.uploadDocToKnowledgeBase(filePath, dataset)

		if err != nil {
			continue
		}
	}

	// remove temporary files
	// ensure it is a dir that unzip from an uploaded file.
	// 26 is the length of uuid in unzip dir name.
	if len(unzipDir)-len(tempDir) > 26 {
		//os.RemoveAll(unzipDir)
	}

	return
}

func (s *KnowledgeBaseService) uploadDocToKnowledgeBase(pth, dataset string) (err error) {
	if dataset == "" {
		dataset = os.Getenv("AI_DATASET_ID")
	}

	url := ""
	if config.CONFIG.Ai.PlatformType == consts.Dify {
		url = _http.AddSepIfNeeded(config.CONFIG.Ai.PlatformUrl) +
			fmt.Sprintf(kbCreateDocUri, dataset)
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

func (s *KnowledgeBaseService) ClearAll(dataset string) (err error) {
	if dataset == "" {
		dataset = os.Getenv("AI_DATASET_ID")
	}

	queryUrl := ""
	if config.CONFIG.Ai.PlatformType == consts.Dify {
		queryUrl = _http.AddSepIfNeeded(config.CONFIG.Ai.PlatformUrl) +
			fmt.Sprintf(kbQueryDocUri, dataset)
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
				fmt.Sprintf(kbRemoveDocUri, dataset, doc.Id)
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

func (s *KnowledgeBaseService) backupKbFiles(root, dirName, exclude string, materialId uint) (ret []string, err error) {
	infos, err := os.ReadDir(dirName)
	if err != nil {
		return
	}

	for _, info := range infos {
		srcPath := filepath.Join(dirName, info.Name())
		realInfo, err := os.Stat(srcPath)
		if err != nil {
			return nil, err
		}

		if realInfo.Name() == "__MACOSX" {
			continue
		}

		if realInfo.IsDir() {
			children, err := s.backupKbFiles(root, srcPath, exclude, materialId)
			if err != nil {
				return nil, err
			}
			ret = append(ret, children...)

		} else {
			distDocDir := filepath.Join(consts.WorkDir, consts.DirKnowledgeBases,
				strconv.Itoa(int(materialId)))

			distName := strings.TrimPrefix(srcPath, root)
			if _file.GetExtName(distName) == ".md" {
				postFix := fmt.Sprintf("%d", materialId)
				distName = _file.AddFileNamePostfix(distName, postFix)
			}
			distPath := filepath.Join(distDocDir, distName)
			_file.InsureDir(dir.Dir(distPath))

			// copy
			err = _file.CopyFile(srcPath, distPath)
			if err != nil {
				continue
			}

			// update
			if _file.GetExtName(distPath) == ".md" {
				s.UpdateDoc(distPath, materialId)
			}
			if info.Name() != exclude {
				ret = append(ret, distPath)
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

func (s *KnowledgeBaseService) UpdateDoc(pth string, materialId uint) (err error) {
	content := _file.ReadFile(pth)

	prefix := fmt.Sprintf("%d/", materialId)

	reg, err := regexp.Compile(`(\[.+?\])\(([^http].+?)\)`)
	if reg == nil || err != nil {
		return
	}
	content = reg.ReplaceAllString(content, fmt.Sprintf(`$1(%s$2)`, prefix))

	_file.WriteBytes(pth, []byte(content))

	return
}
