package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/deeptest-com/deeptest-next"
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	"github.com/deeptest-com/deeptest-next/pkg/domain"
	"github.com/deeptest-com/deeptest-next/pkg/libs/http"
	"github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"github.com/kataras/iris/v12"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

type ChatbotService struct {
}

func (s *ChatbotService) Chat(req v1.ChatReq, flusher http.Flusher, ctx iris.Context) (ret _domain.PageData, err error) {
	url := ""
	if config.CONFIG.Ai.PlatformType == consts.Dify {
		url = _http.AddSepIfNeeded(config.CONFIG.Ai.PlatformUrl) + "v1/chat-messages"
	}
	_logs.Infof("%s url = %s", config.CONFIG.Ai.PlatformType, url)

	req.Query = s.getTemplateResult(req.Query, consts.TemplateKbQuery)

	bts, err := json.Marshal(req)
	reader := bytes.NewReader(bts)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.CONFIG.Ai.ApiKey))
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	transport := &http.Transport{}
	transport.DisableCompression = true
	client.Transport = transport

	resp, err := client.Do(request)
	if err != nil {
		return
	}

	r := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	for {
		bytes, err1 := r.ReadBytes('\n')
		if err1 != nil && err1 != io.EOF {
			err = err1
			break
		}
		if err1 == io.EOF {
			break
		}

		fmt.Println("\n====== " + string(bytes))
		str := s.genResp(bytes, req.ResponseMode)
		if str == "" {
			continue
		}

		fmt.Println("------ " + str + "\n")

		// must with prefix "data:" for openai response
		// must add a postfix "\n\n"
		ctx.Writef("%s\n\n", str)
		flusher.Flush()
	}

	return
}

func (s *ChatbotService) genResp(input []byte, typ consts.LlmResponseMode) (ret string) {
	str := strings.TrimSpace(string(input))
	if str == "" {
		return
	}

	str = strings.TrimPrefix(str, "data:")

	output := make([]byte, 0)
	if typ == consts.Blocking {
		resp := v1.ChatRespBlocking{}
		json.Unmarshal([]byte(str), &resp)
		output, _ = json.Marshal(resp)

	} else if typ == consts.Streaming {
		resp := v1.ChatRespStreamingAnswer{}
		json.Unmarshal([]byte(str), &resp)

		if resp.Answer != "" { // answer
			simple := v1.ChatRespStreamingAnswer{
				Answer:         resp.Answer,
				ConversationId: resp.ConversationId,
			}
			output, _ = json.Marshal(simple)
		} else { // data
			resp := v1.ChatRespStreamingData{}
			json.Unmarshal([]byte(str), &resp)

			output, _ = json.Marshal(resp)
		}
	}

	ret = "data:" + string(output)

	return
}

func (s *ChatbotService) getTemplateResult(query, tmpl string) (ret string) {
	if tmpl != "" {
		pth := filepath.Join("res", "tmpl", tmpl+".tmpl")

		bts, err := deeptest.ReadResData(pth)
		if err == nil {
			str := string(bts)
			ret = fmt.Sprintf(str, query)
		}
	}

	return
}
