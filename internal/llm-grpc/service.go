package llm_grpc

import (
	"context"
	suggest "darius/pkg/proto/suggest"
	"log"
)

type Service interface {
	Generate(context.Context, string) (string, error)
}

func NewService(client suggest.ArceusClient, llm_model string) *service {
	return &service{
		client:    client,
		llm_model: llm_model,
	}
}

type service struct {
	client    suggest.ArceusClient
	llm_model string
}

func (s *service) Generate(ctx context.Context, text string) (string, error) {
	res, err := s.client.GenerateText(ctx, &suggest.GenerateTextRequest{
		Content: text,
		Model:   s.llm_model,
	})

	if err != nil {
		log.Printf("Error calling Arceus service: %v", err)
		return "", err
	}

	return res.Content, nil
}
