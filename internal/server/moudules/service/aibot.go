package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	v1 "github.com/deeptest-com/deeptest-next/cmd/server/v1/domain"
	"github.com/deeptest-com/deeptest-next/internal/pkg/config"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	_http "github.com/deeptest-com/deeptest-next/pkg/libs/http"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"github.com/kataras/iris/v12"
	"io"
	"net/http"
	"strings"
)

type ChatbotService struct {
}

func (s *ChatbotService) ChatCompletion(req v1.ChatCompletionReq, flusher http.Flusher, ctx iris.Context) (ret _domain.PageData, err error) {
	if len(req.Messages) > 0 && strings.TrimSpace(req.Messages[len(req.Messages)-1].Content) == "小深" {
		str := s.genResp("您好，有什么可以帮助您的？")

		ctx.Writef("%s\n\n", str)
		flusher.Flush()
		return
	}

	req.Model = "qwen2.5-coder:1.5b-instruct"

	url := _http.AddSepIfNeeded(config.CONFIG.System.ChatchatUrl) + "chat/chat/completions"
	_logUtils.Info("url=" + url)
	bts, err := json.Marshal(req)

	reader := bytes.NewReader(bts)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return
	}

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
		str := string(bytes)

		if err1 != nil && err1 != io.EOF {
			err = err1
			break
		}
		if err1 == io.EOF {
			break
		}

		fmt.Println("\n>>>" + str + "\n")

		// must with prefix "data:" for openai response
		// must add a postfix "\n\n"
		ctx.Writef("%s\n\n", str)
		flusher.Flush()
	}

	return
}

func (s *ChatbotService) KnowledgeBaseChat(req v1.KnowledgeBaseChatReq, flusher http.Flusher, ctx iris.Context) (ret _domain.PageData, err error) {
	if len(req.Messages) > 0 && strings.TrimSpace(req.Messages[len(req.Messages)-1].Content) == "小深" {
		str := s.genResp("您好，有什么可以帮助您的？")

		ctx.Writef("%s\n\n", str)
		flusher.Flush()
		return
	}

	if strings.TrimSpace(req.Model) == "" {
		req.Model = "glm4-chat"

		//models, _ := s.ListValidModel("llm")
		//
		//if len(models) > 0 {
		//	req.Model = models[0].ModelName
		//} else {
		//	str := s.genResp("没有可使用的大模型，请联系管理员。")
		//
		//	ctx.Writef("%s\n\n", str)
		//	flusher.Flush()
		//	return
		//}
	}

	url := _http.AddSepIfNeeded(config.CONFIG.System.ChatchatUrl) +
		fmt.Sprintf("knowledge_base/local_kb/%s/chat/completions", req.KbName)
	bts, err := json.Marshal(req)

	reader := bytes.NewReader(bts)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return
	}

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

	req.KbName = ""
	r := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	for {
		bytes, err1 := r.ReadBytes('\n')
		str := string(bytes)

		if err1 != nil && err1 != io.EOF {
			err = err1
			break
		}
		if err1 == io.EOF {
			break
		}

		fmt.Println("\n>>>" + str + "\n")

		// must with prefix "data:" for openai response
		// must add a postfix "\n\n"
		ctx.Writef("%s\n\n", str)
		flusher.Flush()
	}

	return
}

func (s *ChatbotService) ListValidModel(typ string) (ret []v1.ChatchatModelData, err error) {
	url := _http.AddSepIfNeeded(config.CONFIG.System.ChatchatUrl) + "v1/models"

	bytes, err := _http.Get(url)
	if err != nil {
		return
	}

	resp := v1.ChatchatModelResp{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return
	}

	for _, item := range resp.Data {
		if strings.ToLower(item.ModelType) == typ {
			ret = append(ret, item)
		}
	}

	return
}

func (s *ChatbotService) ListKnowledgeBase() (ret []v1.ChatchatKnowledgeBaseData, err error) {
	url := _http.AddSepIfNeeded(config.CONFIG.System.ChatchatUrl) + "knowledge_base/list_knowledge_bases"

	bytes, err := _http.Get(url)
	if err != nil {
		return
	}

	resp := v1.ChatchatKnowledgeBaseResp{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return
	}

	ret = resp.Data

	return
}

func (s *ChatbotService) genResp(content string) (ret string) {
	resp := v1.ChatchatResponse{}
	choice := v1.ChatchatChoice{
		Delta: v1.ChatchatDelta{
			Content: content,
		},
	}
	resp.Choices = append(resp.Choices, choice)

	bytes, _ := json.Marshal(resp)

	ret = "data:" + string(bytes)

	return
}
