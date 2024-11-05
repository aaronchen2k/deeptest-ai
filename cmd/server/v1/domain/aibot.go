package v1

type ChatCompletionReq struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream      interface{} `json:"stream"`
	Temperature float64     `json:"temperature"`
}

type KnowledgeBaseChatReq struct {
	KbName   string `json:"kb_name,omitempty"`
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream      interface{} `json:"stream"`
	Temperature float64     `json:"temperature"`
	ExtraBody   struct {
		TopK           int         `json:"top_k"`
		ScoreThreshold float64     `json:"score_threshold"`
		ReturnDirect   interface{} `json:"return_direct"`
	} `json:"extra_body"`
}

type ChatchatModelReq struct {
	Placeholder string `json:"placeholder"`
}

type ChatchatResponse struct {
	Choices []ChatchatChoice `json:"choices"`
}
type ChatchatChoice struct {
	Delta ChatchatDelta `json:"delta"`
}
type ChatchatDelta struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type T struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Id         int    `json:"id"`
		KbName     string `json:"kb_name"`
		KbInfo     string `json:"kb_info"`
		VsType     string `json:"vs_type"`
		EmbedModel string `json:"embed_model"`
		FileCount  int    `json:"file_count"`
		CreateTime string `json:"create_time"`
	} `json:"data"`
}

type ChatchatKnowledgeBaseResp struct {
	Code int                         `json:"code"`
	Msg  string                      `json:"msg"`
	Data []ChatchatKnowledgeBaseData `json:"data"`
}
type ChatchatModelResp struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data []ChatchatModelData `json:"data"`
}

type ChatchatKnowledgeBaseData struct {
	Id         int    `json:"id"`
	KbName     string `json:"kb_name"`
	KbInfo     string `json:"kb_info"`
	VsType     string `json:"vs_type"`
	EmbedModel string `json:"embed_model"`
	FileCount  int    `json:"file_count"`
	CreateTime string `json:"create_time"`
}
type ChatchatModelData struct {
	Id            string   `json:"id"`
	Created       int      `json:"created"`
	Object        string   `json:"object"`
	OwnedBy       string   `json:"owned_by"`
	ModelType     string   `json:"model_type"`
	Address       string   `json:"address"`
	Accelerators  []string `json:"accelerators"`
	ModelName     string   `json:"model_name"`
	Dimensions    int      `json:"dimensions"`
	MaxTokens     int      `json:"max_tokens"`
	Language      []string `json:"language"`
	ModelRevision string   `json:"model_revision"`
	Replica       int      `json:"replica"`
	PlatformName  string   `json:"platform_name"`
}
