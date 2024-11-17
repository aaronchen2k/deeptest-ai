package service

import (
	"encoding/json"
	"errors"
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	_http "github.com/deeptest-com/deeptest-next/pkg/libs/http"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	_str "github.com/deeptest-com/deeptest-next/pkg/libs/string"
	"github.com/fatih/color"
	"io"
	"net/http"
	"strings"
	"time"
)

type NlpService struct {
}

func (s *NlpService) Parse(req v1.NlpReq) (ret v1.NlpResp, err error) {
	if req.Model == "" {
		req.Model = "qwen2.5-coder:1.5b-instruct"
	}

	url := _http.AddSepIfNeeded(config.CONFIG.System.LLmUrl) + "v1/chat/completions"
	_logUtils.Info("url=" + url)

	dataBytes, err := json.Marshal(req.Instruction)
	if err != nil {
		_logUtils.Infof(color.RedString("marshal request failed, error: %s.", err.Error()))
		return
	}

	dataStr := string(dataBytes)

	request, err := http.NewRequest("POST", url, strings.NewReader(dataStr))
	if err != nil {
		_logUtils.Infof(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		_logUtils.Infof(color.RedString("request failed, error: %s.", err.Error()))
		return
	}

	if !_http.IsSuccessCode(resp.StatusCode) {
		_logUtils.Infof(color.RedString("post request return %d - '%s'.", resp.StatusCode, resp.Status))
		err = errors.New(resp.Status)
		return
	}

	unicodeContent, _ := io.ReadAll(resp.Body)
	bytes, err := _str.UnescapeUnicode(unicodeContent)
	if err != nil {
		_logUtils.Infof(color.RedString("request failed, error: %s.", err.Error()))
		return
	}

	err = json.Unmarshal(bytes, &ret)
	if err != nil {
		_logUtils.Infof(color.RedString("unmarshal failed, error: %s.", err.Error()))
		return
	}

	return
}
