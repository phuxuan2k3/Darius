package handler

import (
	"darius/internal/services/bulbasaur"
	llm "darius/internal/services/llm"
	"darius/internal/services/missfortune"
	llmManager "darius/managers/llm"
	hello "darius/pkg/proto/hello"
	suggest "darius/pkg/proto/suggest"
)

type Dependency struct {
	LlmService  llm.LLM
	LLMManager  llmManager.Manager
	Missfortune missfortune.Service
	Bulbasaur   bulbasaur.Service
}

type handler struct {
	hello.HelloServiceServer
	suggest.SuggestServiceServer

	llmService  llm.LLM
	llmManager  llmManager.Manager
	missfortune missfortune.Service
	bulbasaur   bulbasaur.Service
}

func NewHandlerWithDeps(deps Dependency) *handler {
	return &handler{
		llmService:  deps.LlmService,
		llmManager:  deps.LLMManager,
		missfortune: deps.Missfortune,
		bulbasaur:   deps.Bulbasaur,
	}
}
