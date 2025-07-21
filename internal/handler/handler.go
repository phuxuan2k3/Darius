package handler

import (
	"context"
	"darius/internal/services/bulbasaur"
	llm "darius/internal/services/llm"
	"darius/internal/services/missfortune"
	llmManager "darius/managers/llm"
	"darius/pkg/proto/suggest"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Dependency struct {
	LlmService  llm.LLM
	LLMManager  llmManager.Manager
	Missfortune missfortune.Service
	Bulbasaur   bulbasaur.Service
}

type handler struct {
	suggest.UnimplementedSuggestServiceServer

	llmService  llm.LLM
	llmManager  llmManager.Manager
	missfortune missfortune.Service
	bulbasaur   bulbasaur.Service

	cache map[string]interface{}
}

func NewHandlerWithDeps(deps Dependency) *handler {
	return &handler{
		llmService:  deps.LlmService,
		llmManager:  deps.LLMManager,
		missfortune: deps.Missfortune,
		bulbasaur:   deps.Bulbasaur,
		cache:       make(map[string]interface{}),
	}
}

func (h *handler) mustEmbedUnimplementedSuggestServiceServer() {}

type x struct {
}

func (x) SuggestCriteria(context.Context, *suggest.SuggestCriteriaRequest) (*suggest.SuggestCriteriaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SuggestCriteria not implemented")
}
func (x) SuggestOptions(context.Context, *suggest.SuggestOptionsRequest) (*suggest.SuggestOptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SuggestOptions not implemented")
}
func (x) SuggestQuestions(context.Context, *suggest.SuggestQuestionsRequest) (*suggest.SuggestExamQuestionResponseV2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SuggestQuestions not implemented")
}
func (x) SuggestInterviewQuestion(context.Context, *suggest.SuggestInterviewQuestionRequest) (*suggest.SuggestInterviewQuestionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SuggestInterviewQuestion not implemented")
}
func (x) ScoreInterview(context.Context, *suggest.ScoreInterviewRequest) (*suggest.ScoreInterviewResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ScoreInterview not implemented")
}
func (x) SuggestOutlines(context.Context, *suggest.SuggestOutlinesRequest) (*suggest.SuggestOutlinesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SuggestOutlines not implemented")
}
func (x) SuggestExamQuestionV2(context.Context, *suggest.SuggestExamQuestionRequest) (*suggest.SuggestExamQuestionResponseV2, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SuggestExamQuestionV2 not implemented")
}
func (x) mustEmbedUnimplementedSuggestServiceServer() {}
