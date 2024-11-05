package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/deeptest-com/deeptest-next/internal/pkg/go-openai2"
)

func main() {
	config := openai2.DefaultConfig("EMPTY")
	config.BaseURL = "http://192.168.5.134:7861/knowledge_base/local_kb/wiki"
	client := openai2.NewClientWithConfig(config)

	stream, err := client.CreateChatCompletionStream(
		context.Background(),

		openai2.ChatCompletionRequest{
			Model:  "glm4-chat",
			Stream: true,
			//N:           3,
			Temperature: 0.7,
			Messages: []openai2.ChatCompletionMessage{
				{
					Role:    openai2.ChatMessageRoleUser,
					Content: "你好",
				},
				{
					Role:    openai2.ChatMessageRoleAssistant,
					Content: "你好，我是人工智能大模型",
				},
				{
					Role:    openai2.ChatMessageRoleUser,
					Content: "提取器怎么用？",
				},
			},
			ExtraBody: openai2.ExtraBody{
				TopK:           3,
				ScoreThreshold: 2.0,
				ReturnDirect:   false,
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	defer stream.Close()

	fmt.Printf("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			break
		}
		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			break
		}

		fmt.Printf("%v", response.Choices)
	}
}
