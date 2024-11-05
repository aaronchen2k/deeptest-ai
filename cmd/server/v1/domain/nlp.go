package v1

type NlpReq struct {
	Instruction string `json:"instruction"`
	Model       string `json:"model"`
}

type NlpResp struct {
	Intent string                 `json:"intent"`
	Slots  map[string]interface{} `json:"slots"`
}
