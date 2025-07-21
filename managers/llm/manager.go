package managers

import (
	"context"
	llm_grpc "darius/internal/services/llm-grpc"
	databaseService "darius/internal/services/repo"
	"darius/metrics"
	"log"
)

type Manager interface {
	Generate(context.Context, string, string, string, *uint64) (*uint64, string, error)
	GetByRequestKey(context.Context, string) (string, error)
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

func (m *manager) Generate(ctx context.Context, entryPoint string, req string, requestKey string, conversationId *uint64) (*uint64, string, error) {
	resp, err := m.llmService.Generate(ctx, req, conversationId)

	if err != nil {
		log.Printf("[Generate] Error generating text: %v", err)
		return nil, "", err
	}

	err = m.databaseService.CreateLLMCallReport(ctx, entryPoint, req, resp.GetContent(), requestKey, float64(resp.GetUsage().GetTotalTokens()))
	if err != nil {
		log.Printf("[Generate] Error creating LLM call report: %v", err)
		return nil, "", err
	}

	metrics.LLMRequestCounter.WithLabelValues(entryPoint).Inc()
	metrics.LLMTokenCounter.WithLabelValues(entryPoint).Add(float64(resp.GetUsage().GetTotalTokens()))

	conID := resp.GetConversationId()
	return &conID, resp.GetContent(), nil
}

func (m *manager) GetByRequestKey(ctx context.Context, requestKey string) (string, error) {
	report, err := m.databaseService.GetByRequestKey(ctx, requestKey)
	if err != nil {
		log.Printf("[GetByRequestKey] Error getting report by request key: %v", err)
		return "", err
	}
	return report, nil
}
