package handler

import (
	llm "darius/internal/services/llm"
	llmManager "darius/managers/llm"
	hello "darius/pkg/proto/hello"
	suggest "darius/pkg/proto/suggest"
)

type Dependency struct {
	LlmService llm.LLM
	LLMManager llmManager.Manager
}

type handler struct {
	hello.HelloServiceServer
	suggest.SuggestServiceServer

	llmService llm.LLM
	llmManager llmManager.Manager
}

func NewHandlerWithDeps(deps Dependency) *handler {
	return &handler{
		llmService: deps.LlmService,
		llmManager: deps.LLMManager,
	}
}
