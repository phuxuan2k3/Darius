package bulbasaur

import (
	"context"
	"darius/pkg/proto/deps/bulbasaur"
	"log"
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
		log.Printf("[CheckCallingLLM] Error checking calling LLM: %v", err)
		return false, "Error checking calling LLM: " + err.Error()
	}

	log.Printf("[CheckCallingLLM]: UserID: %d, Amount: %f, Description: %s, Response: %+v", uid, amount, description, res)

	return res.GetTransactionCode() != "", res.GetTransactionCode()
}
func (s *service) ChargeCallingLLM(ctx context.Context, code string) bool {
	res, err := s.client.CommitTransaction(ctx, &bulbasaur.CommitTransactionRequest{
		TransactionCode: code,
	})

	if err != nil {
		log.Printf("[ChargeCallingLLM] Error charging calling LLM: %v", err)
		return false
	}

	log.Printf("[ChargeCallingLLM]: Transaction Code: %s, Response: %+v", code, res)

	return res.GetSuccess()
}
