package llm

type LLMRequest struct {
	Model   string `json:"model"`
	Content string `json:"content"`
}

type LLMResponse struct {
	ConversationID string `json:"conversation_id"`
	Created_at     string `json:"created_at"`
	Content        string `json:"content"`
}
