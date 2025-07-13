package llm_grpc

import (
	"context"
	arceus "darius/pkg/proto/deps/arceus"
	"log"
)

type Service interface {
	Generate(context.Context, string, *uint64) (*arceus.GenerateTextResponse, error)
}

func NewService(client arceus.ArceusClient, llm_model string) Service {
	return &service{
		client:    client,
		llm_model: llm_model,
	}
}

type service struct {
	client    arceus.ArceusClient
	llm_model string
}

func (s *service) Generate(ctx context.Context, text string, conversationId *uint64) (resp *arceus.GenerateTextResponse, err error) {
	res, err := s.client.GenerateText(ctx, &arceus.GenerateTextRequest{
		Content:        text,
		Model:          s.llm_model,
		ConversationId: conversationId,
	})

	if err != nil {
		log.Printf("Error calling Arceus service: %v", err)
		return &arceus.GenerateTextResponse{}, err
	}

	log.Printf("[Generate] LLM request: %s, LLM response %s", text, res)

	return res, err
}
