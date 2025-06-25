package bulbasaur

import (
	"context"
	"darius/pkg/proto/deps/bulbasaur"
)

type Service interface {
	CheckCallingLLM(ctx context.Context, uid uint64, amount float32, description string) (bool, string)
	ChargeCallingLLM(ctx context.Context, code string) bool
}

type service struct {
	client bulbasaur.VenusaurClient
}

func NewService(client bulbasaur.VenusaurClient) Service {
	return &service{
		client: client,
	}
}

func (s *service) CheckCallingLLM(ctx context.Context, uid uint64, amount float32, description string) (bool, string) {
	res, err := s.client.StartTransaction(ctx, &bulbasaur.StartTransactionRequest{
		UserId: uid,
		Amount: amount,
		Note:   description,
	})

	if err != nil {
		return false, "Error checking calling LLM: " + err.Error()
	}

	return res.GetTransactionCode() != "", res.GetTransactionCode()
}
func (s *service) ChargeCallingLLM(ctx context.Context, code string) bool {
	res, err := s.client.CommitTransaction(ctx, &bulbasaur.CommitTransactionRequest{
		TransactionCode: code,
	})

	if err != nil {
		return false
	}

	return res.GetSuccess()
}
