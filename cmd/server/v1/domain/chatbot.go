package v1

import "github.com/deeptest-com/deeptest-next/internal/pkg/consts"

type ChatReq struct {
	Query          string                 `json:"query"`
	Inputs         interface{}            `json:"inputs"`
	ResponseMode   consts.LlmResponseMode `json:"response_mode"`
	HasThoughts    bool                   `json:"has_thoughts"`
	User           string                 `json:"user"`
	ConversationId string                 `json:"conversation_id"`

	Temperature float64 `json:"temperature"`
}

type ChatRespBlocking struct {
	Event          string `json:"event"`
	TaskId         string `json:"task_id"`
	Id             string `json:"id"`
	MessageId      string `json:"message_id"`
	ConversationId string `json:"conversation_id"`
	Mode           string `json:"mode"`
	Answer         string `json:"answer"`
	Metadata       struct {
		RetrieverResources []struct {
			DatasetId      string  `json:"dataset_id"`
			DatasetName    string  `json:"dataset_name"`
			DocumentId     string  `json:"document_id"`
			DocumentName   string  `json:"document_name"`
			DataSourceType string  `json:"data_source_type"`
			SegmentId      string  `json:"segment_id"`
			RetrieverFrom  string  `json:"retriever_from"`
			Score          float64 `json:"score"`
			Content        string  `json:"content"`
			Position       int     `json:"position"`
		} `json:"retriever_resources"`
		Usage struct {
			PromptTokens        int     `json:"prompt_tokens"`
			PromptUnitPrice     string  `json:"prompt_unit_price"`
			PromptPriceUnit     string  `json:"prompt_price_unit"`
			PromptPrice         string  `json:"prompt_price"`
			CompletionTokens    int     `json:"completion_tokens"`
			CompletionUnitPrice string  `json:"completion_unit_price"`
			CompletionPriceUnit string  `json:"completion_price_unit"`
			CompletionPrice     string  `json:"completion_price"`
			TotalTokens         int     `json:"total_tokens"`
			TotalPrice          string  `json:"total_price"`
			Currency            string  `json:"currency"`
			Latency             float64 `json:"latency"`
		} `json:"usage"`
	} `json:"metadata"`
	CreatedAt int `json:"created_at"`
}

type ChatRespStreamingAnswer struct {
	Event                string      `json:"event,omitempty"`
	ConversationId       string      `json:"conversation_id,omitempty"`
	MessageId            string      `json:"message_id,omitempty"`
	CreatedAt            int         `json:"created_at,omitempty"`
	TaskId               string      `json:"task_id,omitempty"`
	Id                   string      `json:"id,omitempty"`
	Answer               string      `json:"answer"`
	FromVariableSelector interface{} `json:"from_variable_selector,omitempty"`
}
type ChatRespStreamingData struct {
	Event          string `json:"event"`
	ConversationId string `json:"conversation_id"`
	MessageId      string `json:"message_id"`
	CreatedAt      int    `json:"created_at"`
	TaskId         string `json:"task_id"`
	Id             string `json:"id"`
	Metadata       struct {
		RetrieverResources []struct {
			DatasetId      string  `json:"dataset_id"`
			DatasetName    string  `json:"dataset_name"`
			DocumentId     string  `json:"document_id"`
			DocumentName   string  `json:"document_name"`
			DataSourceType string  `json:"data_source_type"`
			SegmentId      string  `json:"segment_id"`
			RetrieverFrom  string  `json:"retriever_from"`
			Score          float64 `json:"score"`
			Content        string  `json:"content"`
			Position       int     `json:"position"`
		} `json:"retriever_resources"`
		Usage struct {
			PromptTokens        int     `json:"prompt_tokens"`
			PromptUnitPrice     string  `json:"prompt_unit_price"`
			PromptPriceUnit     string  `json:"prompt_price_unit"`
			PromptPrice         string  `json:"prompt_price"`
			CompletionTokens    int     `json:"completion_tokens"`
			CompletionUnitPrice string  `json:"completion_unit_price"`
			CompletionPriceUnit string  `json:"completion_price_unit"`
			CompletionPrice     string  `json:"completion_price"`
			TotalTokens         int     `json:"total_tokens"`
			TotalPrice          string  `json:"total_price"`
			Currency            string  `json:"currency"`
			Latency             float64 `json:"latency"`
		} `json:"usage"`
	} `json:"metadata"`
	Files interface{} `json:"files"`
}
