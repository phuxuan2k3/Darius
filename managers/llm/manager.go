package managers

import (
	"context"
	llm_grpc "darius/internal/services/llm-grpc"
	databaseService "darius/internal/services/repo"
	"darius/metrics"
	"log"
)

type Manager interface {
	Generate(context.Context, string, string) (string, error)
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

func (m *manager) Generate(ctx context.Context, entryPoint string, req string) (string, error) {
	resp, err := m.llmService.Generate(ctx, req)
	if err != nil {
		log.Fatalf("Error generating text: %v", err)
		return "", err
	}

	err = m.databaseService.CreateLLMCallReport(ctx, entryPoint, req, resp.GetContent(), float64(resp.GetUsage().GetTotalTokens()))
	if err != nil {
		log.Fatalf("Error creating LLM call report: %v", err)
		return "", err
	}

	metrics.LLMRequestCounter.WithLabelValues(entryPoint).Inc()
	metrics.LLMTokenCounter.WithLabelValues(entryPoint).Add(float64(resp.GetUsage().GetTotalTokens()))
	return resp.GetContent(), nil
}
