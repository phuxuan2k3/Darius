package managers

import (
	"context"
	llm_grpc "darius/internal/services/llm-grpc"
	databaseService "darius/internal/services/repo"
	"darius/metrics"
	"log"
)

type Manager interface {
	Generate(context.Context, string, string, *uint64) (uint64, string, error)
}

type manager struct {
	llmService      llm_grpc.Service
	databaseService databaseService.Service
}

func NewManager(llmService llm_grpc.Service, databaseService databaseService.Service) Manager {
	return &manager{
		llmService:      llmService,
		databaseService: databaseService,
	}
}

func (m *manager) Generate(ctx context.Context, entryPoint string, req string, conversationId *uint64) (uint64, string, error) {
	resp, err := m.llmService.Generate(ctx, req, conversationId)

	// If conversationId is nil, we create a new one
	var convId uint64
	if conversationId == nil {
		convId = 0
		conversationId = &convId
	}

	if err != nil {
		log.Printf("[Generate] Error generating text: %v", err)
		return *conversationId, "", err
	}

	err = m.databaseService.CreateLLMCallReport(ctx, entryPoint, req, resp.GetContent(), float64(resp.GetUsage().GetTotalTokens()))
	if err != nil {
		log.Printf("[Generate] Error creating LLM call report: %v", err)
		return *conversationId, "", err
	}

	metrics.LLMRequestCounter.WithLabelValues(entryPoint).Inc()
	metrics.LLMTokenCounter.WithLabelValues(entryPoint).Add(float64(resp.GetUsage().GetTotalTokens()))
	return resp.GetConversationId(), resp.GetContent(), nil
}
